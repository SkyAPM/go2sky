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
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	mockCollectorDockerfile = `
FROM openjdk:8

WORKDIR /tests

ARG COMMIT_HASH=8a48c49b4420df5c9576d2aea178b2ebcb7ecd09

ADD https://github.com/apache/skywalking-agent-test-tool/archive/${COMMIT_HASH}.tar.gz .

RUN tar -xf ${COMMIT_HASH}.tar.gz --strip 1

RUN rm ${COMMIT_HASH}.tar.gz

RUN ./mvnw -B -DskipTests package

FROM openjdk:8

EXPOSE 19876 12800

WORKDIR /tests

COPY --from=0 /tests/dist/skywalking-mock-collector.tar.gz /tests

RUN tar -xf skywalking-mock-collector.tar.gz --strip 1

RUN chmod +x bin/collector-startup.sh

ENTRYPOINT bin/collector-startup.sh
`
)

func (p *TestPlugin) runMockCollector(ctx context.Context, t *testing.T) {
	// dynamic build context, see more https://golang.testcontainers.org/features/build_from_dockerfile/
	var buf bytes.Buffer
	tarWriter := tar.NewWriter(&buf)
	hdr := &tar.Header{
		Name: "Dockerfile",
		Mode: 0600,
		Size: int64(len(mockCollectorDockerfile)),
	}
	if err := tarWriter.WriteHeader(hdr); err != nil {
		t.Error(err)
	}
	if _, err := tarWriter.Write([]byte(mockCollectorDockerfile)); err != nil {
		t.Error(err)
	}
	if err := tarWriter.Close(); err != nil {
		t.Error(err)
	}
	reader := bytes.NewReader(buf.Bytes())

	//  create mock collector container
	fromDockerfile := testcontainers.FromDockerfile{
		ContextArchive: reader,
	}

	req := testcontainers.ContainerRequest{
		FromDockerfile: fromDockerfile,
		ExposedPorts:   []string{"12800/tcp", "19876/tcp"},
		WaitingFor:     wait.ForLog("Started"),
	}

	collectorC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Error(err)
	}
	p.addTerminateFunc(func() {
		_ = collectorC.Terminate(ctx)
	})

	//
	ip, err := collectorC.Host(ctx)
	if err != nil {
		t.Error(err)
	}
	grpcPort, err := collectorC.MappedPort(ctx, "19876")
	if err != nil {
		t.Error(err)
	}
	p.grpcServerAddr = fmt.Sprintf("%s:%s", ip, grpcPort.Port())

	httpPort, err := collectorC.MappedPort(ctx, "12800")
	if err != nil {
		t.Error(err)
	}
	p.httpServerAddr = fmt.Sprintf("%s:%s", ip, httpPort.Port())
}
