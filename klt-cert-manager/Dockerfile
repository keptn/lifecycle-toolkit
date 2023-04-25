# Build the manager binary
FROM --platform=$BUILDPLATFORM golang:1.20.3-alpine3.16 AS builder

ENV CGO_ENABLED=0

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.sum ./
RUN go mod download

# Copy the go source
COPY ./ ./

# renovate: datasource=github-releases depName=kubernetes-sigs/controller-tools
ARG CONTROLLER_TOOLS_VERSION=v0.11.3
RUN go install sigs.k8s.io/controller-tools/cmd/controller-gen@$CONTROLLER_TOOLS_VERSION

ARG GIT_HASH
ARG RELEASE_VERSION
ARG BUILD_TIME
ARG TARGETOS
ARG TARGETARCH

# Build
RUN controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..." && \
  GOOS=$TARGETOS GOARCH=$TARGETARCH \
  go build -ldflags "\
    -w \
    -X main.gitCommit=${GIT_HASH} \
    -X main.buildTime=${BUILD_TIME} \
    -X main.buildVersion=${RELEASE_VERSION}" \
    -a -o manager main.go

FROM gcr.io/distroless/static-debian11:debug-nonroot AS debug

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-toolkit" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Lifecycle Certificate Manager" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.licenses="Apache-2.0"

WORKDIR /
COPY --from=builder /workspace/manager .

ENTRYPOINT ["/manager"]

FROM gcr.io/distroless/static-debian11:nonroot AS production

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-toolkit/klt-cert-manager" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Lifecycle Certificate Manager" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.licenses="Apache-2.0"

WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
