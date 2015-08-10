package rest

import (
	//"bytes"
	"time"
	"net/http"
	//"net/url"
	"github.com/mindhash/goFlow/db"
)

type ServerContext struct {
	config      *ServerConfig
	databases_  map[string]*db.DatabaseContext
	lock        sync.RWMutex
//	statsTicker *time.Ticker
	HTTPClient  *http.Client
}

//get existing or add database from config
func (sc *ServerContext) getOrAddDatabaseFromConfig(config *DbConfig, useExisting bool) (*db.DatabaseContext, error) {
	return nil,nil
}


// Adds a database to the ServerContext given its configuration.  If an existing config is found
// for the name, returns an error.
func (sc *ServerContext) AddDatabaseFromConfig(config *DbConfig) (*db.DatabaseContext, error) {
	return sc.getOrAddDatabaseFromConfig(config, false)
}

func NewServerContext(config *ServerConfig) *ServerContext {
	sc := &ServerContext{
		config:     config, 
		HTTPClient: http.DefaultClient,
	}
	return sc
}


func (sc *ServerContext) GetDatabase(name string) (*db.DatabaseContext, error) {
//	sc.lock.RLock()
//	dbc := sc.databases_[name]
//	sc.lock.RUnlock()
//	if dbc != nil {
//		return dbc, nil
//	} else if db.ValidateDatabaseName(name) != nil {
//		return nil, base.HTTPErrorf(http.StatusBadRequest, "invalid database name %q", name)
//	} else if sc.config.ConfigServer == nil {
//		return nil, base.HTTPErrorf(http.StatusNotFound, "no such database %q", name)
//	} else {
//		// Let's ask the config server if it knows this database:
//		base.Logf("Asking config server %q about db %q...", *sc.config.ConfigServer, name)
//		config, err := sc.getDbConfigFromServer(name)
//		if err != nil {
//			return nil, err
//		}
//		if dbc, err = sc.getOrAddDatabaseFromConfig(config, true); err != nil {
//			return nil, err
//		}
//	}
//	return dbc, nil
return nil, nil
}