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

package reporter

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	swAgentAuthentication                         = "SW_AGENT_AUTHENTICATION"
	swAgentLayer                                  = "SW_AGENT_LAYER"
	swAgentCollectorHeartbeatPeriod               = "SW_AGENT_COLLECTOR_HEARTBEAT_PERIOD"
	swAgentCollectorGetAgentDynamicConfigInterval = "SW_AGENT_COLLECTOR_GET_AGENT_DYNAMIC_CONFIG_INTERVAL"
	swAgentCollectorBackendServices               = "SW_AGENT_COLLECTOR_BACKEND_SERVICES"
	swAgentCollectorMaxSendQueueSize              = "SW_AGENT_COLLECTOR_MAX_SEND_QUEUE_SIZE"
	swAgentProcessStatusHookEnable                = "SW_AGENT_PROCESS_STATUS_HOOK_ENABLE"
	swAgentProcessLabels                          = "SW_AGENT_PROCESS_LABELS"
)

// serverAddrFormEnv read the backend service address in the environment variable
func serverAddrFormEnv(serverAddr string) string {
	if value := os.Getenv(swAgentCollectorBackendServices); value != "" {
		return value
	}
	return serverAddr
}

// gRPCReporterOptionsFormEnv read the options in the environment variable
func gRPCReporterOptionsFormEnv() (opts []GRPCReporterOption, err error) {
	if auth := os.Getenv(swAgentAuthentication); auth != "" {
		opts = append(opts, WithAuthentication(auth))
	}

	if layer := os.Getenv(swAgentLayer); layer != "" {
		opts = append(opts, WithLayer(layer))
	}

	if value := os.Getenv(swAgentCollectorHeartbeatPeriod); value != "" {
		period, err1 := strconv.ParseInt(value, 0, 64)
		if err1 != nil {
			return nil, errors.Wrap(err1, fmt.Sprintf("%s=%s", swAgentCollectorHeartbeatPeriod, value))
		}
		opts = append(opts, WithCheckInterval(time.Duration(period)*time.Second))
	}

	if value := os.Getenv(swAgentCollectorGetAgentDynamicConfigInterval); value != "" {
		interval, err1 := strconv.ParseInt(value, 0, 64)
		if err1 != nil {
			return nil, errors.Wrap(err1, fmt.Sprintf("%s=%s", swAgentCollectorGetAgentDynamicConfigInterval, value))
		}
		opts = append(opts, WithCDS(time.Duration(interval)*time.Second))
	}

	if value := os.Getenv(swAgentCollectorMaxSendQueueSize); value != "" {
		size, err1 := strconv.ParseInt(value, 0, 64)
		if err1 != nil {
			return nil, err
		}
		opts = append(opts, WithMaxSendQueueSize(int(size)))
	}

	if value := os.Getenv(swAgentProcessStatusHookEnable); value != "" {
		enable, err1 := strconv.ParseBool(value)
		if err1 != nil {
			return nil, err
		}
		opts = append(opts, WithProcessStatusHook(enable))
	}

	if value := os.Getenv(swAgentProcessLabels); value != "" {
		labels := strings.Split(value, ",")
		opts = append(opts, WithProcessLabels(labels))
	}
	return
}
