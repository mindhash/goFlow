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
//	databases_  map[string]*db.DatabaseContext
//	lock        sync.RWMutex
	statsTicker *time.Ticker
	HTTPClient  *http.Client
}

//get existing or add database from config
func (sc *ServerContext) getOrAddDatabaseFromConfig(config *DbConfig, useExisting bool) (*db.DatabaseContext, error) {
	return sc
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