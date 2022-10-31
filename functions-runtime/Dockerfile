FROM denoland/deno:1.27.0

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-controller" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Functions Runtime" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.licenses="Apache-2.0"

COPY entrypoint.sh /entrypoint.sh

USER deno

ENTRYPOINT /entrypoint.sh
