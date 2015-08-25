package rest

import (
	//"bytes"
	//"time"
	"net/http"
	//"net/url"
	"sync"
	"github.com/mindhash/goFlow/db"
	"github.com/mindhash/goFlow/base"
)

type ServerContext struct {
	config      *ServerConfig
	database_  *db.DatabaseContext// databases_ map[string]*db.DatabaseContext
	lock        sync.RWMutex
//	statsTicker *time.Ticker
	HTTPClient  *http.Client
}

//get existing or add database from config
func (sc *ServerContext) getOrAddDatabaseFromConfig(config *DbConfig) (*db.DatabaseContext, error) {
	// Obtain write lock during add database, to avoid race condition when creating based on ConfigServer
	sc.lock.Lock()
	defer sc.lock.Unlock()

	//return existing DB if opened already
	if sc.database_ != nil {
		 return sc.database_, nil
	}
	
		
	//derive values from config
	bucketName := config.Name		// Default bucket
	dbName := config.Name			// Name used to identify DB
	dbDirectory := config.Name + ".db" //directory + DB file name
	
	if config.Bucket != nil {
		bucketName = *config.Bucket
	}
	
	if dbName == "" {
		dbName = bucketName
	}
	
	base.Logf("Opening db /%s as bucket %q ",
		dbName, bucketName)

	if err := db.ValidateDatabaseName(dbName); err != nil {
		return nil, err
	}
	 
	spec := base.DatabaseSpec {Directory: dbDirectory,
	BucketName: bucketName,
	Name: dbName }
	
	dbHandle,err := db.OpenDatabase (spec)
	if err != nil {
		return nil, err
	}
	
	dbcontext, err := db.NewDatabaseContext(dbName, dbHandle)
	if err != nil {
		return nil, err
	}
	
	// Register it so HTTP handlers can find it:
	sc.database_ = dbcontext

	// Save the config
	sc.config.Database = config
	
	return dbcontext, nil
}


// Adds a database to the ServerContext given its configuration.  If an existing config is found
// for the name, returns an error.
func (sc *ServerContext) AddDatabaseFromConfig(config *DbConfig) (*db.DatabaseContext, error) {
	return sc.getOrAddDatabaseFromConfig(config)
}

func NewServerContext(config *ServerConfig) *ServerContext {
	sc := &ServerContext{
		config:     config, 
		HTTPClient: http.DefaultClient,
	}
	return sc
}

// 
func (sc *ServerContext) GetDatabase() (*db.DatabaseContext, error) {
	sc.lock.RLock()
	dbc := sc.database_
	sc.lock.RUnlock()
	
	if dbc != nil {
		return dbc, nil
	} else {
		return nil, base.HTTPErrorf(http.StatusBadRequest, "database name is invalid or it is not open ")
	}
	
	return dbc, nil
}

func (sc *ServerContext) CloseDatabase() {
	
}