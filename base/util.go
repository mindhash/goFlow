package base
 

import (
	"strconv"
	"crypto/rand" 
	"fmt"
	"io"  
	"sync"
)

//used for password generation
func GenerateRandomSecret() string {
	randomBytes := make([]byte, 20)
	n, err := io.ReadFull(rand.Reader, randomBytes)
	if n < len(randomBytes) || err != nil {
		panic("RNG failed, can't create password")
	}
	return fmt.Sprintf("%x", randomBytes)
}

// Returns a cryptographically-random 160-bit number encoded as a hex string.
func CreateUUID() string {
	bytes := make([]byte, 16)
	n, err := rand.Read(bytes)
	if n < 16 {
		LogPanic("Failed to generate random ID: %s", err)
	}
	return fmt.Sprintf("%x", bytes)
}


// IntMax is an expvar.Value that tracks the maximum value it's given.
// used in various places including http listner
type IntMax struct {
	i  int64
	mu sync.RWMutex
}

func (v *IntMax) String() string {
	v.mu.RLock()
	defer v.mu.RUnlock()
	return strconv.FormatInt(v.i, 10)
}

func FloatToString(input_num float64) string {
    // to convert a float number to a string
    return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func (v *IntMax) SetIfMax(value int64) {
	v.mu.Lock()
	defer v.mu.Unlock()
	if value > v.i {
		v.i = value
	}
	
}