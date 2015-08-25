package db

import (
	"time"
	"net/http"
	"regexp"
	"expvar"
	"github.com/mindhash/goFlow/auth"
	"github.com/mindhash/goFlow/base"
)

// Basic description of a database. Shared between all Database objects on the same database.
// This object is thread-safe so it can be shared between HTTP handlers.
type DatabaseContext struct {
	Name               string                  // Database name
	DbHandle             base.DbHandle             // Storage
	StartTime          time.Time               // Timestamp when context was instantiated
}


// Represents a database. A new instance is created for each HTTP request,
// so this struct does not have to be thread-safe.
type Database struct {
	*DatabaseContext
	user auth.User 
}


// All special/internal documents have this prefix in their keys.
const kWfKeyPrefix = "_goflow:"

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

// open new database connection handle
func OpenDatabase(spec base.DatabaseSpec) (dbhandle base.DbHandle, err error) {
	
	dbhandle, err = base.GetDbHandle(spec)
	
	if err != nil {
		err = base.HTTPErrorf(http.StatusBadGateway,
			"Unable to connect to server: %s", err)
	}  
	
	return
}

// Creates a new DatabaseContext on a database. The DB and bucket will be closed when this context closes.
func NewDatabaseContext(dbName string, dbhandle base.DbHandle) (*DatabaseContext, error) {
	
	if err := ValidateDatabaseName(dbName); err != nil {
		return nil, err
	}
	
	context := &DatabaseContext{
		Name:       dbName,
		DbHandle:     dbhandle,
		StartTime:  time.Now(),
	}
	return context, nil
}

// close database. important to close connection with DB to avoid locks
func (context *DatabaseContext) Close() {
	context.DbHandle.Close()
	context.DbHandle = nil
}

// check if DB is closed
func (context *DatabaseContext) IsClosed() bool {
	return context.DbHandle == nil
}


// Makes a database object given context and user.
func GetDatabase(context *DatabaseContext, user auth.User) (*Database, error) {
	return &Database{context, user}, nil
}

// create new database object
func CreateDatabase(context *DatabaseContext) (*Database, error) {
	return &Database{context, nil}, nil
}

func (db *Database) SameAs(otherdb *Database) bool {
	return db != nil && otherdb != nil &&
		db.DbHandle == otherdb.DbHandle
}

