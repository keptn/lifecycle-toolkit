FROM denoland/deno:alpine-1.40.2 AS production

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-toolkit" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Functions Runtime" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.licenses="Apache-2.0"

RUN apk update && \
    apk add libcrypto3~=3.0.12 libssl3~=3.0.12

COPY entrypoint.sh /entrypoint.sh

USER deno

ENTRYPOINT /entrypoint.sh
