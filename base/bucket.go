package base

// Filename: bucket.go
// Description: performs operations on either bucket or DB directory of Database
// if using Tiedot DB, bucket represents data directory


import( 
	"encoding/json" 
	"time"
	"log"
	"fmt"
	"github.com/boltdb/bolt"
)


type DbHandleInterface interface{
		GetName() string
		GetRaw(k string) (rv []byte, err error)
		PutRaw(k string, v []byte) (added bool, err error)
		//Delete(k string) error
		Close() (err error)
}

type DbHandle DbHandleInterface


// Full specification of how to connect to a bucket
type DatabaseSpec struct {
	Directory string
	BucketName string
	Name string 
}

type BoltDbHandle struct {
	*bolt.DB   
	Bucket string
	spec DatabaseSpec // keep a copy of the BucketSpec for DCP usage
}

func (dbhandle BoltDbHandle) GetName() string {
	return dbhandle.spec.Name		// return Database name from spec
}

func (dbhandle BoltDbHandle) PutRaw(k string, v []byte) (added bool, err error){
	tx, err := dbhandle.DB.Begin(true)
	if err != nil {
	    return false, err
	}
	defer tx.Rollback()

	// Use the transaction...
	b := tx.Bucket([]byte(dbhandle.Bucket))			// need to check error handling here
	
	//if err != nil {
	//    return false, err
	//}

	err = b.Put([]byte(k), []byte(v)) 
	if err != nil {
	    return false, err
	}
	  
	// Commit the transaction and check for error.
	if err = tx.Commit(); err != nil {
	    return false, err
	}
	return true, nil
}

// get  value for a key
// needs work 
func (dbhandle BoltDbHandle) Get(k string) (v interface{},err error){
	data,err := dbhandle.GetRaw(k)
	if err != nil {
	    return nil,err
	}
	
	err = json.Unmarshal([]byte(data), v)
	if err != nil {
		Warn("Error unmarshaling body of doc %q: %s", k, err)
		return nil, err
	}
	return v,nil
}

// get json value for a key
func (dbhandle BoltDbHandle) GetRaw(k string) (v []byte, err error) {
	//start bolt transaction
	tx, err := dbhandle.DB.Begin(true)
	if err != nil {
	    return nil, err
	}
	defer tx.Rollback()
	
	// Use the transaction...
	b := tx.Bucket([]byte(dbhandle.Bucket))
	if err != nil {
	    return  nil, err
	}
	
	v = b.Get([]byte(k))
	
	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
	    return  nil, err
	}
	return v, nil
}

func (dbhandle BoltDbHandle) close() error {
	err:=dbhandle.close()  // TO DO: need to handle error here
    if err != nil {
        log.Fatal(err)
    }
	return nil
}


// Creates a Bucket that talks to boltDB
func GetBoltDbHandle(spec DatabaseSpec) (dbhandle DbHandle, err error) {
	
	//open database 
	dbconn, err := bolt.Open(spec.Directory, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
        log.Fatal(err)
	}
	
	// Start a writable transaction.
	tx, err:= dbconn.Begin(true) 
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
			
	_, err = tx.CreateBucketIfNotExists([]byte(spec.BucketName))
	if err != nil {
	    return nil, fmt.Errorf("create bucket: %s", err)
	}

	// Commit the transaction and check for error.
	 
	if err := tx.Commit(); err != nil {
		return nil,err
	}
	 
	dbhandle = BoltDbHandle{dbconn, spec.BucketName,spec}

	return
}

//get handle to database 
func GetDbHandle(spec DatabaseSpec) (dbhandle DbHandle, err error) {
	//TO DO: validate spec
	
	dbhandle, err = GetBoltDbHandle(spec)
	if err != nil {
			panic(err)
	}
	return
}