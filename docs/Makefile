ROOT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
VOLUMES := -v $(ROOT_DIR):/src 
# renovate: datasource=docker depName=klakegg/hugo
HUGO_VERSION := 0.107.0
IMAGE := klakegg/hugo:$(HUGO_VERSION)-ext
PORT := 1314
DOCKER_CMD := docker run --rm -t -e HUGO_CACHEDIR=/src/tmp/.hugo -e HUGOxPARAMSxGITHUB_REPO=""

.PHONY: build server clean shell

build:
	$(DOCKER_CMD) $(VOLUMES) $(IMAGE) -D -v

shell:
	$(DOCKER_CMD) -i $(VOLUMES) $(IMAGE) shell

server:
	$(DOCKER_CMD) $(VOLUMES) -p  $(PORT):$(PORT) $(IMAGE) server -D -p $(PORT)

clean:
	$(DOCKER_CMD) $(VOLUMES) $(IMAGE) --cleanDestinationDir

.PHONY: htmltest

# renovate: datasource=docker depName=wjdp/htmltest
HTMLTEST_VERSION := v0.17.0
htmltest: build
	$(DOCKER_CMD) -v $(ROOT_DIR):/test wjdp/htmltest:$(HTMLTEST_VERSION) -c .htmltest.yml public

.PHONY: lint lint-fix
lint: markdownlint
lint: lint-fix

# Markdown lint configuration
#
# - .markdownlintignore holds the configuration for files to be ignored
# - .markdownlint.yaml contains the rules for markdownfiles
# renovate: datasource=docker depName=davidanson/markdownlint-cli2-rules
MDL_DOCKER_VERSION := v0.6.0 # TODO: remove next line on version increase from v0.6.0
MDL_DOCKER_VERSION := next
MDL_CMD := docker run -v $(ROOT_DIR)../:/workdir --rm 

.PHONY: markdownlint markdownlint-fix
markdownlint:
	$(MDL_CMD) davidanson/markdownlint-cli2-rules:${MDL_DOCKER_VERSION} "**/*.md" 

markdownlint-fix:
	$(MDL_CMD) --entrypoint="markdownlint-cli2-fix" davidanson/markdownlint-cli2-rules:${MDL_DOCKER_VERSION} "**/*.md" 

# Define the directories
OPERATOR_DIR := operator
CONFIG_DIR := $(OPERATOR_DIR)/config
HELM_TEMPLATES_DIR := helm/chart/templates
METRICS_OPERATOR_DIR := metrics-operator

.PHONY: resources

# New target to automate the steps for both operators
resources: generate-manifests copy-manifests update-crd-docs

generate-manifests:
	@echo "Generating manifests..."
	cd $(OPERATOR_DIR) && make generate && make manifests
	cd $(METRICS_OPERATOR_DIR) && make generate && make manifests

copy-manifests:
	@echo "Copying manifests to Helm..."
	cp $(CONFIG_DIR)/crd/bases/*.yaml $(HELM_TEMPLATES_DIR)/
	cp $(METRICS_OPERATOR_DIR)/config/crd/bases/*.yaml $(HELM_TEMPLATES_DIR)/
	@for file in $(HELM_TEMPLATES_DIR)/*.yaml; do \
		tmpfile="$${file}.tmp"; \
		awk '/metadata:/{print;print "  labels:\n    app.kubernetes.io/part-of: keptn-lifecycle-toolkit\n    crdGroup: metrics.keptn.sh\n    keptn.sh/inject-cert: \"true\"";next}1' $$file > $$tmpfile; \
		awk '/metadata:/{print;print "  annotations:\n    cert-manager.io/inject-ca-from: \"{{ .Release.Namespace }}/klt-certs\"";next}1' $$tmpfile > $$file; \
		rm $$tmpfile; \
	done

update-crd-docs:
	@echo "Updating CRD docs..."
	cd .github/scripts/generate-crd-docs/ && ./generate-crd-docs.sh
