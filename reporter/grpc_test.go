// Copyright 2019 Tetrate Labs
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package reporter

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/tetratelabs/go2sky/reporter/grpc/common"
	"github.com/tetratelabs/go2sky/reporter/grpc/register"
	"github.com/tetratelabs/go2sky/reporter/grpc/register/mock_register"
)

func Test_gRPCReporter_Register(t *testing.T) {
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
	aServiceID, aInstanceID, err := reporter.Register(serviceName, instanceName)
	if err != nil || serviceID != aServiceID || instanceID != aInstanceID {
		t.Errorf("register service and instance error")
	}
}
