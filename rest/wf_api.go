package rest

import (
	"net/http"
	"encoding/json"
	"github.com/mindhash/goFlow/base"
	"github.com/mindhash/goFlow/db" 
	"strconv"
)

// HTTP handler for Flow creation
func (h *handler) handlePostFlowDef() error {
	
	// derive flow name from URL path
	flowDefName := h.PathVar("flowDef")
	base.Logf("...Flow Definition Name...", flowDefName)
	

	//generate flow def key
	flowDefKey,err := NextFlowDefKey(h.db,flowDefName )
	base.Logf("Flow Def Key:%s",flowDefKey)
	
	if flowDefKey == ""{
		return base.HTTPErrorf(http.StatusBadRequest, "Could not generate Flow Def Key") 
	}
	
	wf := db.NewWorkflowDef(flowDefKey)
	
	_, err = h.readObject(wf) 
	if err != nil { 
		return base.HTTPErrorf(http.StatusBadRequest, "Flow Definition not found")   //TO DO: Need to relook at error
	}  
	
	//base.Logf("Encoding Workflow Def ..%s",wf) 
	data, _ := json.Marshal(wf)
	
	//base.Logf("Saving Workflow Def ...%s",wf) 
	_, _ = h.db.PutDocRaw (flowDefKey, data)
	//base.Logf("Done save...") 
	
	// this should be inside db trx TO DO
	//_, err =writeLastVersion(h.db , wf.Name , )
	
	h.writeJSONStatus(http.StatusCreated, db.Body{"ok": true, "flowDefKey": flowDefKey})
	return nil
}

// HTTP handler for Flow Update
func (h *handler) handlePutFlowDef() error {
	
	// derive flow def key from  URL path
	flowDefKey := h.PathVar("flowDefKey")
	base.Logf("Flow Definition Key: ", flowDefKey)
	
	if flowDefKey == "" {
		return base.HTTPErrorf(http.StatusBadRequest, "Invalid Flow Def Key") 
	}
	// 
	wf := db.NewWorkflowDef(flowDefKey)
	
	// read JSON HTTP input into Object	
	_, err := h.readObject(&wf)
 	
	//base.Logf("Encoding Workflow Def ..%s",wf) 
	data, _ := json.Marshal(wf)
	
	//base.Logf("Saving Workflow Def ...%s",wf) 
	_, err := h.db.PutDocRaw (flowDefKey, data)
	//base.Logf("Done save...") 
	 
	if err != nil { 
		return base.HTTPErrorf(http.StatusBadRequest, "Flow Definition Could not be updated")   //TO DO: Need to relook at error
	}  
	// this should be inside db trx TO DO
	//_, err =writeLastVersion(h.db , wf.Name , )
	
	//db.Body{"ok": true, "flowDefKey": flowDefKey}
	h.writeJSONStatus(http.StatusCreated, &wf)
	return nil
}

// HTTP handler for Flow query
func (h *handler) handleGetFlowDef() error {
	// derive flow def key from  URL path
	flowDefKey := h.PathVar("flowDefKey")
	base.Logf("Flow Definition Key: ", flowDefKey)
	
	if flowDefKey == "" {
		return base.HTTPErrorf(http.StatusBadRequest, 
			"Invalid Flow Def Key") 
	}
	
	data, err := h.db.GetDocRaw (flowDefKey)
	if (err != nil){
		return base.HTTPErrorf(http.StatusBadRequest, 
			"WF Definition Query Failed")
	}
	
	if (data == nil){
		return base.HTTPErrorf(http.StatusBadRequest, 
			"WF Definition Not Found. Check Flow Def Key again.")
	}
	
	wf := db.NewWorkflowDef(flowDefKey) 
	  
	err = json.Unmarshal([]byte(data),&wf)
 	
	if (err != nil) {
		return err
	}
	
	h.writeJSONStatus(http.StatusCreated,wf )
	return nil
}

func writeLastVersion(db *Database , flowName string, version string) error{
	saved, err := d.PutValue("_flow:" + flowName + ":_lastversion" , version)
	return err
}
 
// get new key ID for flow definition
func NextFlowDefKey (d *db.Database, flowName string) (string,error) {
	newVersionNum := 1.0
	
	// get flow def last version    
	lastVersionStr,err := d.GetValue("_flow:" + flowName + ":_lastversion")	
	
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