//
// Copyright 2021 SkyAPM org
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

package go2sky

import (
	"testing"
)

func TestConstSampler_IsSampled(t *testing.T) {
	sampler := NewConstSampler(true)
	operationName := "op"
	sampled := sampler.IsSampled(operationName)
	if sampled != true {
		t.Errorf("const sampler should be sampled")
	}
	samplerNegative := NewConstSampler(false)
	sampledNegative := samplerNegative.IsSampled(operationName)
	if sampledNegative != false {
		t.Errorf("const sampler should not be sampled")
	}
}

func TestRandomSampler_IsSampled(t *testing.T) {
	randomSampler := NewRandomSampler(0.5)
	operationName := "op"

	//just for test case
	randomSampler.threshold = 100
	sampled := randomSampler.IsSampled(operationName)
	if sampled != true {
		t.Errorf("random sampler should be sampled")
	}

	randomSampler.threshold = 0
	sampled = randomSampler.IsSampled(operationName)
	if sampled != false {
		t.Errorf("random sampler should not be sampled")
	}
}
