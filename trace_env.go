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

import (
	"fmt"
	"os"
	"strconv"

	"github.com/pkg/errors"
)

const (
	swAgentName         = "SW_AGENT_NAME"
	swAgentInstanceName = "SW_AGENT_INSTANCE_NAME"
	swAgentSample       = "SW_AGENT_SAMPLE"
)

// serviceFormEnv read the service in the environment variable
func serviceFormEnv(service string) string {
	if value := os.Getenv(swAgentName); value != "" {
		return value
	}
	return service
}

// traceOptionsFormEnv read the options in the environment variable
func traceOptionsFormEnv() (opts []TracerOption, err error) {
	// SW_AGENT_INSTANCE_NAME
	if instance := os.Getenv(swAgentInstanceName); instance != "" {
		opts = append(opts, WithInstance(instance))
	}

	// SW_AGENT_SAMPLE
	if value := os.Getenv(swAgentSample); value != "" {
		samplingRate, err1 := strconv.ParseFloat(value, 64)
		if err1 != nil {
			return nil, errors.Wrap(err1, fmt.Sprintf("%s=%s", swAgentSample, value))
		}
		opts = append(opts, WithSampler(samplingRate))
	}
	return
}
