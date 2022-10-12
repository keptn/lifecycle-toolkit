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
FROM golang:1.18.0 as builder

ARG ARCH
ARG RELEASE_VERSION
ARG GIT_HASH

WORKDIR /scheduler

COPY go.mod go.sum Makefile ./

RUN go mod download

COPY cmd cmd
COPY pkg pkg

RUN RELEASE_VERSION=${RELEASE_VERSION} HASH=${GIT_HASH} make build-scheduler.$ARCH

FROM gcr.io/distroless/base-debian11:nonroot

LABEL org.opencontainers.image.source="https://github.com/keptn-sandbox/lifecycle-controller" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Lifecycle Scheduler" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.licenses="Apache-2.0"

COPY --from=builder /scheduler/bin/kube-scheduler /bin/kube-scheduler
USER 65532:65532

WORKDIR /bin
CMD ["kube-scheduler"]
