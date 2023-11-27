# Keptn - Deno Runtime

## Build

```shell
docker build -t keptnsandbox/klc-runtime:${VERSION} .
```

## Usage

The Keptn deno runtime uses [Deno](https://deno.com/)
to execute Javascript/Typescript code.
Keptn uses this runtime to run [KeptnTask](https://lifecycle.keptn.sh/docs/tasks/write-tasks/)
for pre- and post-checks.

### Environment Variables

Keptn passes the following environment variables to the runtime:

* `DATA`: JSON encoded object containing the parameters specified in `spec.parameters` of a `KeptnTask`.
* `SECURE_DATA`: Contains the value of the secret referenced in the `spec.secureParameters` field of a `KeptnTask`.
* `KEPTN_CONTEXT`: JSON encoded object containing context information for the task.

You can then read the data with the following snippet of code.

```js
const data = Deno.env.get("DATA")!;
const secret = Deno.env.get("SECURE_DATA")!;
const context = Deno.env.get("KEPTN_CONTEXT")!;
console.log(data);
console.log(secret);
console.log(context);
```

`KeptnTask`s can be tested locally with the runtime using the following command.
Replace `${VERSION}` with the Keptn version of your choice.

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
