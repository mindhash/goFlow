package rest

import (
	"net/http"
	"encoding/json"
	"github.com/mindhash/goFlow/base"
	"github.com/mindhash/goFlow/db" 
	//"strconv"
)



// post received a flow Txn Request which is used to generate flow instance

func (h *handler) handlePostFlowTxn() error {
	
	base.Logf("handlePostFlowTxn...")
	 
	flowInstanceReq := db.NewFlowTxnRequest()
	
	_, err := h.readObject(flowInstanceReq)  //read JSON request into ftr object
	base.Logf("Read Object...")
	 
	if err != nil { 
		return base.HTTPErrorf(http.StatusBadRequest, "Could not map JSON to FlowTxnRequest Object")   
	}  

	//use flow def key to retrive workflow raw data and convert it to wf object 
	flowDefData,_ := h.db.GetDocRaw(flowInstanceReq.FlowDefKey)
	flowDef := &db.WorkflowDef{}
	_ =json.Unmarshal(flowDefData, flowDef)


    flowInstance := db.NewFlowInstance(flowInstanceReq,flowDef)
	base.Logf("Flow Instance Key:%s",flowInstance.InstanceKey )
	

	// Save Flow Instance to database
	data, _ := json.Marshal(flowInstance) 
	_, _ = h.db.PutDocRaw (flowInstance.InstanceKey , data)

	//JSON output created flow instance  
	h.writeJSONStatus(http.StatusCreated, &flowInstance)
	return nil
}



func (h *handler) handlePutFlowTxn() error {
	return nil
}
	

func (h *handler) handleGetFlowTxn() error {
return nil
}
	
	
