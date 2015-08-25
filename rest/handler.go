package rest


import (
	"bytes"
	"encoding/json"
	"expvar"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gorilla/mux"
 
	"github.com/mindhash/goFlow/base"
	"github.com/mindhash/goFlow/db"
	"github.com/mindhash/goFlow/auth"
)


var kNotFoundError = base.HTTPErrorf(http.StatusNotFound, "missing")
var kBadMethodError = base.HTTPErrorf(http.StatusMethodNotAllowed, "Method Not Allowed")
var kBadRequestError = base.HTTPErrorf(http.StatusMethodNotAllowed, "Bad Request")

var restExpvars = expvar.NewMap("goflow_rest")

// If set to true, JSON output will be pretty-printed.
var PrettyPrint bool = false

// If set to true, diagnostic data will be dumped if there's a problem with MIME multipart data
var DebugMultipart bool = false

var lastSerialNum uint64 = 0

func init() {
	DebugMultipart = (os.Getenv("GatewayDebugMultipart") != "")
}

// Encapsulates the state of handling an HTTP request.
type handler struct {
	server         *ServerContext
	rq             *http.Request
	response       http.ResponseWriter
	status         int
	statusMessage  string
	requestBody    io.ReadCloser
	db             *db.Database
	user           auth.User
	privs          handlerPrivs
	startTime      time.Time
	serialNumber   uint64
	loggedDuration bool
}

type handlerPrivs int

const (
	regularPrivs = iota // Handler requires authentication
	publicPrivs         // Handler checks auth but doesn't require it
	adminPrivs          // Handler ignores auth, always runs with root/admin privs
)

type handlerMethod func(*handler) error

// Creates an http.Handler that will run a handler with the given method
func makeHandler(server *ServerContext, privs handlerPrivs, method handlerMethod) http.Handler {
	return http.HandlerFunc(func(r http.ResponseWriter, rq *http.Request) {
		h := newHandler(server, privs, r, rq)
		err := h.invoke(method)
		h.writeError(err)
		h.logDuration(true) 
	})
}

func newHandler(server *ServerContext, privs handlerPrivs, r http.ResponseWriter, rq *http.Request) *handler {
	return &handler{
		server:       server,
		privs:        privs,
		rq:           rq,
		response:     r,
		status:       http.StatusOK,
		serialNumber: atomic.AddUint64(&lastSerialNum, 1),
		startTime:    time.Now(),
	}
}

// Top-level handler call. It's passed a pointer to the specific method to run.
func (h *handler) invoke(method handlerMethod) error {
	 
	restExpvars.Add("requests_total", 1)
	restExpvars.Add("requests_active", 1)
	defer restExpvars.Add("requests_active", -1)

	switch h.rq.Header.Get("Content-Encoding") {
	case "":
		h.requestBody = h.rq.Body
	default:
		return base.HTTPErrorf(http.StatusUnsupportedMediaType, "Unsupported Content-Encoding;")
	}
	
	var err error
	
	// If there is a "db" path variable, look up the database context:
	var dbContext *db.DatabaseContext 
	if dbContext, err = h.server.GetDatabase(); err != nil {
			h.logRequestLine()
			return err
		}
	
		
	// Authenticate, if not on admin port:
	//TO DO: Add authorization for DB Context 
	//if h.privs != adminPrivs {
	//	if err = h.checkAuth(dbContext); err != nil {
	//		h.logRequestLine()
	//		return err
	//	}
	//}

	h.logRequestLine()

	// Now set the request's Database (i.e. context + user)
	if dbContext != nil {
		h.db, err = db.GetDatabase(dbContext, nil)
		
		if err != nil {
			return err
		}
	}
	 
	return method(h) // Call the actual handler code
	
}

func (h *handler) logRequestLine() {
	if !base.LogKeys["HTTP"] {
		return
	}
	as := ""
	if h.privs == adminPrivs {
		as = "  (ADMIN)"
	} else if h.user != nil && h.user.Name() != "" {
		as = fmt.Sprintf("  (as %s)", h.user.Name())
	}
	base.LogTo("HTTP", " #%03d: %s %s%s", h.serialNumber, h.rq.Method, h.rq.URL, as)
}

func (h *handler) logDuration(realTime bool) {
	if h.loggedDuration {
		return
	}
	h.loggedDuration = true

	var duration time.Duration
	if realTime {
		duration = time.Since(h.startTime)
		bin := int(duration/(100*time.Millisecond)) * 100
		restExpvars.Add(fmt.Sprintf("requests_%04dms", bin), 1)
	}

	logKey := "HTTP+"
	if h.status >= 300 {
		logKey = "HTTP"
	}
	base.LogTo(logKey, "#%03d:     --> %d %s  (%.1f ms)",
		h.serialNumber, h.status, h.statusMessage,
		float64(duration)/float64(time.Millisecond))
}

// Used for indefinitely-long handlers like _changes that we don't want to track duration of
func (h *handler) logStatus(status int, message string) {
	h.setStatus(status, message)
	h.logDuration(false) // don't track actual time
}

func (h *handler) checkAuth(context *db.DatabaseContext) error {
	h.user = nil
	if context == nil {
		return nil
	}
	// add more details later
	return nil
}


func (h *handler) PathVar(name string) string {
	v := mux.Vars(h.rq)[name]

	//Escape special chars i.e. '+' otherwise they are removed by QueryUnescape()
	v = strings.Replace(v, "+", "%2B", -1)

	// Before routing the URL we explicitly disabled expansion of %-escapes in the path
	// (see function fixQuotedSlashes). So we have to unescape them now.
	v, _ = url.QueryUnescape(v)
	return v
}

func (h *handler) SetPathVar(name string, value string) {
	mux.Vars(h.rq)[name] = url.QueryEscape(value)
}

func (h *handler) getQuery(query string) string {
	return h.rq.URL.Query().Get(query)
}

func (h *handler) getBoolQuery(query string) bool {
	return h.getOptBoolQuery(query, false)
}

func (h *handler) getOptBoolQuery(query string, defaultValue bool) bool {
	q := h.getQuery(query)
	if q == "" {
		return defaultValue
	}
	return q == "true"
}

// Returns the integer value of a URL query, defaulting to 0 if unparseable
func (h *handler) getIntQuery(query string, defaultValue uint64) (value uint64) {
	return getRestrictedIntQuery(h.rq.URL.Query(), query, defaultValue, 0, 0, false)
}

func (h *handler) getJSONQuery(query string) (value interface{}, err error) {
	valueJSON := h.getQuery(query)
	if valueJSON != "" {
		err = json.Unmarshal([]byte(valueJSON), &value)
	}
	return
}

func (h *handler) userAgentIs(agent string) bool {
	userAgent := h.rq.Header.Get("User-Agent")
	return len(userAgent) > len(agent) && userAgent[len(agent)] == '/' && strings.HasPrefix(userAgent, agent)
}

// Returns the request body as a raw byte array.
func (h *handler) readBody() ([]byte, error) {
	return ioutil.ReadAll(h.requestBody)
}

// Parses a JSON request body, returning it as a Body map.
func (h *handler) readJSON() (db.Body, error) {
	var body db.Body
	return body, h.readJSONInto(&body)
}

// Parses a JSON request body into a custom structure.
func (h *handler) readJSONInto(into interface{}) error {
	
	contentType := h.rq.Header.Get("Content-Type")
	if contentType != "" && !strings.HasPrefix(contentType, "application/json") {
		return base.HTTPErrorf(http.StatusUnsupportedMediaType, "Invalid content type %s", contentType)
	}
 
 	//TO DO: zip version to be added
	
	decoder := json.NewDecoder(h.requestBody)
	if err := decoder.Decode(into); err != nil {
		base.Warn("Couldn't parse JSON in HTTP request: %v", err)
		return base.HTTPErrorf(http.StatusBadRequest, "Bad JSON")
	}
	return nil
}

//TO DO: Need to add multi part reads
//This function handles marshaling of input JSON into Struct
func (h *handler) readObject(obj interface{}) ( interface{}, error){
	
	contentType, _ , _ := mime.ParseMediaType(h.rq.Header.Get("Content-Type"))
	
	//process JSON Documents only
	switch contentType {
		
	case "", "application/json":
		return obj, h.readJSONInto(obj)
	default:
		return nil, base.HTTPErrorf(http.StatusUnsupportedMediaType, "Invalid content type %s", contentType)
	}
}

// Reads & parses the request body, handling either JSON or multipart.
func (h *handler) readDocument() (db.Body, error) {
	//!contentType,  attrs, _ := mime.ParseMediaType(h.rq.Header.Get("Content-Type"))
	contentType, _ , _ := mime.ParseMediaType(h.rq.Header.Get("Content-Type"))
	switch contentType {
	case "", "application/json":
		return h.readJSON()
	//case "multipart/related":
	//	if DebugMultipart {
	//		raw, err := h.readBody()
	//		if err != nil {
	//			return nil, err
	//		}
	//		reader := multipart.NewReader(bytes.NewReader(raw), attrs["boundary"])
	//		body, err := db.ReadMultipartDocument(reader)
	//		if err != nil {
	//			ioutil.WriteFile("GatewayPUT.mime", raw, 0600)
	//			base.Warn("Error reading MIME data: copied to file GatewayPUT.mime")
	//		}
	//		return body, err
	//	} else {
	//		reader := multipart.NewReader(h.requestBody, attrs["boundary"])
	//		return db.ReadMultipartDocument(reader)
	//	}
	default:
		return nil, base.HTTPErrorf(http.StatusUnsupportedMediaType, "Invalid content type %s", contentType)
	}
}

func (h *handler) requestAccepts(mimetype string) bool {
	accept := h.rq.Header.Get("Accept")
	return accept == "" || strings.Contains(accept, mimetype) || strings.Contains(accept, "*/*")
}

//Responses

func (h *handler) setHeader(name string, value string) {
	h.response.Header().Set(name, value)
}

func (h *handler) setStatus(status int, message string) {
	h.status = status
	h.statusMessage = message
}


// Writes an object to the response in JSON format.
// If status is nonzero, the header will be written with that status.
func (h *handler) writeJSONStatus(status int, value interface{}) {
	if !h.requestAccepts("application/json") {
		base.Warn("Client won't accept JSON, only %s", h.rq.Header.Get("Accept"))
		h.writeStatus(http.StatusNotAcceptable, "only application/json available")
		return
	}

	jsonOut, err := json.Marshal(value)
	if err != nil {
		base.Warn("Couldn't serialize JSON for %v : %s", value, err)
		h.writeStatus(http.StatusInternalServerError, "JSON serialization failed")
		return
	}
	if PrettyPrint {
		var buffer bytes.Buffer
		json.Indent(&buffer, jsonOut, "", "  ")
		jsonOut = append(buffer.Bytes(), '\n')
	}
	h.setHeader("Content-Type", "application/json")
	if h.rq.Method != "HEAD" {
		//if len(jsonOut) < 1000 {
		//	h.disableResponseCompression()
		//}
		h.setHeader("Content-Length", fmt.Sprintf("%d", len(jsonOut)))
		if status > 0 {
			h.response.WriteHeader(status)
			h.setStatus(status, "")
		}
		h.response.Write(jsonOut)
	} else if status > 0 {
		h.response.WriteHeader(status)
		h.setStatus(status, "")
	}
}

func (h *handler) writeJSON(value interface{}) {
	h.writeJSONStatus(http.StatusOK, value)
}

func (h *handler) addJSON(value interface{}) {
	encoder := json.NewEncoder(h.response)
	err := encoder.Encode(value)
	if err != nil {
		base.Warn("Couldn't serialize JSON for %v : %s", value, err)
		panic("JSON serialization failed")
	}
}

func (h *handler) writeMultipart(subtype string, callback func(*multipart.Writer) error) error {
	if !h.requestAccepts("multipart/") {
		return base.HTTPErrorf(http.StatusNotAcceptable, "Response is multipart")
	}

	// Get the output stream. Due to a CouchDB bug, if we're sending to it we need to buffer the
	// output in memory so we can trim the final bytes.
	var output io.Writer
	var buffer bytes.Buffer
	if h.userAgentIs("CouchDB") {
		output = &buffer
	} else {
		output = h.response
	}

	writer := multipart.NewWriter(output)
	h.setHeader("Content-Type",
		fmt.Sprintf("multipart/%s; boundary=%q", subtype, writer.Boundary()))

	err := callback(writer)
	writer.Close()

	if err == nil && output == &buffer {
		// Trim trailing newline; CouchDB is allergic to it:
		_, err = h.response.Write(bytes.TrimRight(buffer.Bytes(), "\r\n"))
	}
	return err
}

func (h *handler) flush() {
	switch r := h.response.(type) {
	case http.Flusher:
		r.Flush()
	}
}

// If the error parameter is non-nil, sets the response status code appropriately and
// writes a CouchDB-style JSON description to the body.
func (h *handler) writeError(err error) {
	if err != nil {
		status, message := base.ErrorAsHTTPStatus(err)
		h.writeStatus(status, message)
	}
}

// Writes the response status code, and if it's an error writes a JSON description to the body.
func (h *handler) writeStatus(status int, message string) {
	if status < 300 {
		h.response.WriteHeader(status)
		h.setStatus(status, message)
		return
	}
	// Got an error:
	var errorStr string
	switch status {
	case http.StatusNotFound:
		errorStr = "not_found"
	case http.StatusConflict:
		errorStr = "conflict"
	default:
		errorStr = http.StatusText(status)
		if errorStr == "" {
			errorStr = fmt.Sprintf("%d", status)
		}
	}

//	h.disableResponseCompression()
	h.setHeader("Content-Type", "application/json")
	h.response.WriteHeader(status)
	h.setStatus(status, message)
	jsonOut, _ := json.Marshal(db.Body{"error": errorStr, "reason": message})
	h.response.Write(jsonOut)
}


// Returns the integer value of a URL query, restricted to a min and max value,
// but returning 0 if missing or unparseable.  If allowZero is true, values coming in
// as zero will stay zero, instead of being set to the minValue.
func getRestrictedIntQuery(values url.Values, query string, defaultValue, minValue, maxValue uint64, allowZero bool) uint64 {
	return getRestrictedIntFromString(
		values.Get(query),
		defaultValue,
		minValue,
		maxValue,
		allowZero,
	)
}

func getRestrictedIntFromString(rawValue string, defaultValue, minValue, maxValue uint64, allowZero bool) uint64 {
	var value *uint64
	if rawValue != "" {
		intValue, err := strconv.ParseUint(rawValue, 10, 64)
		if err != nil {
			value = nil
		} else {
			value = &intValue
		}
	}

	return getRestrictedInt(
		value,
		defaultValue,
		minValue,
		maxValue,
		allowZero,
	)
}



func getRestrictedInt(rawValue *uint64, defaultValue, minValue, maxValue uint64, allowZero bool) uint64 {

	var value uint64

	// Only use the defaultValue if rawValue isn't specified.
	if rawValue == nil {
		value = defaultValue
	} else {
		value = *rawValue
	}

	// If value is zero and allowZero=true, leave value at zero rather than forcing it to the minimum value
	validZero := (value == 0 && allowZero)
	if value < minValue && !validZero {
		value = minValue
	}

	if value > maxValue && maxValue > 0 {
		value = maxValue
	}

	return value
}
