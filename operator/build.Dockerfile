FROM --platform=$BUILDPLATFORM golang:1.20.4-alpine3.16 AS builder

ENV CGO_ENABLED=0

WORKDIR /workspace

# Copy the go source
COPY operator/ operator
COPY klt-cert-manager/ klt-cert-manager
COPY metrics-operator/ metrics-operator

# Setup go workspace
RUN go work init
RUN go work use ./operator
RUN go work use ./klt-cert-manager
RUN go work use ./metrics-operator

# renovate: datasource=github-releases depName=kubernetes-sigs/controller-tools
ARG CONTROLLER_TOOLS_VERSION=v0.12.0
RUN go install sigs.k8s.io/controller-tools/cmd/controller-gen@$CONTROLLER_TOOLS_VERSION

ARG GIT_HASH
ARG RELEASE_VERSION
ARG BUILD_TIME
ARG TARGETOS
ARG TARGETARCH

# Build
RUN cd operator && controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..." && \
  GOOS=$TARGETOS GOARCH=$TARGETARCH \
  go build -ldflags "\
    -w \
    -X common.gitCommit=${GIT_HASH} \
    -X common.buildTime=${BUILD_TIME} \
    -X common.buildVersion=${RELEASE_VERSION}" \
    -o bin/manager main.go

FROM gcr.io/distroless/static-debian11:debug-nonroot AS debug

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-toolkit" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Lifecycle Operator" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.licenses="Apache-2.0"

WORKDIR /
COPY --from=builder /workspace/operator/bin/manager .

ENTRYPOINT ["/manager"]

FROM gcr.io/distroless/static-debian11:nonroot AS production

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-toolkit" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Lifecycle Operator" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.licenses="Apache-2.0"

WORKDIR /
COPY --from=builder /workspace/bin/manager .
USER 65532:65532

ENTRYPOINT ["/manager"]
