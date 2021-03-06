# Copyright 2020 SUSE
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

FROM golang:1.14 AS builder

ARG DUMB_INIT_VERSION="1.2.2"
ARG CF_CLI_VERSION="7.0.2"

RUN curl -L -o /usr/local/bin/dumb-init \
      "https://github.com/Yelp/dumb-init/releases/download/v${DUMB_INIT_VERSION}/dumb-init_${DUMB_INIT_VERSION}_amd64"
RUN chmod +x /usr/local/bin/dumb-init

RUN go get -u github.com/onsi/ginkgo/ginkgo
RUN curl -L "https://packages.cloudfoundry.org/stable?release=linux64-binary&version=${CF_CLI_VERSION}" \
      | tar zx -C /usr/local/bin cf cf7

WORKDIR /minibroker-integration-tests
# Copy the go.mod over so docker can cache the module downloads if possible.
COPY go.mod go.sum ./
RUN go mod download

COPY mits/ ./mits/
RUN ginkgo build ./mits

###############################################################################

FROM opensuse/leap:15.1

WORKDIR /minibroker-integration-tests
COPY --from=builder /go/bin/ginkgo /usr/local/bin/ginkgo
COPY --from=builder /usr/local/bin/cf /usr/local/bin/cf
COPY --from=builder /minibroker-integration-tests/mits/mits.test ./mits/mits.test
COPY --from=builder /usr/local/bin/dumb-init /usr/local/bin/dumb-init
COPY ./assets/ ./mits/assets/

ENTRYPOINT ["/usr/local/bin/dumb-init", "--"]
