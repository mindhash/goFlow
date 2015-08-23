package db

import ("testing"
		"log"
		"runtime"
		"strings"
		"fmt"
		"encoding/json"
		"github.com/mindhash/go.assert"
		"github.com/mindhash/goFlow/base"
)





func setupTestDB(t *testing.T) *Database {
	handle,err := openDatabase(base.DatabaseSpec{
		Directory:     "testDB.db",
		BucketName: "TestBucket",
		Name: "db"})
	 
	if err != nil {
			log.Fatalf("Couldn't connect to bucket: %v", err)
	}
		
	context, err := NewDatabaseContext(handle.GetName(), handle)
	assertNoError(t, err, "Couldn't create context for database 'db'")
	db, err := CreateDatabase(context)
	assertNoError(t, err, "Couldn't create database 'db'")
	return db
}

func tearDownTestDB(t *testing.T, db *Database) {
	db.Close()
}

func TestDatabase(t *testing.T) {
	db := setupTestDB(t)
	defer tearDownTestDB(t, db)

	// Test creating & updating a document:
	log.Printf("Create rev 1...")
	body := Body{"key1": "value1", "key2": 1234}
	
	docid1, err := db.PostDoc(body)
	assertNoError(t, err, "Couldn't create document")
	log.Printf("Create rev 1 Doc ID...%s",docid1)
	
	log.Printf("Retrieve doc...")
	doc,err := db.GetDoc(docid1)
	assertNoError(t, err, "Couldn't get Doc Body")
	 
	assert.DeepEquals(t, doc.body, body) 

}

func TestWorkflow (t *testing.T){
	db := setupTestDB(t)
	defer tearDownTestDB(t, db)
	
	wfid := "sndklee3fndsdsf"
	
	
	log.Printf("Create Workflow Def 1...")
	wf := newWorkflowDef (wfid)
	
	wf.Name ="Test Workflow"
	
	log.Printf("Encoding Workflow Def 1...%s",wf) 
	data, _ := json.Marshal(wf)
	
	log.Printf("Saving Workflow Def 1...%s",wf) 
	_, _ = db.PutDocRaw (wfid, data)
	log.Printf("Done save...") 
	
	saveddata,_ := db.GetDocRaw  (wfid)
	
	log.Printf("Retrive Encoded WF Data...%s", saveddata)
	
	_ = json.Unmarshal(saveddata, wf)
	
	log.Printf("Retrive  WF Data...%s", wf)	
	 
}

func TestWorkflow1 (t *testing.T){
	db := setupTestDB(t)
	defer tearDownTestDB(t, db)
	
	body := []byte(`{ "Name":"Test Workflow 1","Version":"1.0"}`)
	
	wf := newWorkflowDef("123dsasf")
	
	_ = json.Unmarshal(body, wf)
	log.Printf("Retrive  WF Data...%s", wf)	
	
	
}

//////// HELPERS:

func assertFailed(t *testing.T, message string) {
	_, file, line, ok := runtime.Caller(2) // assertFailed + assertNoError + public function.
	if ok {
		// Truncate file name at last file name separator.
		if index := strings.LastIndex(file, "/"); index >= 0 {
			file = file[index+1:]
		} else if index = strings.LastIndex(file, "\\"); index >= 0 {
			file = file[index+1:]
		}
	} else {
		file = "???"
		line = 1
	}
	t.Fatalf("%s:%d: %s", file, line, message)
}

func assertNoError(t *testing.T, err error, message string) {
	if err != nil {
		assertFailed(t, fmt.Sprintf("%s: %v", message, err))
	}
}

func assertTrue(t *testing.T, success bool, message string) {
	if !success {
		assertFailed(t, message)
	}
}

func assertFalse(t *testing.T, failure bool, message string) {
	if failure {
		assertFailed(t, message)
	}
}