# Keptn Lifecycle Controller - Function Runtime

## Build
```
docker build -t keptnsandbox/klc-runtime:${VERSION} .
```

## Usage

### Docker
```
docker run -e SCRIPT=https://deno.land/std/examples/welcome.ts -it keptnsandbox/klc-runtime:${VERSION}
```

