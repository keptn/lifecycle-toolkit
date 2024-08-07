FROM --platform=$BUILDPLATFORM golang:1.22.5-alpine3.19 AS builder

ENV CGO_ENABLED=0

WORKDIR /keptn-gateway

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./

ARG RELEASE_VERSION
ARG GIT_HASH
ARG BUILD_TIME
ARG TARGETOS
ARG TARGETARCH

RUN GOOS=$TARGETOS GOARCH=$TARGETARCH go build -ldflags "\
  -X main.gitCommit=${GIT_HASH} \
  -X main.buildTime=${BUILD_TIME} \
  -X main.buildVersion=${RELEASE_VERSION} \
  -w" \
  -o bin/keptn-gateway ./main.go

FROM gcr.io/distroless/static-debian11:debug-nonroot AS debug

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-toolkit" \
  org.opencontainers.image.url="https://keptn.sh" \
  org.opencontainers.image.title="Keptn Gateway" \
  org.opencontainers.image.vendor="Keptn" \
  org.opencontainers.image.licenses="Apache-2.0"

COPY --from=builder /keptn-gateway/bin/keptn-gateway /bin/keptn-gateway

WORKDIR /bin

EXPOSE 8080

CMD ["keptn-gateway"]

FROM gcr.io/distroless/static-debian11:nonroot AS production

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-toolkit" \
  org.opencontainers.image.url="https://keptn.sh" \
  org.opencontainers.image.title="Keptn Gateway" \
  org.opencontainers.image.vendor="Keptn" \
  org.opencontainers.image.licenses="Apache-2.0"

COPY --from=builder /keptn-gateway/bin/keptn-gateway /bin/keptn-gateway
USER 65532:65532

WORKDIR /bin

EXPOSE 8080

CMD ["keptn-gateway"]
