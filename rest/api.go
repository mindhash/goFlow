package rest

import (
	"encoding/json"
	"net/http"
)

const ServerName = "GoFlow - Simple Workflow"
const VersionNumber float64 = 1.0                    // API/feature level


var LongVersionString string

func init(){
	LongVersionString = fmt.Sprintf("%s/unofficial", ServerName)
}

// HTTP handler for the root ("/")
func (h *handler) handleRoot() error {
	response := map[string]interface{}{
		"goflow": "Welcome",
		"version": LongVersionString,
		"vendor":  db.Body{"name": ServerName, "version": VersionNumber},
	}
	//if h.privs == adminPrivs {
	//	response["ADMIN"] = true
	//}
	h.writeJSON(response)
	return nil
}