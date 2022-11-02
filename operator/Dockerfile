# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.19.3-alpine3.16 as builder

ARG CONTROLLER_TOOLS_VERSION=v0.9.2
ARG GIT_HASH
ARG RELEASE_VERSION
ARG BUILD_TIME
ARG TARGETOS
ARG TARGETARCH

ENV CGO_ENABLED=0

WORKDIR /workspace
RUN go install sigs.k8s.io/controller-tools/cmd/controller-gen@$CONTROLLER_TOOLS_VERSION

COPY go.mod go.sum ./
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY hack/ hack/
COPY api/ api/
COPY controllers/ controllers/
COPY webhooks/ webhooks/

# Build
RUN controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..." && \
  GOOS=$TARGETOS GOARCH=$TARGETARCH \
  go build -ldflags '\
    -w \
    -X main.gitCommit=$GIT_HASH \
    -X main.buildTime=$BUILD_TIME \
    -X main.buildVersion=$RELEASE_VERSION' \
    -o bin/manager main.go

FROM gcr.io/distroless/base-debian11:debug-nonroot as debug

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-controller" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Lifecycle Operator" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.licenses="Apache-2.0"

WORKDIR /
COPY --from=builder /workspace/bin/manager .

ENTRYPOINT ["/manager"]

FROM gcr.io/distroless/base-debian11:nonroot as production

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-controller" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Lifecycle Operator" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.licenses="Apache-2.0"

WORKDIR /
COPY --from=builder /workspace/bin/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
