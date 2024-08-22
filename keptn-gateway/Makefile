# Simple Makefile for a Go project
ARCHS = amd64 arm64
COMMONENVVAR=GOOS=$(shell uname -s | tr A-Z a-z)
BUILDENVVAR=CGO_ENABLED=0
TAG?="v1.0.1" # x-release-please-version
CHART_APPVERSION?=""

LOCAL_REGISTRY=localhost:5000/keptn-gateway
LOCAL_IMAGE=keptn-gateway:latest

# RELEASE_REGISTRY is the container registry to push
# into.
RELEASE_REGISTRY?=ghcr.io/keptn
RELEASE_IMAGE:=keptn-gateway:$(TAG)

# Build the application
all: build

build:
	@echo "Building..."
	
	
	@go build -o keptn-gateway cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go



# Test the application
test:
	@echo "Testing..."
	@go test ./... -v



# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
watch:
	@if command -v air > /dev/null; then \
			air --build.cmd "go build -o keptn-gateway cmd/api/main.go" --build.bin "./keptn-gateway"; \
			echo "Watching...";\
	else \
			read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
			if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
					go install github.com/air-verse/air@latest; \
					air --build.cmd "go build -o keptn-gateway cmd/api/main.go" --build.bin "./keptn-gateway"; \
					echo "Watching...";\
			else \
					echo "You chose not to install air. Exiting..."; \
					exit 1; \
			fi; \
	fi

.PHONY: all build run test clean watch

#### developer build   ####

.PHONY: build-and-push-local
build-and-push-local: release-local

.PHONY: release-local
release-local: release-local.amd64 release-local.arm64
	for arch in $(ARCHS); do \
			docker push $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-$${arch} ;\
		done
		DOCKER_CLI_EXPERIMENTAL=enabled docker manifest create $(RELEASE_REGISTRY)/$(RELEASE_IMAGE) $(addprefix --amend $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-, $(ARCHS))
		for arch in $(ARCHS); do \
			DOCKER_CLI_EXPERIMENTAL=enabled docker manifest annotate --arch $${arch} $(RELEASE_REGISTRY)/$(RELEASE_IMAGE) $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-$${arch} ;\
		done
		DOCKER_CLI_EXPERIMENTAL=enabled docker manifest push $(RELEASE_REGISTRY)/$(RELEASE_IMAGE) ;\

.PHONY: release-local.amd64
release-local.amd64: clean
	DOCKER_BUILDKIT=1 docker build \
		--build-arg GIT_HASH="$(HASH)" \
		--build-arg RELEASE_VERSION="$(RELEASE_VERSION)" \
		--build-arg BUILD_TIME="$(BUILD_TIME)" \
		-t $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-amd64 \
		--platform linux/amd64 \
		--target debug \
		.
	docker tag $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-amd64 $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)

.PHONY: release-local.arm64
release-local.arm64: clean
	DOCKER_BUILDKIT=1 docker build \
		--build-arg GIT_HASH="$(HASH)" \
		--build-arg RELEASE_VERSION="$(RELEASE_VERSION)" \
		--build-arg BUILD_TIME="$(BUILD_TIME)" \
		--build-arg BUILDPLATFORM="linux/arm64" \
		-t $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-arm64 \
		--platform linux/arm64 \
		--target debug \
		.
	docker tag $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-arm64 $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)

.PHONY: push-local
push-local:
	if [[ "${ARCH}" == "amd64" ]]; then \
		docker push $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-amd64; \
	elif [[ "${ARCH}" == "arm64" ]]; then \
		docker push $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-arm64 ;\
	fi
	docker push $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)

##golangci-lint
.PHONY: lint
lint:
	${GOPATH}/bin/golangci-lint run --config ../.golangci.yml -v
