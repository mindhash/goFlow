package base

// Filename: bucket.go
// Description: performs operations on either bucket or DB directory of Database
// if using Tiedot DB, bucket represents data directory


import(
	"github.com/boltdb/bolt"
)


type DbHandleInterface interface{
		GetName() string
		GetRaw(k string) (rv []byte, err error)
		PutRaw(k string, v []byte) (added bool, err error)
		//Delete(k string) error
		Close()
}

type dbhandle DbHandleInterface


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
	return dbhandle.Name		// probably need to use spec
}

func (dbhandle BoltDbHandle) PutRaw(k string, v []byte) (added bool, err error){
	tx, err := dbhandle.DB.Begin(true)
	if err != nil {
	    return err
	}
	defer tx.Rollback()

	// Use the transaction...
	_, err := tx.Bucket([]byte(dbhandle.Bucket))
	if err != nil {
	    return err
	}

	err := b.Put([]byte(k), []byte(v))
	
	if err != nil {
	    return err
	}
	  
	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
	    return err
	}
	return 
}



func (dbhandle BoltDbHandle) GetRaw(k string) (v []byte, err error) {
	//start bolt transaction
	tx, err := dbhandle.DB.Begin(true)
	if err != nil {
	    return err
	}
	defer tx.Rollback()
	
	// Use the transaction...
	_, err := tx.Bucket([]byte(dbhandle.Bucket))
	if err != nil {
	    return err
	}
	
	v := b.Get([]byte(k))
	
	fmt.Printf("The Key Value is: %s\n", v)
	
	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
	    return err
	}
	return 
}

func (dbhandle BoltDbHandle) func close(){
	dbhandle.DB.close()
    if err != nil {
        log.Fatal(err)
    }
}


// Creates a Bucket that talks to boltDB
func GetBoltDbHandle(spec DatabaseSpec) (dbhandle DbHandle, err error) {
	
	//open database 
	dbconn, err := bolt.Open(spec.Directory, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
        log.Fatal(err)
	}
	
	// Start a writable transaction.
	tx, err := dbconn.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
			
	_, err := tx.CreateBucketIfNotExists([]byte(spec.BucketName))
	if err != nil {
	    return fmt.Errorf("create bucket: %s", err)
	}

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
		return err
	}
	 
	dbhandle = BoltDbHandle{dbconn, spec.BucketName,spec}

	return
}

//get handle to database 
func GetDbHandle(spec DatabaseSpec) (dbhandle DbHandle, err error) {
	//TO DO: validate spec
	
	dbhandle, err := GetBoltDbHandle(spec)
	if err != nil {
			panic(err)
	}
	return
}