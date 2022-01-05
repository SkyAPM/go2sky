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
	"testing"

	commonv3 "skywalking.apache.org/repo/goapi/collect/common/v3"
)

type TestAgentConfigChangeWatcher struct {
	currentValue string
	lastEvent    AgentConfigEventType
}

func (s *TestAgentConfigChangeWatcher) Key() string {
	return "test.key"
}

func (s *TestAgentConfigChangeWatcher) Notify(eventType AgentConfigEventType, newValue string) {
	s.currentValue = newValue
	s.lastEvent = eventType
}

func (s *TestAgentConfigChangeWatcher) Value() string {
	return s.currentValue
}

func TestHandleCommand(t *testing.T) {
	configDiscoveryService := NewConfigDiscoveryService()
	testWatcher := &TestAgentConfigChangeWatcher{}
	configDiscoveryService.BindWatchers([]AgentConfigChangeWatcher{
		testWatcher,
	})
	tests := []struct {
		name string
		// ready to handle command
		command *commonv3.Command
		// verify uuid
		uuid string
		// verify current value
		value string
		// verify last event
		lastEvent AgentConfigEventType
	}{
		{
			name: "first add",
			command: &commonv3.Command{Args: []*commonv3.KeyStringValuePair{
				{Key: "test.key", Value: "a"},
				{Key: "UUID", Value: "uuid1"},
			}},
			uuid:      "uuid1",
			value:     "a",
			lastEvent: MODIFY,
		},
		{
			name: "modify key",
			command: &commonv3.Command{Args: []*commonv3.KeyStringValuePair{
				{Key: "test.key", Value: "b"},
				{Key: "UUID", Value: "uuid2"},
			}},
			uuid:      "uuid2",
			value:     "b",
			lastEvent: MODIFY,
		},
		{
			name: "same uuid",
			command: &commonv3.Command{Args: []*commonv3.KeyStringValuePair{
				{Key: "test.key", Value: "b"},
				{Key: "UUID", Value: "uuid2"},
			}},
			uuid:      "uuid2",
			value:     "b",
			lastEvent: MODIFY,
		},
		{
			name: "delete key",
			command: &commonv3.Command{Args: []*commonv3.KeyStringValuePair{
				{Key: "UUID", Value: "uuid3"},
			}},
			uuid:      "uuid3",
			value:     "",
			lastEvent: DELETED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// handle command
			configDiscoveryService.HandleCommand(tt.command)

			// verify result
			if tt.value != testWatcher.currentValue {
				t.Errorf("error validate current value, current is: %s, excepted is: %s", testWatcher.currentValue, tt.value)
			}
			if tt.uuid != configDiscoveryService.UUID {
				t.Errorf("error validate current uuid, current is: %s, excepted is: %s", configDiscoveryService.UUID, tt.uuid)
			}
			if tt.lastEvent != testWatcher.lastEvent {
				t.Errorf("error validate current last event type, current is: %d, excepted is: %d", testWatcher.lastEvent, tt.lastEvent)
			}
		})
	}
}
