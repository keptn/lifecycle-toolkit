# Image URL to use all building/pushing image targets
IMG ?= controller:latest
CERT_MANAGER_VERSION ?= v1.8.0

# RELEASE_REGISTRY is the container registry to push
# into.
RELEASE_REGISTRY?=ghcr.io/keptn-sandbox
RELEASE_VERSION?=$(shell date +%Y%m%d%s)-v0.24.3#$(shell git describe --tags --match "v*")
TAG?=latest

ARCHS = amd64 arm64
COMMONENVVAR=GOOS=$(shell uname -s | tr A-Z a-z)
BUILDENVVAR=CGO_ENABLED=0

.PHONY: build-and-push-dev-images
build-and-push-dev-images:
	$(MAKE) -C operator build-and-push-local ARCH=amd64
	$(MAKE) -C scheduler build-and-push-local ARCH=amd4

.PHONY: build-dev-manifests
build-dev-manifests:
	$(MAKE) -C operator release-manifests ARCH=amd64
	$(MAKE) -C scheduler release-manifests ARCH=arm64
	if [[ ! -d manifests ]]; then mkdir manifests; fi
	cat operator/config/rendered/release.yaml > manifests/dev.yaml
	echo "---" >> manifests/dev.yaml
	cat scheduler/config/rendered/release.yaml >> manifests/dev.yaml

.PHONY: build-deploy-dev-environment
build-deploy-dev-environment: build-and-push-dev-images build-dev-manifests
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v${CERT_MANAGRER_VERSION}/cert-manager.yaml
	kubectl apply -f manifests/dev.yaml