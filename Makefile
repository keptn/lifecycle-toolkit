# Image URL to use all building/pushing image targets

# renovate: datasource=github-releases depName=cert-manager/cert-manager
CERT_MANAGER_VERSION ?= v1.11.0
TAG ?= "$(shell date +%Y%m%d%s)"
TAG := $(TAG)

# RELEASE_REGISTRY is the container registry to push
# into.
RELEASE_REGISTRY?=ghcr.io/keptn
ARCH?=amd64

.PHONY: integration-test #these tests should run on a real cluster!
integration-test:
	kubectl kuttl test --start-kind=false ./test/integration/ --config=kuttl-test.yaml

.PHONY: integration-test-local #these tests should run on a real cluster!
integration-test-local:
	kubectl kuttl test --start-kind=false ./test/integration/ --config=kuttl-test-local.yaml

.PHONY: cleanup-manifests
cleanup-manifests:
	rm -rf manifests

.PHONY: build-deploy-operator
build-deploy-operator: deploy-cert-manager
	$(MAKE) -C operator release-local.$(ARCH) RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C operator push-local RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C operator release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG) ARCH=$(ARCH)

	kubectl apply -f operator/config/rendered/release.yaml

.PHONY: deploy-cert-manager
deploy-cert-manager:
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/$(CERT_MANAGER_VERSION)/cert-manager.yaml

.PHONY: build-deploy-dev-environment
build-deploy-dev-environment: build-deploy-operator
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/$(CERT_MANAGER_VERSION)/cert-manager.yaml
