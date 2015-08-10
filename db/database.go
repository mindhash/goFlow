package db


// Basic description of a database. Shared between all Database objects on the same database.
// This object is thread-safe so it can be shared between HTTP handlers.
type DatabaseContext struct {
	Name               string                  // Database name
	DbHandle             base.DbHandle             // Storage
	StartTime          time.Time               // Timestamp when context was instantiated
}


// Represents a simulated CouchDB database. A new instance is created for each HTTP request,
// so this struct does not have to be thread-safe.
type Database struct {
	*DatabaseContext
	user string 
}

type Body map[string]interface{}

// All special/internal documents the gateway creates have this prefix in their keys.
const kWfKeyPrefix = "_sync:")

var dbExpvars = expvar.NewMap("goflow_db")

//general validation of DB name. 
// can be skipped / changed according to database 
func ValidateDatabaseName(dbName string) error {
	if match, _ := regexp.MatchString(`^[a-z][-a-z0-9_$()+/]*$`, dbName); !match {
		return base.HTTPErrorf(http.StatusBadRequest,
			"Illegal database name: %s", dbName)
	}
	return nil
}

func openDatabase(spec base.DatabaseSpec) (dbhandle base.DbHandle, err error) {
	
	dbhandle, err = base.GetDbHandle(spec)
	
	if err != nil {
		err = base.HTTPErrorf(http.StatusBadGateway,
			"Unable to connect to server: %s", err)
	}  
	
	return
}



// Makes a Database object given its name and bucket.
func GetDatabase(context *DatabaseContext, user string) (*Database, error) {
	return &Database{context, user}, nil
}


