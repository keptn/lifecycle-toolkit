# Keptn Lifecycle Controller - Function Runtime

## Build

```shell
docker build -t keptnsandbox/klc-runtime:${VERSION} .
```

## Usage

### Docker with function on webserver (function in this repo)

```shell
docker run \
  -e SCRIPT=https://raw.githubusercontent.com/keptn/lifecycle-toolkit/main/functions-runtime/samples/ts/hello-world.ts \
  -it \
  keptnsandbox/klc-runtime:${VERSION}
```

### Docker with function and external data - scheduler

```shell
docker run \
  -e SCRIPT=https://raw.githubusercontent.com/keptn/lifecycle-toolkit/main/functions-runtime/samples/ts/scheduler.ts \
  -e DATA='{ "targetDate":"2025-04-16T06:55:31.820Z" }' \
  -it \
  keptnsandbox/klc-runtime:${VERSION}
```

### Docker with function and external secure data - slack

```shell
docker run \
  -e SCRIPT=https://raw.githubusercontent.com/keptn/lifecycle-toolkit/main/functions-runtime/samples/ts/slack.ts \
  -e SECURE_DATA='{ "slack_hook":"hook/parts","text":"this is my test message" }' \
  -it \
  keptnsandbox/klc-runtime:${VERSION}
```

### Docker with function and external data - prometheus

```shell
docker run \
  -e SCRIPT=https://raw.githubusercontent.com/keptn/lifecycle-toolkit/main/functions-runtime/samples/ts/prometheus.ts \
  -e DATA='{ "url":"http://localhost:9090", "metrics": "up{service=\"kubernetes\"}", "expected_value": "1" }' \
  -it \
  ghcr.keptn.sh/keptn/functions-runtime:${VERSION}
```

<!-- markdownlint-disable-next-line MD033 MD013 -->
<img referrerpolicy="no-referrer-when-downgrade" src="https://static.scarf.sh/a.png?x-pxid=858843d8-8da2-4ce5-a325-e5321c770a78" />
