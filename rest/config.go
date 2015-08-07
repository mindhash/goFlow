package rest

import "fmt"

var config *ServerConfig

const {
	// Default value of ServerConfig.MaxIncomingConnections
	DefaultMaxIncomingConnections = 0
}

// JSON object that defines the server configuration.
type ServerConfig struct {
	Interface                      *string         // Interface to bind REST API to, default ":4984"
	SSLCert				   		   *string
	SSLKey						   *string
	ServerReadTimeout              *int            // maximum duration.Second before timing out read of the HTTP(S) request
	ServerWriteTimeout             *int            // maximum duration.Second before timing out write of the HTTP(S) response
	//Log                            []string        // Log keywords to enable
	//LogFilePath                    *string         // Path to log file, if missing write to stderr 
	Databases                      DbConfigMap     // Pre-configured databases, mapped by name
}

type DbConfig struct {
	Name               string
	Bucket             *string						//Bucket or Directory name for DB  
}

type DbConfigMap map[string]*DbConfig

func ParseCommandLine() {

	
	
	dbName := flag.String("dbName","flowDB","Default Database Name")
	addr   := flag.String("addr","localhost:4984","HTTP Server Address")
	dbBucket := flag.String("dbBucket","/Users/amolumbarkar/GoProjects")
	flag.Parse()
	
	config = &ServerConfig { Interface: &addr, Databases: map[string]*DbConfig{ dbName:{Name: dbName,Bucket:&dbBucket}}}
}

func (config *ServerConfig) serve(addr string, handler http.Handler) {
	maxConns := DefaultMaxIncomingConnections
	if config.MaxIncomingConnections != nil {
		maxConns = *config.MaxIncomingConnections
	}

	err := base.ListenAndServeHTTP(addr, maxConns, config.SSLCert, config.SSLKey, handler, config.ServerReadTimeout, config.ServerWriteTimeout)
	if err != nil {
		base.LogFatal("Failed to start HTTP server on %s: %v", addr, err)
	}
}

func RunServer(config *ServerConfig) {
	
	sc := NewServerContext(config)
	
	for _, dbConfig := range config.Databases {
			if _, err := sc.AddDatabaseFromConfig(dbConfig); err != nil {
				base.LogFatal("Error opening database: %v", err)
			}
	}
	
	base.Logf("Starting server on %s ...", *config.Interface)
	config.serve(*config.Interface, CreatePublicHandler(sc))
	
}



func ServerMain(){
	fmt.Println("Initiating Server..")
	ParseCommandLine() 
	RunServer(config)
}