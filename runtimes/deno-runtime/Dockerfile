FROM denoland/deno:alpine-1.38.1 AS production

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-toolkit" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Functions Runtime" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.licenses="Apache-2.0"

RUN apk update && apk upgrade --no-cache libcrypto3 libssl3

COPY entrypoint.sh /entrypoint.sh

USER deno

ENTRYPOINT /entrypoint.sh
