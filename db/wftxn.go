package db

import (
	"time"
	"github.com/mindhash/goFlow/base"
)

type Activity struct {
	ActivityKey string
	Name string
	ActivityInstanceKey string
	AccessKey string  // used for authorization. shared with worker at the time of job assignment
	CreationDate time.Time
	UpdateDate time.Time
	PlannedStart time.Time 
	PlannedEnd time.Time 
	ActualStart time.Time 
	ActualEnd time.Time 
	Status string
}



type Attribute struct {
	Name string
	Type string    // String, Byte
	StrinValue string   /// need 
	JsonValue []byte
}

type FlowInstance struct{
	InstanceKey string
	Version string
	FlowDefKey string
	FlowName string
	FlowVersion string
	UserInstanceKey string
	Status string
	CreationDate time.Time
	UpdateDate time.Time
	StartDate time.Time
	EndDate time.Time
	Activities []Activity
	Attributes []Attribute
}



type FlowTxnRequest struct {
	Name string
	InstanceKey string
	FlowDefKey string
	FlowVersion string
	UserInstanceKey string 
	AttributeValues []Attribute
}

func NewFlowTxnRequest() *FlowTxnRequest{
	return &FlowTxnRequest{}
}

// try to make this immuatable TO DO
func NewFlowInstance (ftr *FlowTxnRequest, Activites []Activity ) *FlowInstance{	
	newInstanceKey    := base.CreateUUID() 

	return &FlowInstance { FlowDefKey: ftr.FlowDefKey,  InstanceKey: newInstanceKey, Status: "Scheduled",CreationDate: time.Now()}
}



func NewActivity(actkey string) *Activity{
	// need to derive ActivityInstanceKey 
	newActInstanceId := base.CreateUUID()
	return &Activity {ActivityInstanceKey: newActInstanceId, ActivityKey: actkey, Status: "Scheduled",CreationDate: time.Now()}
} 