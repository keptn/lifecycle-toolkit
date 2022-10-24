# Build the manager binary
FROM golang:1.18 as builder

ARG ARCH
ARG GIT_HASH
ARG RELEASE_VERSION

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.sum Makefile ./

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY hack/ hack/
COPY api/ api/
COPY controllers/ controllers/
COPY webhooks/ webhooks/

# Build
RUN make build.$ARCH HASH=${GIT_HASH} TAG=${RELEASE_VERSION}

# Use distroless as minimal base image to package the manager binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/base-debian11:nonroot

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-controller" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Lifecycle Operator" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.licenses="Apache-2.0"

WORKDIR /
COPY --from=builder /workspace/bin/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
