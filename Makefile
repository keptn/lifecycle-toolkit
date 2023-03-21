# Image URL to use all building/pushing image targets

# renovate: datasource=github-releases depName=cert-manager/cert-manager
CERT_MANAGER_VERSION ?= v1.11.0
# renovate: datasource=github-tags depName=kubernetes-sigs/kustomize
KUSTOMIZE_VERSION?=v4.5.7
# renovate: datasource=github-tags depName=helm/helm
HELM_VERSION ?= v3.11.2
CHART_APPVERSION ?= v0.7.0 # x-release-please-version

# renovate: datasource=docker depName=cytopia/yamllint
YAMLLINT_VERSION ?= alpine-1-0.14

# RELEASE_REGISTRY is the container registry to push
# into.
RELEASE_REGISTRY?=ghcr.keptn.sh/keptn
ARCH?=amd64
TAG ?= "$(shell date +%Y%m%d%s)"
TAG := $(TAG)

## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

## Tool Binaries
KUSTOMIZE ?= $(LOCALBIN)/kustomize
HELMIFY ?=  $(LOCALBIN)/helmify

.PHONY: helmify
helmify: $(HELMIFY) ## Download helmify locally if necessary.
$(HELMIFY): $(LOCALBIN)
	test -s $(LOCALBIN)/helmify || GOBIN=$(LOCALBIN) go install github.com/keptn/helmify/cmd/helmify@b1da2bb756ec4328bac7645da037a6fb4e6f30cf

.PHONY: integration-test #these tests should run on a real cluster!
integration-test:
	kubectl kuttl test --start-kind=false ./test/integration/ --config=kuttl-test.yaml

.PHONY: integration-test-local #these tests should run on a real cluster!
integration-test-local:
	kubectl kuttl test --start-kind=false ./test/integration/ --config=kuttl-test-local.yaml

.PHONY: load-test
load-test:
	kubectl apply -f ./test/load/assets/templates/namespace.yaml
	kubectl apply -f ./test/load/assets/templates/provider.yaml
	kube-burner init -c ./test/load/cfg.yml --metrics-profile ./test/load/metrics.yml

.PHONY: cleanup-manifests
cleanup-manifests:
	rm -rf manifests

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	test -s $(LOCALBIN)/kustomize || { curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN); }

.PHONY: release-helm-manifests
release-helm-manifests: helmify
	echo "building helm overlay"
	kustomize build ./helm/overlay  > helmchart.yaml
	envsubst < helmchart.yaml > tmp.yaml; mv tmp.yaml helmchart.yaml
	cat helmchart.yaml | $(HELMIFY) -probes=true -image-pull-secrets=true -vv helm/chart

.PHONY: helm-package
helm-package:
	$(MAKE) build-release-manifests CHART_APPVERSION=$(CHART_APPVERSION) RELEASE_REGISTRY=$(RELEASE_REGISTRY)
	$(MAKE) release-helm-manifests CHART_APPVERSION=$(CHART_APPVERSION) RELEASE_REGISTRY=$(RELEASE_REGISTRY)


.PHONY: build-release-manifests
build-release-manifests:
	$(MAKE) -C operator generate
	$(MAKE) -C klt-cert-manager generate
	$(MAKE) -C metrics-operator generate

	$(MAKE) -C operator release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG) ARCH=$(ARCH) CHART_APPVERSION=$(CHART_APPVERSION)
	$(MAKE) -C scheduler release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG) ARCH=$(ARCH) CHART_APPVERSION=$(CHART_APPVERSION)
	$(MAKE) -C klt-cert-manager release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG) ARCH=$(ARCH) CHART_APPVERSION=$(CHART_APPVERSION)
	$(MAKE) -C metrics-operator release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG) ARCH=$(ARCH) CHART_APPVERSION=$(CHART_APPVERSION)

.PHONY: build-deploy-operator
build-deploy-operator:
	$(MAKE) -C operator release-local.$(ARCH) RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C operator push-local RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C operator release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) CHART_APPVERSION=$(TAG) ARCH=$(ARCH)

	kubectl apply -f operator/config/rendered/release.yaml

.PHONY: build-deploy-metrics-operator
build-deploy-metrics-operator:
	$(MAKE) -C metrics-operator release-local.$(ARCH) RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C metrics-operator push-local RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C metrics-operator release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) CHART_APPVERSION=$(TAG) ARCH=$(ARCH)

	kubectl apply -f metrics-operator/config/rendered/release.yaml

.PHONY: build-deploy-scheduler
build-deploy-scheduler:
	$(MAKE) -C scheduler release-local.$(ARCH) RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C scheduler push-local RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C scheduler release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) CHART_APPVERSION=$(TAG) ARCH=$(ARCH)
	kubectl create namespace keptn-lifecycle-toolkit-system --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f scheduler/config/rendered/release.yaml

.PHONY: build-deploy-certmanager
build-deploy-certmanager:
	$(MAKE) -C klt-cert-manager release-local.$(ARCH) RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C klt-cert-manager push-local RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C klt-cert-manager release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) CHART_APPVERSION=$(TAG) ARCH=$(ARCH)
	kubectl create namespace keptn-lifecycle-toolkit-system --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f klt-cert-manager/config/rendered/release.yaml

.PHONY: build-deploy-dev-environment
build-deploy-dev-environment: build-deploy-certmanager build-deploy-operator build-deploy-metrics-operator build-deploy-scheduler


include docs/Makefile

yamllint:
	@docker run --rm -t -v $(PWD):/data cytopia/yamllint:$(YAMLLINT_VERSION) .github docs
