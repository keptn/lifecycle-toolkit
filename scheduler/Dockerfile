FROM --platform=$BUILDPLATFORM golang:1.20.4-alpine3.16 AS builder

ENV CGO_ENABLED=0

WORKDIR /scheduler

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

ARG SCHEDULER_COMPATIBLE_K8S_VERSION
ARG RELEASE_VERSION
ARG GIT_HASH
ARG BUILD_TIME
ARG TARGETOS
ARG TARGETARCH

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags "\
    -X k8s.io/component-base/version.gitVersion=${SCHEDULER_COMPATIBLE_K8S_VERSION} \
    -X main.gitCommit=${GIT_HASH} \
    -X main.buildTime=${BUILD_TIME} \
    -X main.buildVersion=${RELEASE_VERSION} \
    -w" \
    -o bin/kube-scheduler cmd/scheduler/main.go

FROM gcr.io/distroless/static-debian11:debug-nonroot AS debug

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-toolkit" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Lifecycle Scheduler" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.licenses="Apache-2.0"

COPY --from=builder /scheduler/bin/kube-scheduler /bin/kube-scheduler

WORKDIR /bin
CMD ["kube-scheduler"]

FROM gcr.io/distroless/static-debian11:nonroot AS production

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-toolkit" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Lifecycle Scheduler" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.licenses="Apache-2.0"

COPY --from=builder /scheduler/bin/kube-scheduler /bin/kube-scheduler
USER 65532:65532

WORKDIR /bin
CMD ["kube-scheduler"]
