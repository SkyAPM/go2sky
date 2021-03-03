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
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
)

func (p *TestPlugin) validateExpectedData(ctx context.Context, t *testing.T) {
	client := http.Client{}

	expectedDataFile, err := os.Open(p.expectedDataFile)
	if err != nil {
		t.Error(err)
	}
	defer expectedDataFile.Close()
	expectedData, err := ioutil.ReadAll(expectedDataFile)
	if err != nil {
		t.Error(err)
	}

	dataValidateReq, err := http.NewRequest("POST", fmt.Sprintf("http://%s/dataValidate", p.httpServerAddr), bytes.NewReader(expectedData))
	if err != nil {
		t.Error(err)
	}

	dataValidateResp, err := client.Do(dataValidateReq)
	if err != nil {
		t.Error(err)
	}
	defer dataValidateResp.Body.Close()
	if dataValidateResp.StatusCode == http.StatusOK {
		return
	}

	// validate failed, print actual data
	actualDataReq, err := http.NewRequest("POST", fmt.Sprintf("http://%s/receiveData", p.httpServerAddr), nil)
	if err != nil {
		t.Error(err)
	}
	actualDataResp, err := client.Do(actualDataReq)
	if err != nil {
		t.Error(err)
	}
	defer actualDataResp.Body.Close()
	actualData, err := ioutil.ReadAll(actualDataResp.Body)
	if err != nil {
		t.Error(err)
	}
	t.Errorf("\nexpectedData:\n%s\nactualData:\n%s\n", expectedData, actualData)
}
