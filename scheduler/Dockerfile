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
FROM --platform=$BUILDPLATFORM golang:1.18.7-alpine3.16 as builder

ARG SCHEDULER_COMPATIBLE_K8S_VERSION
ARG RELEASE_VERSION
ARG GIT_HASH
ARG BUILD_TIME
ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0

WORKDIR /scheduler

COPY go.mod go.sum ./
RUN go mod download

COPY cmd cmd
COPY pkg pkg

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags '\
    -X k8s.io/component-base/version.gitVersion=$SCHEDULER_COMPATIBLE_K8S_VERSION \
    -X main.gitCommit=$GIT_HASH \
    -X main.buildTime=$BUILD_TIME \
    -X main.buildVersion=$RELEASE_VERSION \
    -w' \
    -o bin/kube-scheduler cmd/scheduler/main.go

FROM gcr.io/distroless/base-debian11:debug-nonroot as debug

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-controller" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Lifecycle Scheduler" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.licenses="Apache-2.0"

COPY --from=builder /scheduler/bin/kube-scheduler /bin/kube-scheduler

WORKDIR /bin
CMD ["kube-scheduler"]

FROM gcr.io/distroless/base-debian11:nonroot as production

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-controller" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Lifecycle Scheduler" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.licenses="Apache-2.0"

COPY --from=builder /scheduler/bin/kube-scheduler /bin/kube-scheduler
USER 65532:65532

WORKDIR /bin
CMD ["kube-scheduler"]
