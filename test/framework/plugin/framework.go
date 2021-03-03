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

package plugin

import (
	"context"
	"testing"
	"time"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter"
)

// NewService create a service that generates the trace data
type NewService func(ctx context.Context, tracer *go2sky.Tracer) error

type serviceHolder struct {
	service        string
	newServiceFunc NewService
}

func (h serviceHolder) run(ctx context.Context, oapAddr string) error {
	r, err := reporter.NewGRPCReporter(oapAddr)
	if err != nil {
		return err
	}
	defer r.Close()
	tracer, err := go2sky.NewTracer(h.service, go2sky.WithReporter(r))
	if err != nil {
		return err
	}
	return h.newServiceFunc(ctx, tracer)
}

type terminateFunc func()

// TestPlugin
type TestPlugin struct {
	expectedDataFile string
	services         []*serviceHolder
	terminateFuncs   []terminateFunc
	grpcServerAddr   string
	httpServerAddr   string
}

// CreateTestPlugin
func CreateTestPlugin() *TestPlugin {
	return &TestPlugin{}
}

// WithExpectedDataFile
func (p *TestPlugin) WithExpectedDataFile(filepath string) *TestPlugin {
	p.expectedDataFile = filepath
	return p
}

// AddService
func (p *TestPlugin) AddService(service string, serviceFunc NewService) *TestPlugin {
	p.services = append(p.services, &serviceHolder{
		service:        service,
		newServiceFunc: serviceFunc,
	})
	return p
}

func (p *TestPlugin) addTerminateFunc(terminateFunc terminateFunc) *TestPlugin {
	p.terminateFuncs = append(p.terminateFuncs, terminateFunc)
	return p
}

func (p *TestPlugin) verification(t *testing.T) {
	if len(p.services) < 1 {
		t.Fatal("at least one service needs to be added to generate the trace data")
	}

	if p.expectedDataFile == "" {
		p.expectedDataFile = "expected.data.yml"
	}
}

func (p *TestPlugin) runServices(ctx context.Context, t *testing.T) {
	for _, holder := range p.services {
		holder := holder
		go func() {
			e := holder.run(ctx, p.grpcServerAddr)
			if e != nil {
				t.Error(e)
			}
		}()
	}
	time.Sleep(time.Second * 10)
}

func (p *TestPlugin) Run(t *testing.T) {
	// step1. verify the test plugin configuration
	p.verification(t)

	ctx, cancel := context.WithCancel(context.Background())
	// step2. create mock collector container
	p.runMockCollector(ctx, t)

	// step3. run services
	p.runServices(ctx, t)

	// step4. validate expected data
	p.validateExpectedData(ctx, t)

	defer cancel()
	defer func() {
		for _, terminateFunc := range p.terminateFuncs {
			terminateFunc()
		}
	}()
}
