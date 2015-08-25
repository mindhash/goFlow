package rest

import ("testing" 
"fmt"
		"github.com/mindhash/go.assert"  
)



func TestDBAccess(t *testing.T) {

	sc := NewServerContext(&ServerConfig{})
	bname:= "TestBucket"
	dname:="db"
	dbConfig := &DbConfig{
		Bucket: &bname ,
		Name: dname,
	}
	
	_, err := sc.AddDatabaseFromConfig(dbConfig)
	if err != nil {
		panic(fmt.Sprintf("Error from AddDatabaseFromConfig: %v", err))
	}
	dbContext,err := sc.GetDatabase()
	
	assert.True(t, (dbContext != nil))
	
	
}
