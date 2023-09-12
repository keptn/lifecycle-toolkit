# renovate: datasource=github-tags depName=jaegertracing/jaeger-operator
JAEGER_VERSION ?= v1.49.0
TOOLKIT_NAMESPACE ?= keptn-lifecycle-toolkit-system
PODTATO_NAMESPACE ?= podtato-kubectl
GRAFANA_PORT_FORWARD ?= 3000

.PHONY: install
install: install-observability
	@echo "-----------------------------------"
	@echo "Create Namespace and install Keptn-lifecycle-toolkit"
	@echo "-----------------------------------"
	helm repo add klt https://charts.lifecycle.keptn.sh
	helm repo update
	helm upgrade --install keptn klt/klt -n $(TOOLKIT_NAMESPACE) --create-namespace --wait
	kubectl apply -f support/keptn/keptnconfig.yaml -n $(TOOLKIT_NAMESPACE)

.PHONY: install-observability
install-observability:
	kubectl create namespace $(TOOLKIT_NAMESPACE) --dry-run=client -o yaml | kubectl apply -f -
	make -C support/observability install

.PHONY: install-argo
install-argo:
	make -C support/argo install

.PHONY: port-forward-jaeger
port-forward-jaeger:
	make -C support/observability port-forward-jaeger

.PHONY: port-forward-grafana
port-forward-grafana:
	make -C support/observability port-forward-grafana GRAFANA_PORT_FORWARD=$(GRAFANA_PORT_FORWARD)

.PHONY: deploy-version-1
deploy-version-1:
	kubectl create namespace "$(PODTATO_NAMESPACE)" --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -k sample-app/version-1

.PHONY: deploy-version-2
deploy-version-2:
	kubectl create namespace "$(PODTATO_NAMESPACE)" --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -k sample-app/version-2

.PHONY: deploy-version-3
deploy-version-3:
	kubectl create namespace "$(PODTATO_NAMESPACE)" --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -k sample-app/version-3

.PHONY: cleanup
cleanup:
	kubectl delete ns "$(PODTATO_NAMESPACE)" --ignore-not-found=true

	@echo "######################"
	@echo "PodTatoHead undeployed"
	@echo "######################"

.PHONY: undeploy-podtatohead
undeploy-podtatohead: cleanup


.PHONY: uninstall-observability
uninstall-observability: undeploy-podtatohead
	make -C support/observability uninstall

.PHONY: uninstall-argo
uninstall-argo: undeploy-podtatohead
	make -C support/argo uninstall

.PHONY: uninstall
uninstall: uninstall-observability uninstall-argo
	@echo "-----------------------------------"
	@echo "Uninstall Keptn-lifecycle-toolkit"
	@echo "-----------------------------------"
	helm uninstall keptn -n $(TOOLKIT_NAMESPACE)
	kubectl delete ns $(TOOLKIT_NAMESPACE) --ignore-not-found=true

.PHONY: port-forward-argocd
port-forward-argocd:
	@echo ""
	@echo "Open ArgoCD in your Browser: http://localhost:8080"
	@echo "CTRL-c to stop port-forward"

	@echo ""
	kubectl port-forward svc/argocd-server -n "$(ARGO_NAMESPACE)" 8080:443

.PHONY: argo-get-password
argo-get-password:
	@echo $(ARGO_SECRET)

.PHONY: restart-lifecycle-toolkit
restart-lifecycle-toolkit:
	@echo ""
	@echo "----------------------------------"
	@echo "Restart Keptn Lifecycle Controller"
	@echo "----------------------------------"
	kubectl rollout restart deployment -n "$(TOOLKIT_NAMESPACE)" -l control-plane=lifecycle-operator
	kubectl rollout status deployment -n "$(TOOLKIT_NAMESPACE)" -l control-plane=lifecycle-operator --watch
	kubectl rollout restart deployment -n "$(TOOLKIT_NAMESPACE)" -l component=scheduler
	kubectl rollout status deployment -n "$(TOOLKIT_NAMESPACE)" -l component=scheduler --watch

