//
// Copyright 2022 SkyAPM org
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package idgen

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"
)

var (
	seededIDGen  = rand.New(rand.NewSource(time.Now().UnixNano()))
	seededIDLock sync.Mutex
)

func oldGenerateGlobalID() string {
	seededIDLock.Lock()
	seededID := seededIDGen.Int63()
	seededIDLock.Unlock()
	id := []int64{time.Now().UnixNano(), 0, seededID}
	ii := make([]string, len(id))
	for i, v := range id {
		ii[i] = fmt.Sprint(v)
	}
	return strings.Join(ii, ".")
}

// old method eg: 1586620117067704000.0.5762895797676669347
func BenchmarkOldGenerateGlobalID(t *testing.B) {
	t.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			oldGenerateGlobalID()
		}
	})
}

// UUID eg: e38bf4d27c0b11ea9217acde48001122
func BenchmarkGenerateGlobalID(t *testing.B) {
	t.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := GenerateGlobalID()
			if err != nil {
				t.Fail()
			}
		}
	})
}
