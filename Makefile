# Image URL to use all building/pushing image targets

# renovate: datasource=github-releases depName=cert-manager/cert-manager
CERT_MANAGER_VERSION ?= v1.10.1
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

.PHONY: load-test
load-test:
	kubectl port-forward -n monitoring svc/prometheus-k8s 9090 &
	sleep 10 # TODO make this reliable and not as dumb as that
	kube-burner init -c ./test/load/cfg.yml --metrics-profile ./test/load/metrics.yml

.PHONY: cleanup-manifests
cleanup-manifests:
	rm -rf manifests

.PHONY: build-deploy-operator
build-deploy-operator: deploy-cert-manager
	$(MAKE) -C operator release-local.$(ARCH) RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C operator push-local RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C operator release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG) ARCH=$(ARCH)

	kubectl apply -f operator/config/rendered/release.yaml

.PHONY: build-deploy-scheduler
build-deploy-scheduler:
	$(MAKE) -C scheduler release-local.$(ARCH) RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C scheduler push-local RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C scheduler release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG) ARCH=$(ARCH)
	kubectl create namespace keptn-lifecycle-toolkit-system --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f scheduler/config/rendered/release.yaml

.PHONY: deploy-cert-manager
deploy-cert-manager:
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/$(CERT_MANAGER_VERSION)/cert-manager.yaml

.PHONY: build-deploy-dev-environment
build-deploy-dev-environment: build-deploy-operator build-deploy-scheduler
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/$(CERT_MANAGER_VERSION)/cert-manager.yaml
