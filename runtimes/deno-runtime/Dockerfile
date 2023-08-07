FROM denoland/deno:alpine-1.36.1 AS production

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-toolkit" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Functions Runtime" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.licenses="Apache-2.0"

COPY entrypoint.sh /entrypoint.sh

USER deno

ENTRYPOINT /entrypoint.sh
