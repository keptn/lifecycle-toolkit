# Image URL to use all building/pushing image targets
CERT_MANAGER_VERSION ?= v1.8.0
TAG ?= "$(shell date +%Y%m%d%s)"
TAG := $(TAG)

# RELEASE_REGISTRY is the container registry to push
# into.
RELEASE_REGISTRY?=ghcr.io/keptn-sandbox
ARCH?=amd64

.PHONY: build-and-push-dev-images
build-and-push-dev-images:
	RELEASE_TAG=$(TAG)
	$(MAKE) -C operator release-local.$(ARCH) TAG=$(TAG)
	$(MAKE) -C scheduler release-local.$(ARCH) TAG=$(TAG)
	$(MAKE) -C operator push-local TAG=$(TAG)
	$(MAKE) -C scheduler push-local TAG=$(TAG)

.PHONY: build-dev-manifests
build-dev-manifests:
	$(MAKE) -C operator release-manifests TAG=$(TAG) ARCH=$(ARCH)
	$(MAKE) -C scheduler release-manifests TAG=$(TAG) ARCH=$(ARCH)
	if [[ ! -d manifests ]]; then mkdir manifests; fi
	cat operator/config/rendered/release.yaml > manifests/dev.yaml
	echo "---" >> manifests/dev.yaml
	cat scheduler/config/rendered/release.yaml >> manifests/dev.yaml

.PHONY: build-deploy-dev-environment
build-deploy-dev-environment: build-and-push-dev-images build-dev-manifests
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/$(CERT_MANAGER_VERSION)/cert-manager.yaml
	kubectl apply -f manifests/dev.yaml