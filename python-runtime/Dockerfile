FROM python:3.11-alpine AS production

LABEL org.opencontainers.image.source="https://github.com/keptn/lifecycle-toolkit" \
    org.opencontainers.image.url="https://keptn.sh" \
    org.opencontainers.image.title="Keptn Python Runtime" \
    org.opencontainers.image.vendor="Keptn" \
    org.opencontainers.image.licenses="Apache-2.0"

RUN pip install -q --disable-pip-version-check pyyaml GitPython requests

COPY entrypoint.sh /entrypoint.sh

USER 1000:1000

ENV CMD_ARGS=""
ENV SCRIPT=""

ENTRYPOINT /entrypoint.sh
