# renovate: datasource=github-tags depName=jaegertracing/jaeger
JAEGER_VERSION ?= v1.43.0
LFC_NAMESPACE ?= keptn-lifecycle-controller-system
PODTATO_NAMESPACE ?= podtato-kubectl
GRAFANA_PORT_FORWARD ?= 3000

.PHONY: install
install:
	make -C ../support/observability install

.PHONY: port-forward-jaeger
port-forward-jaeger:
	make -C ../support/observability port-forward-jaeger

.PHONY: port-forward-grafana
port-forward-grafana:
	make -C ../support/observability port-forward-grafana GRAFANA_PORT_FORWARD=$(GRAFANA_PORT_FORWARD)

.PHONY: deploy-version-1
deploy-version-1:
	kubectl apply -k version-1

.PHONY: deploy-version-2
deploy-version-2:
	kubectl apply -k version-2

.PHONY: deploy-version-3
deploy-version-3:
	kubectl apply -k version-3

.PHONY: undeploy-podtatohead
undeploy-podtatohead:
	kubectl delete ns "$(PODTATO_NAMESPACE)" --ignore-not-found=true

	@echo "######################"
	@echo "PodTatoHead undeployed"
	@echo "######################"

.PHONY: uninstall
uninstall: undeploy-podtatohead
	make -C ../support/observability uninstall
