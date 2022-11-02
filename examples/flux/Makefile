JAEGER_VERSION ?= v1.38.0
LFC_NAMESPACE ?= keptn-lifecycle-toolkit-system
PODTATO_NAMESPACE ?= podtato-kubectl
REPO_PATH ?= flux-demo
REPO_CLUSTER_PATH ?= clusters/demo-cluster
GITHUB_REPO ?= ""
GITHUB_BRANCH ?= "main"
GITHUB_USER ?= ""
GITHUB_TOKEN ?= ""


.PHONY: install
install: check-repo-url-empty check-github-user-empty check-flux-binary
	@echo "----------------------------------------"
	@echo "Preparing your flux repo and set up flux"
	@echo "----------------------------------------"
	@echo "(i) Make sure that you have a Personal Access Token prepared for your GitHub account (https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)"
	@echo "(i) You will need all permissions under repo (according to https://fluxcd.io/flux/installation/#github-and-github-enterprise)"
	@echo ""
	@sleep 5
	flux bootstrap github \
		--owner=$(GITHUB_USER) \
		--repository=$(GITHUB_REPO) \
		--branch=$(GITHUB_BRANCH) \
		--path=$(REPO_CLUSTER_PATH) \
		--personal
	rm -rf $(REPO_CLUSTER_PATH)

	@echo "----------------------------------------"
	@echo "Cloning your flux repo"
	@echo "----------------------------------------"
	git clone $(GITHUB_REPO) $(REPO_PATH)

	@echo "-----------------------------------------------------------------------------"
	@echo "Flux is prepared, you can find the repo at $(REPO_PATH)"
	@echo "Then you can proceed with the next step: make manifests <YOUR_REPO_PATH>"
	@echo "-----------------------------------------------------------------------------"

.PHONY: manifests
manifests:
	cp config/flux-config.yaml $(REPO_PATH)/$(REPO_CLUSTER_PATH)
	cd $(REPO_PATH) && git add . && git commit -m "Add podtatohead manifests" && git push

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

.PHONY: uninstall
uninstall: check-flux-binary
	@echo "----------------------------------------"
	@echo "Uninstalling flux"
	@echo "----------------------------------------"
	flux uninstall