# Copyright 2022 SkyAPM org
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -ex

if ! command -v yq &>/dev/null; then
  # prepare base dir
  BASE_DIR=/tmp/skywalking-infra-e2e/swctl
  BIN_DIR=/usr/local/bin
  mkdir -p $BASE_DIR && cd $BASE_DIR

  curl -kLo yq.tar.gz https://github.com/mikefarah/yq/archive/v4.14.1.tar.gz
  tar -zxf yq.tar.gz --strip=1
  go install && go build -ldflags -s && cp yq $BIN_DIR/
fi
