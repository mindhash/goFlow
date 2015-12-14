package rest

import (
	"github.com/gorilla/mux"
	"net/http" 
	"regexp"
	"strings"
)


// Regexes that match doc ID component of a path. 
const docRegex = "[^_/][^/]*"

// Regex that matches a URI containing a regular doc ID with an escaped "/" character
var docWithSlashPathRegex *regexp.Regexp

func init() {
	docWithSlashPathRegex, _ = regexp.Compile("/"  + "/[^_].*%2[fF]")
}


func createHandler(sc *ServerContext, privs handlerPrivs) (*mux.Router) {
	r := mux.NewRouter()
	r.StrictSlash(true)
	
	// Global operations:
	r.Handle("/", makeHandler(sc, privs, (*handler).handleRoot)).Methods("GET", "HEAD")
	
	//flow definition
	r.Handle("/_flow/{flowDefKey}", makeHandler(sc, privs, (*handler).handlePutFlowDef)).Methods("PUT")
	r.Handle("/_flow/create/{flowName}", makeHandler(sc, privs, (*handler).handlePostFlowDef)).Methods("POST")
	r.Handle("/_flow/{flowDefKey}", makeHandler(sc, privs, (*handler).handleGetFlowDef)).Methods("GET")
	//r.Handle("/_flow/name/{flowDef}", makeHandler(sc, privs, (*handler).handleGetByNameFlowDef)).Methods("GET")

	//handle flow instance creation
	r.Handle("/_flow/{flowName}/_transaction/create", makeHandler(sc, privs, (*handler).handlePostFlowTxn)).Methods("POST")
	
	return r
}

// Creates the HTTP handler for the public API of a gateway server.
func CreatePublicHandler(sc *ServerContext) http.Handler {
	r := createHandler(sc, regularPrivs)
	
	return wrapRouter(sc, regularPrivs, r)
}

// Returns a top-level HTTP handler for a Router. This adds behavior for URLs that don't
// match anything -- it handles the OPTIONS method as well as returning either a 404 or 405
// for URLs that don't match a route.
func wrapRouter(sc *ServerContext, privs handlerPrivs, router *mux.Router) http.Handler {
	return http.HandlerFunc(func(response http.ResponseWriter, rq *http.Request) {
		fixQuotedSlashes(rq)
		var match mux.RouteMatch

		// Inject CORS if enabled and requested and not admin port
		//originHeader := rq.Header["Origin"]
		//!if privs != adminPrivs && sc.config.CORS != nil && len(originHeader) > 0 {
		//	origin := matchedOrigin(sc.config.CORS.Origin, originHeader)
		//	response.Header().Add("Access-Control-Allow-Origin", origin)
		//	response.Header().Add("Access-Control-Allow-Credentials", "true")
		//	response.Header().Add("Access-Control-Allow-Headers", strings.Join(sc.config.CORS.Headers, ", "))
		//}

		if router.Match(rq, &match) {
			router.ServeHTTP(response, rq)
		} else {
			// Log the request
			h := newHandler(sc, privs, response, rq)
			h.logRequestLine()

			// What methods would have matched?
			var options []string
			for _, method := range []string{"GET", "HEAD", "POST", "PUT", "DELETE"} {
				if wouldMatch(router, rq, method) {
					options = append(options, method)
				}
			}
			if len(options) == 0 {
				h.writeStatus(http.StatusNotFound, "unknown URL")
			} else {
				response.Header().Add("Allow", strings.Join(options, ", "))
				//if privs != adminPrivs && sc.config.CORS != nil && len(originHeader) > 0 {
				//	response.Header().Add("Access-Control-Max-Age", strconv.Itoa(sc.config.CORS.MaxAge))
				//	response.Header().Add("Access-Control-Allow-Methods", strings.Join(options, ", "))
				//}
				if rq.Method != "OPTIONS" {
					h.writeStatus(http.StatusMethodNotAllowed, "")
				} else {
					h.writeStatus(http.StatusNoContent, "")
				}
			}
			h.logDuration(true)
		}
	})
}

func matchedOrigin(allowOrigins []string, rqOrigins []string) string {
	for _, rv := range rqOrigins {
		for _, av := range allowOrigins {
			if rv == av {
				return av
			}
		}
	}
	for _, av := range allowOrigins {
		if av == "*" {
			return "*"
		}
	}
	return ""
}

func fixQuotedSlashes(rq *http.Request) {
	uri := rq.RequestURI
	if docWithSlashPathRegex.MatchString(uri) {
		if stop := strings.IndexAny(uri, "?#"); stop >= 0 {
			uri = uri[0:stop]
		}
		rq.URL.Path = uri
	}
}

func wouldMatch(router *mux.Router, rq *http.Request, method string) bool {
	savedMethod := rq.Method
	rq.Method = method
	defer func() { rq.Method = savedMethod }()
	var matchInfo mux.RouteMatch
	return router.Match(rq, &matchInfo)
}
