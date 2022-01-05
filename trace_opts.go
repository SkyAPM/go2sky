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

package go2sky

// WithReporter setup report pipeline for tracer
func WithReporter(reporter Reporter) TracerOption {
	return func(t *Tracer) {
		t.reporter = reporter
	}
}

// WithInstance setup instance identify
func WithInstance(instance string) TracerOption {
	return func(t *Tracer) {
		t.instance = instance
	}
}

// WithSampler setup sampler
func WithSampler(samplingRate float64) TracerOption {
	return func(t *Tracer) {
		t.sampler = NewDynamicSampler(samplingRate, t)
	}
}

// WithCustomSampler setup custom sampler
func WithCustomSampler(sampler Sampler) TracerOption {
	return func(t *Tracer) {
		t.sampler = sampler
	}
}
