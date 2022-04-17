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
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	commonv3 "skywalking.apache.org/repo/goapi/collect/common/v3"
)

type ProcessReportStatus int8

const (
	ProcessLabelKey = "processLabels"

	NotInit ProcessReportStatus = iota
	Reported
	Closed
)

var process *processStat

type processStat struct {
	basePath     string
	metadataFile string
	status       ProcessReportStatus
	shutdownOnce sync.Once
}

func initProcessStat(r *gRPCReporter) *processStat {
	basePath := path.Join(os.TempDir(), "apache_skywalking", "process", strconv.Itoa(os.Getpid()))
	metaFilePath := path.Join(basePath, "metadata.properties")

	return &processStat{
		basePath:     basePath,
		metadataFile: metaFilePath,
		status:       NotInit,
	}
}

// Report the current process metadata to local file
// using to work with eBPF agent
func reportProcess(r *gRPCReporter) {
	if process == nil {
		process = initProcessStat(r)
	}

	if process.status == NotInit {
		// create metadata file
		if p, err := process.initMetadataFile(r); err != nil {
			r.logger.Warnf("process status file init failure: %s, %v", p, err)
		} else {
			process.status = Reported
		}
	} else if process.status == Reported {
		// keep the metadata file alive(update modify time)
		updateTime := time.Now()
		if err := os.Chtimes(process.metadataFile, updateTime, updateTime); err != nil {
			r.logger.Warnf("keep the process metadata alive failure: %v", err)
		}
	}
}

func (p *processStat) initMetadataFile(r *gRPCReporter) (string, error) {
	// create base directory
	basePath := process.basePath
	if err := os.RemoveAll(basePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return basePath, err
	}
	if err := os.MkdirAll(basePath, 0o700); err != nil {
		return basePath, err
	}

	// create and write metadata file
	metadataFile := process.metadataFile
	if metaFile, err := os.Create(metadataFile); err != nil {
		return metadataFile, err
	} else {
		if content, err := p.buildMetadataContent(r); err != nil {
			return metadataFile, err
		} else if _, err = metaFile.WriteString(content); err != nil {
			return metadataFile, err
		}
	}
	return "", nil
}

func (p *processStat) buildMetadataContent(g *gRPCReporter) (string, error) {
	layer := g.layer
	if layer == "" {
		layer = "GENERAL"
	}

	propertiesJson, err := p.buildPropertiesJson(g)
	if err != nil {
		return "", err
	}

	metadata := map[string]string{
		"layer":         layer,
		"service_name":  g.service,
		"instance_name": g.serviceInstance,
		"process_name":  g.serviceInstance, // process name is same with instance name
		"properties":    propertiesJson,
		"labels":        strings.Join(g.processLabels, ","),
		"language":      "golang",
	}

	result := ""
	for k, v := range metadata {
		result += fmt.Sprintf("%s=%s\n", k, v)
	}
	return result, nil
}

func (p *processStat) buildPropertiesJson(g *gRPCReporter) (string, error) {
	props := buildOSInfo()
	if g.instanceProps != nil {
		for k, v := range g.instanceProps {
			props = append(props, &commonv3.KeyStringValuePair{
				Key:   k,
				Value: v,
			})
		}
	}

	properties := make(map[string]string)
	for _, p := range props {
		properties[p.Key] = p.Value
	}
	bytes, err := json.Marshal(properties)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func cleanupProcessDirectory(r *gRPCReporter) {
	if process == nil {
		return
	}
	process.shutdownOnce.Do(func() {
		if err := os.RemoveAll(process.basePath); err != nil && r != nil {
			r.logger.Warnf("could delete process: %s, %v", process.basePath, err)
		}
		process.status = Closed
	})
}
