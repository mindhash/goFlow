package db


// Basic description of a database. Shared between all Database objects on the same database.
// This object is thread-safe so it can be shared between HTTP handlers.
type DatabaseContext struct {
	Name               string                  // Database name
	Bucket             string             // Storage
}


// Represents a simulated CouchDB database. A new instance is created for each HTTP request,
// so this struct does not have to be thread-safe.
type Database struct {
	*DatabaseContext
	user string
}


type Body map[string]interface{}


// Makes a Database object given its name and bucket.
func GetDatabase(context *DatabaseContext, user string) (*Database, error) {
	return &Database{context, user}, nil
}
