package db

import (
	// "encoding/json"
	//"github.com/mindhash/goFlow/base"
	
	)

type AttributeDef struct {
	Name string
	Type string    // String, Byte 
}

type ActivitiesDef struct {
	ActivityKey string
	Name string
	Inactive bool  
}


type WorkflowDef  struct { 
	FlowDefKey string   // Def key is combination of flow name and version 
	Name string
	Version string
	Activities  []ActivitiesDef 
	 Attributes []AttributeDef
}

// Returns a new empty document.
func NewWorkflowDef(flowdefkey string) *WorkflowDef {
	return &WorkflowDef{FlowDefKey: flowdefkey}
}

// TO DO: Error handling 
// Returns Next Flow Definition Key (Format: _flow: <FLOW NAME> : _version: <VERSION>)
// Uses _flow:<FLOW NAME>:_lastversion:<LAST_VERSION> key to derive recent version from DB
// 


/*
func (wf *WorkflowDef) Save(db *Database) (bool, error) {
	
	data,err := json.Marshal(wf ) 
	
	if err != nil {
		base.Warn("Error marshaling body of workflow %q: %s", wf.Name, err)
		return false,err
	}
	
	saved,err := db.PutDocRaw(wf.ID, data)  
	
	if err != nil {  		
		base.Warn("Error saving workflow %q: %s", &wf.Name,  err)
		return saved, err
	}
	
	return saved, nil
	 
}

func (wf *WorkflowDef) UnmarshalJSON(data []byte) error {
	
	err := json.Unmarshal([]byte(data), &wf)
	if err != nil {
		base.Warn("Error unmarshaling body of workflow : %s",  err)
		return err
	}

	return nil
	
}

func (wf *WorkflowDef) MarshalJSON() ([]byte, error) {
	data,err := json.Marshal(&wf )
	
	if err != nil {
		base.Warn("Error marshaling body of workflow %q: %s", wf.Name, err)
		return nil,err
	}
	return data, err
}*/