# Copyright 2020 The Kubernetes Authors.
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
FROM --platform=$BUILDPLATFORM golang:1.19.4-alpine3.16 as builder

ARG SCHEDULER_COMPATIBLE_K8S_VERSION
ARG RELEASE_VERSION
ARG GIT_HASH
ARG BUILD_TIME
ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0

WORKDIR /adapter

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY pkg pkg

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags "\
    -X main.gitCommit=${GIT_HASH} \
    -X main.buildTime=${BUILD_TIME} \
    -X main.buildVersion=${RELEASE_VERSION} \
    -w" \
    -o bin/adapter cmd/adapter/main.go

FROM gcr.io/distroless/base-debian11:debug-nonroot as debug

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-toolkit" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Metrics Adapter" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.licenses="Apache-2.0"

COPY --from=builder /adapter/bin/adapter /bin/adapter

WORKDIR /bin
CMD ["adapter"]

FROM gcr.io/distroless/base-debian11:nonroot as production

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-toolkit" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Metrics Adapter" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.licenses="Apache-2.0"

COPY --from=builder /adapter/bin/adapter /bin/adapter
USER 65532:65532

WORKDIR /bin
CMD ["adapter"]

