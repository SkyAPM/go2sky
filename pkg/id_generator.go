package pkg

import (
	"math/rand"
	"sync"
	"time"
)

var (
	seededIDGen  = rand.New(rand.NewSource(time.Now().UnixNano()))
	seededIDLock sync.Mutex
)

func generateID() int64 {
	seededIDLock.Lock()
	defer seededIDLock.Unlock()
	return seededIDGen.Int63()
}

// GenerateGlobalID generates global unique id
func GenerateGlobalID() []int64 {
	return []int64{
		time.Now().UnixNano(),
		0,
		generateID(),
	}
}

// GenerateScopedGlobalID generates global unique id with a scopeId prefix
func GenerateScopedGlobalID(scopeID int64) []int64 {
	return []int64{
		scopeID,
		time.Now().UnixNano(),
		generateID(),
	}
}
