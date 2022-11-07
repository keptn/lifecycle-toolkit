JAEGER_VERSION ?= v1.38.0
LFC_NAMESPACE ?= keptn-lifecycle-controller-system
PODTATO_NAMESPACE ?= podtato-kubectl

.PHONY: install
install:
	make -C ../observability install
	make -C ../observability import-grafana-dashboard

.PHONY: port-forward-jaeger
port-forward-jaeger:
	make -C ../observability port-forward-jaeger

.PHONY: port-forward-grafana
port-forward-grafana:
	make -C ../observability port-forward-grafana

.PHONY: deploy-version-1
deploy-version-1:
	kubectl create namespace "$(PODTATO_NAMESPACE)" --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -k version-1

.PHONY: deploy-version-2
deploy-version-2:
	kubectl create namespace "$(PODTATO_NAMESPACE)" --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -k version-2

.PHONY: deploy-version-3
deploy-version-3:
	kubectl create namespace "$(PODTATO_NAMESPACE)" --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -k version-3

.PHONY: undeploy-podtatohead
undeploy-podtatohead:
	kubectl delete ns "$(PODTATO_NAMESPACE)" --ignore-not-found=true

	@echo "######################"
	@echo "PodTatoHead undeployed"
	@echo "######################"

.PHONY: uninstall
uninstall: undeploy-podtatohead
	make -C ../observability uninstall