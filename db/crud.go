package db

import (
	"encoding/json"
	"github.com/mindhash/goFlow/base"
	)

func realDocID(docid string) string {
	if len(docid) > 250 {
		return "" // Invalid doc IDs
	}
	return docid
}

// get raw data from DB 
// will be frequently used by ORM Objects
func (db *DatabaseContext) GetDocRaw(key string) (v []byte, err error) {
	v, err = db.DbHandle.GetRaw (key)
	if err != nil {
		return nil, err
	} 
	return   
}
// put raw data into db
// will be frequently used by ORM objects
func (db *Database) PutDocRaw(key string, v []byte) (added bool, err error) {
	added, err = db.DbHandle.PutRaw (key,v)
	if err != nil {
		return false, err
	} 
	return   
}

// create new object
// return object key
func (db *Database) PostDocRaw(v []byte) (key string, err error) {

	// If there's an incoming _id property, use that as the doc ID.
	//TO DO: will not work as v is byte array not JSON key, idFound := v["_id"].(string)
	//if !idFound {
	//	key = base.CreateUUID()
	//}
	
	key = base.CreateUUID()
	
	_, err = db.DbHandle.PutRaw(key,v )
	if err != nil {
		key = ""
	}
	return key, err
}


// Get Document TO DO: Need to see if this is needed
// need to check if DBC or Database 
func (db *DatabaseContext) GetDoc(docid string) (*document, error) {
	key := realDocID(docid)
	if key == "" {
		return nil, base.HTTPErrorf(400, "Invalid doc ID")
	}
	dbExpvars.Add("document_gets", 1)// need to check this
	doc := newDocument(docid)
	
	//get data from store
	data, err := db.DbHandle.GetRaw(key)
	if err != nil {
		return nil, err
	} 
	//unmarshal JSON into doc.body
	doc.UnmarshalJSON(data)
	
	if err != nil {
		return nil, err
	}
	return doc, nil
}


//TO DO: Need to see if this is needed
// Creates a new document, assigning it a random doc ID.
func (db *Database) PostDoc(body Body) (string, error) {

	// If there's an incoming _id property, use that as the doc ID.
	docid, idFound := body["_id"].(string)
	if !idFound {
		docid = base.CreateUUID()
	}
	jbody,err := json.Marshal(body)
	
	if err != nil {
		docid = ""
		return docid, err
	}
	
	_, err = db.DbHandle.PutRaw(docid,jbody )
	if err != nil {
		docid = ""
	}
	return docid, err
}

//TO DO: Need to revisit
// Deletes a document, by adding a new revision whose "_deleted" property is true.
func (db *Database) DeleteDoc(docid string, revid string) (string, error) {
	body := Body{"_deleted": true, "_rev": revid}
	jbody,err := json.Marshal(body)
	
	if err != nil {
		docid = ""
		return docid, err
	}
	
	_, err = db.DbHandle.PutRaw(docid,jbody )
	if err != nil {
		docid = ""
	}
	
	return docid, err
}

// get plain values from DB
func (db *Database) GetValue (key string) (string, error) {
	keyVal, err := db.DbHandle.GetRaw(key )
	return string(keyVal), err
}

// Put plain values from DB
func (db *Database) PutValue (key string, val string) (bool, error) {

	retval, err := db.DbHandle.PutRaw(key,[]byte(val))
	return retval, err
}


