# Keptn Lifecycle Controller - Function Runtime

## Build

```shell
docker build -t keptnsandbox/klc-runtime:${VERSION} .
```

## Usage

The Keptn function runtime allows the execution of code written in various programming languages. The Keptn Lifecycle Toolkit supports the following runtime options:

- [Deno](https://deno.land/) Runtime: Uses Deno to execute 
  JavaScript/TypeScript code.
- Python Runtime: Uses Python to execute Python code.
- Container Runtime: Executes code inside a containerized 
  environment.

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
docker run -v $(pwd)/test.ts:/test.ts -e SCRIPT=/test.ts -e DATA='{ "url":"http://localhost:9090" }' -e SECURE_DATA='{ "token": "myToken"}' -it ghcr.io/keptn/functions-runtime:${VERSION}
```
## Testing with Different Runtimes

### Deno Runtime

To test KeptnTasks locally using the Deno runtime, you can use the following command:

```shell
docker run -v $(pwd)/test.ts:/test.ts -e SCRIPT=/test.ts -e DATA='{ "url":"http://localhost:9090" }' -e SECURE_DATA='{ "token": "myToken" }' -it ghcr.io/keptn/functions-runtime:${VERSION}
```

### Python Runtime

To test KeptnTasks using the Python runtime,
you can use the following command:

```shell
docker run -v $(pwd)/test.py:/test.py -e SCRIPT=/test.py -e DATA='{ "url":"http://localhost:9090" }' -e SECURE_DATA='{ "token": "myToken" }' -it ghcr.io/keptn/python-runtime:${VERSION}
```

### Container Runtime

To test KeptnTasks using the Container runtime, you can use the following command:

```shell
docker run -v $(pwd)/test.sh:/test.sh -e SCRIPT=/test.sh -e DATA='{ "url":"http://localhost:9090" }' -e SECURE_DATA='{ "token": "myToken" }' -it ghcr.io/keptn/container-runtime:${VERSION}
```

Please replace `${VERSION}` with the desired version of the Keptn Lifecycle Toolkit or the respective runtime version.

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
  ghcr.io/keptn/functions-runtime:${VERSION}
```

<!-- markdownlint-disable-next-line MD033 MD013 -->
<img referrerpolicy="no-referrer-when-downgrade" src="https://static.scarf.sh/a.png?x-pxid=858843d8-8da2-4ce5-a325-e5321c770a78" />
