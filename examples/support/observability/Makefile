# renovate: datasource=github-tags depName=jaegertracing/jaeger-operator
JAEGER_VERSION ?= v1.45.0
TOOLKIT_NAMESPACE ?= keptn-lifecycle-toolkit-system
PODTATO_NAMESPACE ?= podtato-kubectl
GRAFANA_PORT_FORWARD ?= 3000

.PHONY: install
install: install-cert-manager
	@echo "-----------------------------------"
	@echo "Create Namespace and install Jaeger"
	@echo "-----------------------------------"
	kubectl create namespace observability --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f https://github.com/jaegertracing/jaeger-operator/releases/download/$(JAEGER_VERSION)/jaeger-operator.yaml -n observability
	kubectl wait --for=condition=available deployment/jaeger-operator -n observability --timeout=300s
	kubectl apply -f config/jaeger.yaml -n "$(TOOLKIT_NAMESPACE)"

	@echo ""
	@echo "------------------"
	@echo "Configure Prometheus"
	@echo "------------------"
	kubectl create namespace monitoring --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply --server-side -f config/prometheus/setup
	kubectl wait --for=condition=Established --all CustomResourceDefinition --namespace=monitoring

	@echo ""
	@echo "-------------------------------"
	@echo "Install OpenTelemetry Collector"
	@echo "-------------------------------"
	kubectl apply -f config/otel-collector.yaml -n "$(TOOLKIT_NAMESPACE)"

	@echo ""
	@echo "-------------------------------"
	@echo "Wait for Resources to get ready"
	@echo "-------------------------------"
	kubectl wait --for=condition=available deployment/jaeger -n "$(TOOLKIT_NAMESPACE)" --timeout=120s
	kubectl wait --for=condition=available deployment/otel-collector -n "$(TOOLKIT_NAMESPACE)" --timeout=120s

	@echo ""
	@echo "------------------"
	@echo "Install Prometheus"
	@echo "------------------"
	kubectl apply -f config/prometheus/
	kubectl wait --for=condition=available deployment/prometheus-operator -n monitoring --timeout=120s
	kubectl wait --for=condition=available deployment/prometheus-adapter -n monitoring --timeout=120s
	kubectl wait --for=condition=available deployment/kube-state-metrics -n monitoring --timeout=120s
	kubectl wait --for=condition=available deployment/grafana -n monitoring --timeout=120s

	# @echo ""
	# @echo "------------------"
	# @echo "Install Prometheus Mockserver"
	# @echo "------------------"
	# kubectl apply -n "$(TOOLKIT_NAMESPACE)" -f config/prometheus-mock/mockserver-config.yaml
	# helm repo add mockserver "https://www.mock-server.com"
	# helm upgrade --install -n "$(TOOLKIT_NAMESPACE)" --version 5.13.0 mockserver mockserver/mockserver --set service.type=ClusterIP --wait


	@echo ""
	@echo "#######################################################"
	@echo "Observability Demo installed"
	@echo "- Install PodTatoHead: make deploy-podtatohead"
	@echo "- Port-Forward Grafana: make port-forward-grafana"
	@echo "#######################################################"

.PHONY: install-cert-manager
install-cert-manager:
	@echo "-----------------------------------"
	@echo "Create Namespace and install CertManager"
	@echo "-----------------------------------"
	kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.11.0/cert-manager.yaml
	kubectl wait --for=condition=Available deployment/cert-manager-webhook -n cert-manager --timeout=60s

.PHONY: port-forward-jaeger
port-forward-jaeger:
	@echo ""
	@echo "Open Jaeger in your Browser: http://localhost:16686"
	@echo "CTRL-c to stop port-forward"

	@echo ""
	kubectl port-forward -n "$(TOOLKIT_NAMESPACE)" svc/jaeger-query 16686


.PHONY: port-forward-prometheus
port-forward-prometheus:
	@echo ""
	@echo "Open Prometheus in your Browser: http://localhost:9090"
	@echo "CTRL-c to stop port-forward"
	@echo ""
	kubectl -n monitoring port-forward svc/prometheus-k8s 9090

.PHONY: port-forward-grafana
port-forward-grafana:
	@echo ""
	@echo "Open Grafana in your Browser: http://localhost:$(GRAFANA_PORT_FORWARD)"
	@echo "CTRL-c to stop port-forward"
	@echo ""
	kubectl -n monitoring port-forward svc/grafana $(GRAFANA_PORT_FORWARD):3000

.PHONY: import-grafana-dashboard
import-grafana-dashboard:
	@echo ""
	@echo "Importing Grafana Dashboard and Datasources"
	@echo ""
	curl -X POST --insecure http://admin:admin@localhost:$(GRAFANA_PORT_FORWARD)/api/dashboards/import -H "Content-Type: application/json" -d @./assets/grafana_dashboard_overview.json
	curl -X POST --insecure http://admin:admin@localhost:$(GRAFANA_PORT_FORWARD)/api/dashboards/import -H "Content-Type: application/json" -d @./assets/grafana_dashboard_workloads.json
	curl -X POST --insecure http://admin:admin@localhost:$(GRAFANA_PORT_FORWARD)/api/dashboards/import -H "Content-Type: application/json" -d @./assets/grafana_dashboard_applications.json

.PHONY: uninstall
uninstall:
	kubectl delete -f config/jaeger.yaml -n "$(TOOLKIT_NAMESPACE)" --ignore-not-found=true
	kubectl delete -f https://github.com/jaegertracing/jaeger-operator/releases/download/$(JAEGER_VERSION)/jaeger-operator.yaml -n observability --ignore-not-found=true
	#kubectl delete -f config/prometheus/ --ignore-not-found=true
	#kubectl delete -f config/prometheus/setup --ignore-not-found=true
	kubectl delete ns observability --ignore-not-found=true
	kubectl delete ns monitoring --ignore-not-found=true
	@echo ""
	@echo "##########################"
	@echo "Observability Demo deleted"
	@echo "##########################"

.PHONY: deploy-podtatohead
deploy-podtatohead:
	kubectl create namespace "$(PODTATO_NAMESPACE)" --dry-run=client -o yaml | kubectl apply -f -
	kubectl annotate namespace "$(PODTATO_NAMESPACE)" keptn.sh/lifecycle-toolkit=enabled --overwrite
	kubectl apply -f ./assets/podtatohead-deployment-evaluation/.
	kubectl apply -f ./assets/podtatohead-evaluationdefinitions/keptnevaluationdefinition1.yaml
	kubectl apply -f ./assets/podtatohead-evaluationdefinitions/keptnevaluationdefinition2.yaml
	kubectl wait --for=condition=available deployment/podtato-head-entry -n podtato-kubectl --timeout=120s
	@echo "Watch Workload Deployment Progress: kubectl get keptnworkloadinstances -n podtato-kubectl -w"

.PHONY: deploy-podtatohead-failed-post
deploy-podtatohead-failed-post:
	kubectl create namespace "$(PODTATO_NAMESPACE)" --dry-run=client -o yaml | kubectl apply -f -
	kubectl annotate namespace "$(PODTATO_NAMESPACE)" keptn.sh/lifecycle-toolkit=enabled --overwrite
	kubectl apply -f ./assets/podtatohead-deployment-evaluation/.
	kubectl apply -f ./assets/podtatohead-evaluationdefinitions/keptnevaluationdefinition1.yaml
	kubectl apply -f ./assets/podtatohead-evaluationdefinitions/keptnevaluationdefinition3.yaml
	kubectl wait --for=condition=available deployment/podtato-head-entry -n podtato-kubectl --timeout=120s
	@echo "Watch Workload Deployment Progress: kubectl get keptnworkloadinstances -n podtato-kubectl -w"

.PHONY: deploy-podtatohead-failed-pre
deploy-podtatohead-failed-pre:
	kubectl create namespace "$(PODTATO_NAMESPACE)" --dry-run=client -o yaml | kubectl apply -f -
	kubectl annotate namespace "$(PODTATO_NAMESPACE)" keptn.sh/lifecycle-toolkit=enabled --overwrite
	kubectl apply -f ./assets/podtatohead-deployment-evaluation/.
	kubectl apply -f ./assets/podtatohead-evaluationdefinitions/keptnevaluationdefinition3.yaml
	kubectl apply -f ./assets/podtatohead-evaluationdefinitions/keptnevaluationdefinition4.yaml
	kubectl wait --for=condition=available deployment/podtato-head-entry -n podtato-kubectl --timeout=120s
	@echo "Watch Workload Deployment Progress: kubectl get keptnworkloadinstances -n podtato-kubectl -w"

.PHONY: undeploy-podtatohead
undeploy-podtatohead:
	kubectl delete ns "$(PODTATO_NAMESPACE)" --ignore-not-found=true

	@echo "######################"
	@echo "PodTatoHead undeployed"
	@echo "######################"
