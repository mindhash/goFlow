package base


import (
	"encoding/json"
	"fmt"
	"net/http"
 
)

// Simple error implementation wrapping an HTTP response status.
type HTTPError struct {
	Status  int
	Message string
}

func (err *HTTPError) Error() string {
	return fmt.Sprintf("%d %s", err.Status, err.Message)
}

func HTTPErrorf(status int, format string, args ...interface{}) *HTTPError {
	return &HTTPError{status, fmt.Sprintf(format, args...)}
}

// Attempts to map an error to an HTTP status code and message.
// Defaults to 500 if it doesn't recognize the error. Returns 200 for a nil error.
func ErrorAsHTTPStatus(err error) (int, string) {
	if err == nil {
		return 200, "OK"
	}
	switch err := err.(type) {
	case *HTTPError:
		return err.Status, err.Message
	case *json.SyntaxError, *json.UnmarshalTypeError:
		return http.StatusBadRequest, fmt.Sprintf("Invalid JSON: \"%v\"", err)
	}
	Warn("Couldn't interpret error type %T, value %v", err, err)
	return http.StatusInternalServerError, fmt.Sprintf("Internal error: %v", err)
}

// Returns the standard GoFlow error string for an HTTP error status.
// These are important for compatibility, as some REST APIs don't show numeric statuses,
// only these strings.
func GoFlowHTTPErrorName(status int) string {
	switch status {
	case 400:
		return "bad_request"
	case 401:
		return "unauthorized"
	case 404:
		return "not_found"
	case 403:
		return "forbidden"
	case 406:
		return "not_acceptable"
	case 409:
		return "conflict"
	case 412:
		return "file_exists"
	case 415:
		return "bad_content_type"
	}
	return fmt.Sprintf("%d", status)
}

