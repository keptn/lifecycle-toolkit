JAEGER_VERSION ?= v1.38.0
LFC_NAMESPACE ?= keptn-lifecycle-controller-system
PODTATO_NAMESPACE ?= podtato-kubectl
REPO_LOCAL_PATH ?= clusters/demo-cluster
GITHUB_REPO ?= ""
GITHUB_BRANCH ?= "main"
GITHUB_USER ?= ""


.PHONY: install
install: check-repo-url-empty check-github-user-empty check-flux-binary
	@echo "----------------------------------------"
	@echo "Preparing your flux repo and set up flux"
	@echo "----------------------------------------"
	@echo "(i) Make sure that you have a Personal Access Token prepared for your GitHub account (https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)"
	@echo "(i) You will need all permissions under repo (according to https://fluxcd.io/flux/installation/#github-and-github-enterprise)"
	@echo ""
	@sleep 5
	mkdir -p $(REPO_LOCAL_PATH) | true
	flux bootstrap github \
		--owner=$(GITHUB_USER) \
		--repository=$(GITHUB_REPO) \
		--branch=$(GITHUB_BRANCH) \
		--path=$(REPO_LOCAL_PATH) \
		--personal

.PHONY: manifests
manifests: check-flux-binary
	@echo "----------------------------------------"
	@echo "Generating manifests for flux"
	@echo "----------------------------------------"
	flux create source git podtato-head \
 	 --url=https://github.com/keptn-sandbox/lifecycle-controller \
	 --branch=main \
 	 --interval=30s \
  	 --export > $(REPO_LOCAL_PATH)/lfc-source.yaml
	cp -r ../podtatohead-deployment/. $(REPO_LOCAL_PATH)/




.PHONY: check-repo-url-empty
check-repo-url-empty:
	@if [ $(GITHUB_REPO) == "" ]; then \
		echo "ERROR: GITHUB_REPO has not been specified"; \
		exit 1; \
	fi

.PHONY: check-github-user-empty
check-github-user-empty:
	@if [ $(GITHUB_USER) == "" ]; then \
		echo "ERROR: GITHUB_USER has not been specified"; \
		exit 1; \
	fi

.PHONY: check-flux-binary
check-flux-binary:
	@if ! command -v flux &> /dev/null; then \
		echo "ERROR: flux binary not found. Please install flux"; \
		exit 1; \
	fi