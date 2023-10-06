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
TAG?="v0.8.2" # x-release-please-version
CHART_APPVERSION?=""

LOCAL_REGISTRY=localhost:5000/scheduler-plugins
LOCAL_IMAGE=kube-scheduler:latest

# RELEASE_REGISTRY is the container registry to push
# into.
RELEASE_REGISTRY?=ghcr.io/keptn
RELEASE_IMAGE:=scheduler:$(TAG)

# VERSION is the scheduler's version
#
# The SCHEDULER_COMPATIBLE_K8S_VERSION variable can have one of two formats:
# v20201009-v0.18.800-46-g939c1c0 - automated build for a commit(not a tag) and also a local build
# v20200521-v0.18.800             - automated build for a tag
RELEASE_VERSION?=$(shell date +%Y%m%d%s)-v0.24.3
SCHEDULER_COMPATIBLE_K8S_VERSION=$(shell echo $(RELEASE_VERSION) | awk -F - '{print $$2}')

BUILD_TIME=$(shell date -u "+%F_%T")

# Compute the current Git commit hash
HASH?=$(shell git rev-parse HEAD)

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Versions
# renovate: datasource=github-tags depName=kubernetes-sigs/kustomize
KUSTOMIZE_VERSION ?= v5.1.1

## Tool Binaries
KUSTOMIZE ?= $(LOCALBIN)/kustomize
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
ENVTEST ?= $(LOCALBIN)/setup-envtest
KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"

.PHONY: all
all: build

.PHONY: build
build: build-scheduler

.PHONY: build-scheduler
build-scheduler:
	$(COMMONENVVAR) $(BUILDENVVAR) go build -ldflags '\
		-X k8s.io/component-base/version.gitVersion=$(SCHEDULER_COMPATIBLE_K8S_VERSION) \
		-X main.gitCommit=$(HASH) \
		-X main.buildTime=$(BUILD_TIME) \
		-X main.buildVersion=$(TAG) \
		-w' \
		-o bin/kube-scheduler cmd/scheduler/main.go

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: e2e-test
e2e-test: manifests envtest ## Run tests.
	go test ./test/e2e -v -coverprofile cover.out

.PHONY: unit-test
unit-test: manifests fmt vet envtest## Run tests.
	go test ./pkg/... -coverprofile cover-pkg.out
	go test ./cmd/scheduler -coverprofile cover-main.out
	sed -i '/mode: set/d' cover-pkg.out
	sed -i '/mode: set/d' cover-main.out
	echo "mode: set" > cover.out
	cat cover-main.out cover-pkg.out >> cover.out
	rm cover-pkg.out cover-main.out

#command(make test) to run all tests 
.PHONY: test
test: e2e-test unit-test


.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	test -s $(LOCALBIN)/kustomize || { curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN); }

.PHONY: release-manifests
release-manifests:
	kustomize build manifests/install > config/rendered/release.yaml
	envsubst < config/rendered/release.yaml > tmp.yaml; mv tmp.yaml config/rendered/release.yaml

envtest: $(ENVTEST) ## Download envtest-setup locally if necessary.
$(ENVTEST): $(LOCALBIN)
	test -s $(LOCALBIN)/setup-envtest || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-runtime/tools/setup-envtest@latest

.PHONY: local-image
local-image: clean
	docker build \
		--build-arg GIT_HASH="$(HASH)" \
		--build-arg RELEASE_VERSION="$(RELEASE_VERSION)" \
		--build-arg SCHEDULER_COMPATIBLE_K8S_VERSION="$(SCHEDULER_COMPATIBLE_K8S_VERSION)" \
		--build-arg BUILD_TIME="$(BUILD_TIME)" \
		-t $(LOCAL_REGISTRY)/$(LOCAL_IMAGE) .

.PHONY: clean
clean:
	rm -rf ./bin


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
		--build-arg SCHEDULER_COMPATIBLE_K8S_VERSION="$(SCHEDULER_COMPATIBLE_K8S_VERSION)" \
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
		--build-arg SCHEDULER_COMPATIBLE_K8S_VERSION="$(SCHEDULER_COMPATIBLE_K8S_VERSION)" \
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
