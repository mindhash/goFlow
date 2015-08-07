

package samples


import ( 
	"log"
    "github.com/boltdb/bolt"
)


func loaded(){
	
	db, err := bolt.Open("test.db",0600,nil)
	
	if err != nil {
		log.Fatal(err)
	}
	
	tx, err := db.Begin(true)
	if err != nil {
	    log.Fatal(err) 
	}
	defer tx.Rollback()

	// Use the transaction...
	b, err := tx.CreateBucket([]byte("MyBucket"))
	if err != nil {
	    log.Fatal(err)
	}
	
	err = b.Put([]byte("answer"), []byte("42"))
	if err != nil {
	    log.Fatal(err)
	}
	
    v := b.Get([]byte("answer"))
    fmt.Printf("The answer is: %s\n", v)
	   

	// Commit the transaction and check for error.
	if err := tx.Commit(); err != nil {
	    log.Fatal(err)
	}
	
	defer db.Close()
	fmt.Println("Db Open/Close ")
	
	
}


