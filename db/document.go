package db

import 	(
		"encoding/json"
		"github.com/mindhash/goFlow/base"
	 	)

type Body map[string]interface{}

type document struct { 
	body Body
	ID   string `json:"-"`
}

// Returns a new empty document.
func newDocument(docid string) *document {
	return &document{ID: docid}
}

// Unmarshals a document from JSON data. 
func unmarshalDocument(docid string, data []byte) (*document, error) {
	doc := newDocument(docid);
	if len(data) > 0 {
		if err := json.Unmarshal(data, doc); err != nil {
			return nil, err
		}
	}
	return doc, nil
}


func (doc *document) UnmarshalJSON(data []byte) error {
	//need to identify doc first
	if doc.ID == "" {
		base.Warn("Doc was unmarshaled without ID set") // TO DO: may be panic not required
	}
	
	//fetch json into doc.body
	err := json.Unmarshal([]byte(data), &doc.body)
	if err != nil {
		base.Warn("Error unmarshaling body of doc %q: %s", doc.ID, err)
		return err
	}

	return nil
}

func (doc *document) MarshalJSON() ([]byte, error) {
	body := doc.body
	if body == nil {
		body = Body{}
	}
	
	data, err := json.Marshal(body)
	
	return data, err
}
