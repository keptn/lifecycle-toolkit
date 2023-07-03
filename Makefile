# Image URL to use all building/pushing image targets

# renovate: datasource=github-tags depName=kubernetes-sigs/kustomize
KUSTOMIZE_VERSION?=v5.0.3
# renovate: datasource=github-tags depName=helm/helm
HELM_VERSION ?= v3.12.1
CHART_APPVERSION ?= v0.8.0 # x-release-please-version

# renovate: datasource=docker depName=cytopia/yamllint
YAMLLINT_VERSION ?= alpine

# RELEASE_REGISTRY is the container registry to push
# into.
RELEASE_REGISTRY?=ghcr.io/keptn
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
	test -s $(LOCALBIN)/helmify || GOBIN=$(LOCALBIN) go install github.com/keptn/helmify/cmd/helmify@1060b5d08806e40bfd9f38c3e8a9a302ab38e71a

.PHONY: integration-test #these tests should run on a real cluster!
integration-test:	# to run a single test by name use --test eg. --test=expose-keptn-metric
	kubectl kuttl test --start-kind=false ./test/integration/ --config=kuttl-test.yaml
	kubectl kuttl test --start-kind=false ./test/testcertificate/ --config=kuttl-test.yaml



.PHONY: integration-test-local #these tests should run on a real cluster!
integration-test-local: install-prometheus
	kubectl kuttl test --start-kind=false ./test/integration/ --config=kuttl-test-local.yaml
	kubectl kuttl test --start-kind=false ./test/testcertificate/ --config=kuttl-test-local.yaml

.PHONY: load-test
load-test:
	kubectl apply -f ./test/load/assets/templates/namespace.yaml
	kubectl apply -f ./test/load/assets/templates/provider.yaml
	kube-burner init -c ./test/load/cfg.yml --metrics-profile ./test/load/metrics.yml

.PHONY: install-prometheus
install-prometheus:
	kubectl create namespace monitoring --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply --server-side -f test/prometheus/setup
	kubectl wait --for=condition=Established --all CustomResourceDefinition --namespace=monitoring
	kubectl apply -f test/prometheus/
	kubectl wait --for=condition=available deployment/prometheus-operator -n monitoring --timeout=120s
	kubectl wait --for=condition=available deployment/prometheus-adapter -n monitoring --timeout=120s
	kubectl wait --for=condition=available deployment/kube-state-metrics -n monitoring --timeout=120s
	kubectl wait pod/prometheus-k8s-0 --for=condition=ready --timeout=120s -n monitoring



.PHONY: cleanup-manifests
cleanup-manifests:
	rm -rf manifests

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	test -s $(LOCALBIN)/kustomize || { curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN); }

.PHONY: build-deploy-helm
build-deploy-helm: build-dev-environment
	helm upgrade --install klt ./helm/chart \
		--set metrics-operator.metricsOperator.manager.image.repository=$(RELEASE_REGISTRY)/metrics-operator \
		--set metrics-operator.metricsOperator.manager.image.tag=$(TAG) \
		--set lifecycle-operator.operator.manager.image.repository=$(RELEASE_REGISTRY)/lifecycle-operator \
		--set lifecycle-operator.operator.manager.image.tag=$(TAG) \
		--set cert-manager.certificateOperator.manager.image.repository=$(RELEASE_REGISTRY)/certificate-operator \
		--set cert-manager.certificateOperator.manager.image.tag=$(TAG) \
		--set scheduler.scheduler.scheduler.image.repository=$(RELEASE_REGISTRY)/scheduler \
		--set scheduler.scheduler.scheduler.image.tag=$(TAG) \

.PHONY: build-deploy-operator
build-deploy-operator: build-operator
	kubectl apply -f operator/config/rendered/release.yaml

.PHONY: build-operator
build-operator:
	$(MAKE) -C operator release-local.$(ARCH) RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C operator push-local RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C operator release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) CHART_APPVERSION=$(TAG) ARCH=$(ARCH)

.PHONY: build-deploy-metrics-operator
build-deploy-metrics-operator: build-metrics-operator
	kubectl apply -f metrics-operator/config/rendered/release.yaml

.PHONY: build-metrics-operator
build-metrics-operator:
	$(MAKE) -C metrics-operator release-local.$(ARCH) RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C metrics-operator push-local RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C metrics-operator release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) CHART_APPVERSION=$(TAG) ARCH=$(ARCH)

.PHONY: build-deploy-scheduler
build-deploy-scheduler: build-scheduler
	kubectl apply -f scheduler/config/rendered/release.yaml

.PHONY: build-scheduler
build-scheduler:
	$(MAKE) -C scheduler release-local.$(ARCH) RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C scheduler push-local RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C scheduler release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) CHART_APPVERSION=$(TAG) ARCH=$(ARCH)

.PHONY: build-deploy-certmanager
build-deploy-certmanager: build-certmanager
	kubectl apply -f klt-cert-manager/config/rendered/release.yaml

.PHONY: build-certmanager
build-certmanager:
	$(MAKE) -C klt-cert-manager release-local.$(ARCH) RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C klt-cert-manager push-local RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C klt-cert-manager release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) CHART_APPVERSION=$(TAG) ARCH=$(ARCH)

.PHONY: build-deploy-dev-environment
build-deploy-dev-environment: build-deploy-certmanager build-deploy-operator build-deploy-metrics-operator build-deploy-scheduler

.PHONY: build-dev-environment
build-dev-environment: build-certmanager build-operator build-metrics-operator build-scheduler


include docs/Makefile

yamllint:
	@docker run --rm -t -v $(PWD):/data cytopia/yamllint:$(YAMLLINT_VERSION) .
