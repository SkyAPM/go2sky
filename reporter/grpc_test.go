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

package reporter

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter/grpc/common"
	"github.com/SkyAPM/go2sky/reporter/grpc/register"
	"github.com/SkyAPM/go2sky/reporter/grpc/register/mock_register"
)

const header string = "1-MTU1NTY0NDg4Mjk2Nzg2ODAwMC4wLjU5NDYzNzUyMDYzMzg3NDkwODc=" +
	"-NS4xNTU1NjQ0ODgyOTY3ODg5MDAwLjM3NzUyMjE1NzQ0Nzk0NjM3NTg=" +
	"-1-2-3-I2NvbS5oZWxsby5IZWxsb1dvcmxk-Iy9yZXN0L2Fh-Iy9nYXRld2F5L2Nj"

func Test_gRPCReporter_Register(t *testing.T) {
	serviceName, serviceID, instanceName, instanceID, reporter := createMockReporter(t)
	aServiceID, aInstanceID, err := reporter.Register(serviceName, instanceName)
	if err != nil || serviceID != aServiceID || instanceID != aInstanceID {
		t.Errorf("register service and instance error")
	}
}

func Test_e2e(t *testing.T) {
	serviceName, _, instance, _, reporter := createMockReporter(t)
	reporter.sendCh = make(chan *common.UpstreamSegment, 10)
	tracer, err := go2sky.NewTracer(serviceName, go2sky.WithReporter(reporter), go2sky.WithInstance(instance))
	if err != nil {
		t.Error(err)
	}
	tracer.WaitUntilRegister()
	entrySpan, ctx, err := tracer.CreateEntrySpan(context.Background(), "/rest/api", func() (string, error) {
		return header, nil
	})
	if err != nil {
		t.Error(err)
	}
	exitSpan, err := tracer.CreateExitSpan(ctx, "/foo/bar", "foo.svc:8787", func(head string) error {
		if head == "" {
			t.Fail()
		}
		if len(strings.Split(head, "-")) != 9 {
			t.Fail()
		}
		return nil
	})
	if err != nil {
		t.Error(err)
	}
	exitSpan.End()
	entrySpan.End()
	for s := range reporter.sendCh {
		reporter.Close()
		if len(s.GlobalTraceIds) != 1 && len(s.GlobalTraceIds[0].IdParts) != 3 {
			t.Error("trace id format is incorrect")
		}
		if s.Segment == nil {
			t.Error("null segment")
		}
	}
}

func createMockReporter(t *testing.T) (string, int32, string, int32, *gRPCReporter) {
	ctrl := gomock.NewController(t)
	mockRegisterClient := mock_register.NewMockRegisterClient(ctrl)
	reporter := &gRPCReporter{
		registerClient: mockRegisterClient,
		logger:         log.New(os.Stderr, "go2sky", log.LstdFlags),
	}

	serviceID := rand.Int31()
	serviceName := fmt.Sprintf("service-%d", serviceID)
	mockRegisterClient.EXPECT().DoServiceRegister(
		gomock.Any(),
		gomock.Any(),
	).Return(&register.ServiceRegisterMapping{Services: []*common.KeyIntValuePair{{
		Value: serviceID,
		Key:   serviceName,
	}}}, nil)
	instanceID := rand.Int31()
	instanceName := fmt.Sprintf("instance-%d", instanceID)
	mockRegisterClient.EXPECT().DoServiceInstanceRegister(
		gomock.Any(),
		gomock.Any(),
	).Return(&register.ServiceInstanceRegisterMapping{ServiceInstances: []*common.KeyIntValuePair{{
		Value: instanceID,
		Key:   instanceName,
	}}}, nil)
	return serviceName, serviceID, instanceName, instanceID, reporter
}
