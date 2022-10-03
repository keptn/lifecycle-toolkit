# Copyright 2020 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

ARCHS = amd64 arm64
COMMONENVVAR=GOOS=$(shell uname -s | tr A-Z a-z)
BUILDENVVAR=CGO_ENABLED=0
TAG?=latest

LOCAL_REGISTRY=localhost:5000/scheduler-plugins
LOCAL_IMAGE=kube-scheduler:latest
LOCAL_CONTROLLER_IMAGE=controller:latest

# RELEASE_REGISTRY is the container registry to push
# into.
RELEASE_REGISTRY?=ghcr.io/keptn-sandbox
RELEASE_VERSION?=$(shell date +%Y%m%d%s)-v0.24.3#$(shell git describe --tags --match "v*")
RELEASE_IMAGE:=lfc-scheduler:$(TAG)
RELEASE_CONTROLLER_IMAGE:=controller:$(RELEASE_VERSION)

# VERSION is the scheduler's version
#
# The RELEASE_VERSION variable can have one of two formats:
# v20201009-v0.18.800-46-g939c1c0 - automated build for a commit(not a tag) and also a local build
# v20200521-v0.18.800             - automated build for a tag
VERSION=$(shell echo $(RELEASE_VERSION) | awk -F - '{print $$2}')


## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Versions
KUSTOMIZE_VERSION ?= v3.8.7

## Tool Binaries
KUSTOMIZE ?= $(LOCALBIN)/kustomize
KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"

.PHONY: all
all: build

.PHONY: build
build: build-scheduler

.PHONY: build.amd64
build.amd64: build-controller.amd64 build-scheduler.amd64

.PHONY: build.arm64v8
build.arm64v8: build-controller.arm64v8 build-scheduler.arm64v8

.PHONY: build-controller
build-controller:
	$(COMMONENVVAR) $(BUILDENVVAR) go build -ldflags '-w' -o bin/controller cmd/controller/controller.go

.PHONY: build-controller.amd64
build-controller.amd64:
	$(COMMONENVVAR) $(BUILDENVVAR) GOARCH=amd64 go build -ldflags '-w' -o bin/controller cmd/controller/controller.go

.PHONY: build-controller.arm64v8
build-controller.arm64v8:
	GOOS=linux $(BUILDENVVAR) GOARCH=arm64 go build -ldflags '-w' -o bin/controller cmd/controller/controller.go

.PHONY: build-scheduler
build-scheduler:
	$(COMMONENVVAR) $(BUILDENVVAR) go build -ldflags '-X k8s.io/component-base/version.gitVersion=$(VERSION) -w' -o bin/kube-scheduler cmd/scheduler/main.go

.PHONY: build-scheduler.amd64
build-scheduler.amd64:
	$(COMMONENVVAR) $(BUILDENVVAR) GOARCH=amd64 go build -ldflags '-X k8s.io/component-base/version.gitVersion=$(VERSION) -w' -o bin/kube-scheduler cmd/scheduler/main.go

.PHONY: build-scheduler.arm64v8
build-scheduler.arm64v8:
	GOOS=linux $(BUILDENVVAR) GOARCH=arm64 go build -ldflags '-X k8s.io/component-base/version.gitVersion=$(VERSION) -w' -o bin/kube-scheduler cmd/scheduler/main.go

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: fmt vet ## Run tests.
	go test ./... -coverprofile cover.out

.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	test -s $(LOCALBIN)/kustomize || { curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN); }

.PHONY: release-manifests
release-manifests: kustomize
	cd manifests/install && $(KUSTOMIZE) edit set image klfc-scheduler=$(RELEASE_REGISTRY)/$(RELEASE_IMAGE)
	$(KUSTOMIZE) build manifests/install > config/rendered/release.yaml

.PHONY: local-image
local-image: clean
	docker build --build-arg ARCH="amd64" --build-arg RELEASE_VERSION="$(RELEASE_VERSION)" -t $(LOCAL_REGISTRY)/$(LOCAL_IMAGE) .

.PHONY: build-and-push-image
build-and-push-image: release-image push-release-images

.PHONY: release-image
release-image: release-image.amd64 release-image.arm64v8

.PHONY: release-image.amd64
release-image.amd64: clean
	docker build --load --build-arg ARCH="amd64" --build-arg RELEASE_VERSION="$(RELEASE_VERSION)" -t $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-amd64 .

.PHONY: release-image.arm64v8
release-image.arm64v8: clean
	docker build --load --build-arg ARCH="arm64v8" --build-arg RELEASE_VERSION="$(RELEASE_VERSION)" -t $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-arm64 .

.PHONY: push-release-images
push-release-images: release-image.amd64 release-image.arm64v8
	for arch in $(ARCHS); do \
		docker push $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-$${arch} ;\
	done
	DOCKER_CLI_EXPERIMENTAL=enabled docker manifest create $(RELEASE_REGISTRY)/$(RELEASE_IMAGE) $(addprefix --amend $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-, $(ARCHS))
	for arch in $(ARCHS); do \
		DOCKER_CLI_EXPERIMENTAL=enabled docker manifest annotate --arch $${arch} $(RELEASE_REGISTRY)/$(RELEASE_IMAGE) $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-$${arch} ;\
	done
	DOCKER_CLI_EXPERIMENTAL=enabled docker manifest push $(RELEASE_REGISTRY)/$(RELEASE_IMAGE) ;\

.PHONY: clean
clean:
	rm -rf ./bin


#### developer build   ####

.PHONY: build-and-push-local
build-and-push-local: release-local

.PHONY: release-local
release-local: release-local.amd64 release-local.arm64v8
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
	docker build --build-arg ARCH="amd64" --build-arg RELEASE_VERSION="$(RELEASE_VERSION)" -t $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-amd64 .

.PHONY: release-local.arm64v8
release-local.arm64v8: clean
	docker build --build-arg ARCH="arm64v8" --build-arg RELEASE_VERSION="$(RELEASE_VERSION)" -t $(RELEASE_REGISTRY)/$(RELEASE_IMAGE)-arm64 .
