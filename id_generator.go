package go2sky

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

func generateGlobalID() []int64 {
	return []int64{
		time.Now().UnixNano(),
		0,
		generateID(),
	}
}
