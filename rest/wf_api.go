package rest

import (
	"net/http"
	"github.com/mindhash/goFlow/base"
	"github.com/mindhash/goFlow/db" 
	"strconv"
)

// HTTP handler for Flow creation
func (h *handler) handlePostFlowDef() error {
	
	// derive flow name from URL path
	flowDefName := h.PathVar("flowDef")
	
	//generate flow def key
	flowDefKey,err := NextFlowDefKey(h,flowDefName )
	
	if flowDefKey != ""{
		return base.HTTPErrorf(http.StatusBadRequest, "Could not generate Flow Def Key") 
	}
	
	wf := db.NewWorkflowDef(flowDefKey)
	
	_, err = h.readObject(wf) 
	if err != nil { 
		return base.HTTPErrorf(http.StatusBadRequest, "Flow Definition not found")   //TO DO: Need to relook at error
	}  
	
	h.writeJSONStatus(http.StatusCreated, db.Body{"ok": true, "flowDefKey": flowDefKey})
	return nil
}

// HTTP handler for Flow Update
func (h *handler) handlePutFlowDef() error {
	
	// derive flow def key from  URL path
	flowDefKey := h.PathVar("flowDefKey")
	base.Logf("...Flow Definition Key...", flowDefKey)
	
	if flowDefKey != "" {
		return base.HTTPErrorf(http.StatusBadRequest, "Invalid Flow Def Key") 
	}
	
	wf := db.NewWorkflowDef(flowDefKey)
	
	_, err := h.readObject(wf) 
	if err != nil { 
		return base.HTTPErrorf(http.StatusBadRequest, "Flow Definition not found")   //TO DO: Need to relook at error
	}  
	
	h.writeJSONStatus(http.StatusCreated, db.Body{"ok": true, "flowDefKey": flowDefKey})
	return nil
}


func NextFlowDefKey (h *handler, flowName string) (string,error) {
	newVersionNum := 1.0
	
	// get flow def last version    
	lastVersionStr,err := h.db.GetValue("_flow:" + flowName + ":_lastversion")	
	
	if (lastVersionStr != "") {
    	// convert version string to num 
		lastVersionNum,err := strconv.ParseFloat(lastVersionStr, 64)  
		if err != nil {
			return "", err
		}
		// add 1 to last version
		newVersionNum = lastVersionNum + 1
	}  
	
	return "_flow:" + flowName + ":_version:" + base.FloatToString(newVersionNum), err
  
}