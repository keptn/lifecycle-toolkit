# Keptn Lifecycle Controller - Deno Runtime

## Build

```shell
docker build -t keptnsandbox/klc-runtime:${VERSION} .
```

## Usage

The Keptn deno runtime uses [Deno](https://deno.land/)
to execute Javascript/Typescript code.
The Keptn Lifecycle Toolkit uses this runtime to run [KeptnTask](https://lifecycle.keptn.sh/docs/tasks/write-tasks/)
for pre- and post-checks.
The Keptn Lifecycle Toolkit passes parameters of `KeptnTask`s and
[Context](https://lifecycle.keptn.sh/docs/concepts/tasks/#context) information
to the runtime via the special environmental variable `DATA`.
It also supports mounting Kubernetes secrets making them accessible via the `SECURE_DATA` env var.
You can then read the data with the following snippet of code.

```js
const data = Deno.env.get("DATA")!;
const secret = Deno.env.get("SECURE_DATA")!;
console.log(data);
console.log(secret);
```

`KeptnTask`s can be tested locally with the runtime using the following command.
Replace `${VERSION}` with the KLT version of your choice.

```sh
docker run -v $(pwd)/test.ts:/test.ts -e SCRIPT=/test.ts -e DATA='{ "url":"http://localhost:9090" }' -e SECURE_DATA='{ "token": "myToken"}' -it ghcr.io/keptn/deno-runtime:${VERSION}
```

### Docker with function on webserver (function in this repo)

```shell
docker run \
  -e SCRIPT=https://raw.githubusercontent.com/keptn/lifecycle-toolkit/main/runtimes/deno-runtime/samples/ts/hello-world.ts \
  -it \
  keptnsandbox/klc-runtime:${VERSION}
```

### Docker with function and external data - scheduler

```shell
docker run \
  -e SCRIPT=https://raw.githubusercontent.com/keptn/lifecycle-toolkit/main/runtimes/deno-runtime/samples/ts/scheduler.ts \
  -e DATA='{ "targetDate":"2025-04-16T06:55:31.820Z" }' \
  -it \
  keptnsandbox/klc-runtime:${VERSION}
```

### Docker with function and external secure data - slack

```shell
docker run \
  -e SCRIPT=https://raw.githubusercontent.com/keptn/lifecycle-toolkit/main/runtimes/deno-runtime/samples/ts/slack.ts \
  -e SECURE_DATA='{ "slack_hook":"hook/parts","text":"this is my test message" }' \
  -it \
  keptnsandbox/klc-runtime:${VERSION}
```

### Docker with function and external data - prometheus

```shell
docker run \
  -e SCRIPT=https://raw.githubusercontent.com/keptn/lifecycle-toolkit/main/runtimes/deno-runtime/samples/ts/prometheus.ts \
  -e DATA='{ "url":"http://localhost:9090", "metrics": "up{service=\"kubernetes\"}", "expected_value": "1" }' \
  -it \
  ghcr.io/keptn/deno-runtime:${VERSION}
```

<!-- markdownlint-disable-next-line MD033 MD013 -->
<img referrerpolicy="no-referrer-when-downgrade" src="https://static.scarf.sh/a.png?x-pxid=858843d8-8da2-4ce5-a325-e5321c770a78" />
