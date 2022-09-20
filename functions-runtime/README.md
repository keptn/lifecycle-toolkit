# Keptn Lifecycle Controller - Function Runtime

## Build
```
docker build -t keptnsandbox/klc-runtime:${VERSION} .
```

## Usage

### Docker with function on webserver (function in this repo)
```
docker run -e SCRIPT=https://raw.githubusercontent.com/keptn-sandbox/lifecycle-controller/main/functions-runtime/samples/ts/hello-world.ts -it keptnsandbox/klc-runtime:${VERSION}
```

### Docker with function and external data - scheduler
```
docker run -e SCRIPT=https://raw.githubusercontent.com/keptn-sandbox/lifecycle-controller/main/functions-runtime/samples/ts/scheduler.ts -e DATA='{ "targetDate":"2025-04-16T06:55:31.820Z" }' -it keptnsandbox/klc-runtime:${VERSION}
```

### Docker with function and external secure data - slack
```
docker run -e SCRIPT=https://raw.githubusercontent.com/keptn-sandbox/lifecycle-controller/main/functions-runtime/samples/ts/slack.ts -e SECURE_DATA='{ "slack_hook":"hook/parts","text":"this is my test message" }' -it keptnsandbox/klc-runtime:${VERSION}
```

