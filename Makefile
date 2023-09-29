# Image URL to use all building/pushing image targets

# renovate: datasource=github-tags depName=kubernetes-sigs/kustomize
KUSTOMIZE_VERSION?=v5.1.1
CHART_APPVERSION ?= v0.8.2 # x-release-please-version

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

.PHONY: integration-test #these tests should run on a real cluster!
integration-test:	# to run a single test by name use --test eg. --test=expose-keptn-metric
	kubectl kuttl test --start-kind=false ./test/integration/ --config=kuttl-test.yaml
	kubectl kuttl test --start-kind=false ./test/testmetrics/ --config=kuttl-test.yaml
	kubectl kuttl test --start-kind=false ./test/testanalysis/ --config=kuttl-test.yaml
	kubectl kuttl test --start-kind=false ./test/testcertificate/ --config=kuttl-test.yaml

.PHONY: integration-test-local #these tests should run on a real cluster!
integration-test-local: install-prometheus
	kubectl kuttl test --start-kind=false ./test/integration/ --config=kuttl-test-local.yaml
	kubectl kuttl test --start-kind=false ./test/testmetrics/ --config=kuttl-test-local.yaml
	kubectl kuttl test --start-kind=false ./test/testanalysis/ --config=kuttl-test-local.yaml
	kubectl kuttl test --start-kind=false ./test/testcertificate/ --config=kuttl-test-local.yaml

.PHONY: integration-test-scheduling-gates #these tests should run on a real cluster!
integration-test-scheduling-gates:	# to run a single test by name use --test eg. --test=expose-keptn-metric
	kubectl kuttl test --start-kind=false ./test/scheduling-gates/ --config=kuttl-test.yaml

.PHONY: integration-test-scheduling-gates-local #these tests should run on a real cluster!
integration-test-scheduling-gates-local: install-prometheus
	kubectl kuttl test --start-kind=false ./test/scheduling-gates/ --config=kuttl-test-local.yaml

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

.PHONY: metrics-operator-test
metrics-operator-test:
	$(MAKE) -C metrics-operator test

.PHONY: certmanager-test
certmanager-test:
	$(MAKE) -C klt-cert-manager test

.PHONY: operator-test
operator-test:
	$(MAKE) -C lifecycle-operator test

.PHONY: scheduler-test
scheduler-test:
	$(MAKE) -C scheduler test

#command(make test) to run all tests 
.PHONY: test
test: metrics-operator-test certmanager-test operator-test scheduler-test integration-test

.PHONY: cleanup-manifests
cleanup-manifests:
	rm -rf manifests

KUSTOMIZE_INSTALL_SCRIPT ?= "https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh"
.PHONY: kustomize
kustomize: $(KUSTOMIZE) ## Download kustomize locally if necessary.
$(KUSTOMIZE): $(LOCALBIN)
	test -s $(LOCALBIN)/kustomize || { curl -s $(KUSTOMIZE_INSTALL_SCRIPT) | bash -s -- $(subst v,,$(KUSTOMIZE_VERSION)) $(LOCALBIN); }

.PHONY: build-deploy-operator
build-deploy-operator:
	$(MAKE) -C lifecycle-operator release-local.$(ARCH) RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C lifecycle-operator push-local RELEASE_REGISTRY=$(RELEASE_REGISTRY) TAG=$(TAG)
	$(MAKE) -C lifecycle-operator release-manifests RELEASE_REGISTRY=$(RELEASE_REGISTRY) CHART_APPVERSION=$(TAG) ARCH=$(ARCH)

	kubectl apply -f lifecycle-operator/config/rendered/release.yaml

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
	@docker run --rm -t -v $(PWD):/data cytopia/yamllint:$(YAMLLINT_VERSION) .

##Run lint for the subfiles
.PHONY: metrics-operator-lint
metrics-operator-lint:
	$(MAKE) -C metrics-operator lint

.PHONY: certmanager-lint
certmanager-lint:
	$(MAKE) -C klt-cert-manager lint

.PHONY: operator-lint
operator-lint:
	$(MAKE) -C lifecycle-operator lint

.PHONY: scheduler-lint
scheduler-lint:
	$(MAKE) -C scheduler lint

.PHONY: lint
lint: 
	go install -v github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	metrics-operator-lint certmanager-lint operator-lint scheduler-lint
