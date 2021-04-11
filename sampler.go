// Licensed to SkyAPM org under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. SkyAPM org licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package go2sky

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

type Sampler interface {
	IsSampled(operation string) (sampled bool)
}

type ConstSampler struct {
	decision bool
}

// NewConstSampler creates a ConstSampler.
func NewConstSampler(sample bool) *ConstSampler {
	s := &ConstSampler{
		decision: sample,
	}
	return s
}

// IsSampled implements IsSampled() of Sampler.
func (s *ConstSampler) IsSampled(operation string) bool {
	return s.decision
}

type RandomSampler struct {
	samplingRate float64
	rand         *rand.Rand
	threshold    int
}

// IsSampled implements IsSampled() of Sampler.
func (s *RandomSampler) IsSampled(operation string) bool {
	return s.threshold >= s.rand.Intn(100)
}

func (s *RandomSampler) init() {
	s.rand = rand.New(rand.NewSource(time.Now().Unix()))
	s.threshold = int(s.samplingRate * 100)
}

func NewRandomSampler(samplingRate float64) *RandomSampler {
	s := &RandomSampler{
		samplingRate: samplingRate,
	}
	s.init()
	return s
}

type DynamicSampler struct {
	currentRate float64
	defaultRate float64
	sampler     Sampler
}

// IsSampled implements IsSampled() of Sampler.
func (s *DynamicSampler) IsSampled(operation string) bool {
	return s.sampler.IsSampled(operation)
}

func (s *DynamicSampler) Key() string {
	return "agent.sample_rate"
}

func (s *DynamicSampler) Notify(eventType AgentConfigEventType, newValue string) {
	if eventType == DELETED {
		newValue = fmt.Sprintf("%f", s.defaultRate)
	}
	samplingRate, err := strconv.ParseFloat(newValue, 64)
	if err != nil {
		return
	}

	// change sampler
	var sampler Sampler
	if samplingRate <= 0 {
		sampler = NewConstSampler(false)
	} else if samplingRate >= 1.0 {
		sampler = NewConstSampler(true)
	} else {
		sampler = NewRandomSampler(samplingRate)
	}
	s.sampler = sampler
	s.currentRate = samplingRate
}

func (s *DynamicSampler) Value() string {
	return fmt.Sprintf("%f", s.currentRate)
}

func NewDynamicSampler(samplingRate float64, tracer *Tracer) *DynamicSampler {
	s := &DynamicSampler{
		currentRate: samplingRate,
		defaultRate: samplingRate,
	}
	s.Notify(MODIFY, fmt.Sprintf("%f", samplingRate))
	// append watcher
	tracer.dcsWatchers = append(tracer.dcsWatchers, s)
	return s
}
