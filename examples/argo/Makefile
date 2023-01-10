LFC_NAMESPACE ?= keptn-lifecycle-toolkit-system
PODTATO_NAMESPACE ?= podtato-kubectl
ARGO_NAMESPACE ?= argocd
# renovate: datasource=github-tags depName=argoproj/argo-cd
ARGO_VERSION ?= v2.5.5
# renovate: datasource=github-tags depName=cert-manager/cert-manager
CERT_MANAGER_VERSION ?= v1.10.1
ARGO_SECRET = $(shell kubectl -n ${ARGO_NAMESPACE} get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d; echo)

.PHONY: install
install:
	@echo "-----------------------------------"
	@echo "Create Namespace and install Keptn-lifecycle-toolkit"
	@echo "-----------------------------------"
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/$(CERT_MANAGER_VERSION)/cert-manager.yaml
	kubectl wait --for=condition=available deployment/cert-manager -n cert-manager --timeout=300s
	kubectl apply -f https://github.com/keptn/lifecycle-toolkit/releases/download/v0.5.0/manifest.yaml #x-release-please-version
	kubectl wait --for=condition=available deployment/klc-controller-manager -n keptn-lifecycle-toolkit-system --timeout=300s

	@echo "-----------------------------------"
	@echo "Create Namespace and install ArgoCD"
	@echo "-----------------------------------"
	kubectl create namespace $(ARGO_NAMESPACE)
	kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/$(ARGO_VERSION)/manifests/install.yaml

	@echo ""
	@echo "-------------------------------"
	@echo "Wait for Resources to get ready"
	@echo "-------------------------------"
	kubectl wait --for=condition=available deployment/argocd-dex-server -n "$(ARGO_NAMESPACE)" --timeout=120s
	kubectl wait --for=condition=available deployment/argocd-redis -n "$(ARGO_NAMESPACE)" --timeout=120s
	kubectl wait --for=condition=available deployment/argocd-repo-server -n "$(ARGO_NAMESPACE)" --timeout=120s
	kubectl wait --for=condition=available deployment/argocd-server  -n "$(ARGO_NAMESPACE)" --timeout=120s


	@echo ""
	@echo "#######################################################"
	@echo "ArgoCD Demo installed"
	@echo "- Get Admin Password: make argo-get-password"
	@echo "- Port-Forward ArgoCD: make port-forward-argocd"
	@echo "- Get Argo CLI: https://argo-cd.readthedocs.io/en/stable/cli_installation/"
	@echo "- Configure ArgoCD CLI (needs port-forward): make argo-configure-cli"
	@echo "- Install PodTatoHead via ArgoCD: make argo-install-podtatohead"	
	@echo "#######################################################"

.PHONY: argo-configure-cli
argo-configure-cli:
	@argocd login localhost:8080 --username admin --password $(ARGO_SECRET) --insecure

.PHONY: argo-get-password
argo-get-password:
	@echo $(ARGO_SECRET)

.PHONY: port-forward-argocd
port-forward-argocd:
	@echo ""
	@echo "Open ArgoCD in your Browser: http://localhost:8080"
	@echo "CTRL-c to stop port-forward"

	@echo ""
	kubectl port-forward svc/argocd-server -n "$(ARGO_NAMESPACE)" 8080:443

.PHONY: argo-install-podtatohead
argo-install-podtatohead:
	@echo ""
	kubectl apply -f config/app.yaml -n "$(ARGO_NAMESPACE)" 

.PHONY: uninstall
uninstall:
	kubectl delete -n $(ARGO_NAMESPACE) -f https://raw.githubusercontent.com/argoproj/argo-cd/v$(ARGO_VERSION)/manifests/install.yaml --ignore-not-found=true
	kubectl delete namespace $(ARGO_NAMESPACE) --ignore-not-found=true
	@echo ""
	@echo "##########################"
	@echo "Argo Demo deleted"
	@echo "##########################"
