JAEGER_VERSION ?= v1.38.0
LFC_NAMESPACE ?= keptn-lifecycle-controller-system
PODTATO_NAMESPACE ?= podtato-kubectl

.PHONY: install
install:
	@echo "-----------------------------------"
	@echo "Create Namespace and install Jaeger"
	@echo "-----------------------------------"
	kubectl create namespace observability --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f https://github.com/jaegertracing/jaeger-operator/releases/download/$(JAEGER_VERSION)/jaeger-operator.yaml -n observability
	kubectl wait --for=condition=available deployment/jaeger-operator -n observability --timeout=300s
	kubectl apply -f config/jaeger.yaml -n "$(LFC_NAMESPACE)"

	@echo ""
	@echo "------------------"
	@echo "Install Prometheus"
	@echo "------------------"
	kubectl create namespace monitoring --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply --server-side -f config/prometheus/setup
	kubectl apply -f config/prometheus/

	@echo ""
	@echo "-------------------------------"
	@echo "Install OpenTelemetry Collector"
	@echo "-------------------------------"
	kubectl apply -f config/otel-collector.yaml -n "$(LFC_NAMESPACE)"

	@echo ""
	@echo "-------------------------------"
	@echo "Wait for Resources to get ready"
	@echo "-------------------------------"
	kubectl wait --for=condition=available deployment/jaeger -n "$(LFC_NAMESPACE)" --timeout=120s
	kubectl wait --for=condition=available deployment/prometheus-operator -n monitoring --timeout=120s
	kubectl wait --for=condition=available deployment/prometheus-adapter -n monitoring --timeout=120s
	kubectl wait --for=condition=available deployment/kube-state-metrics -n monitoring --timeout=120s
	kubectl wait --for=condition=available deployment/otel-collector -n "$(LFC_NAMESPACE)" --timeout=120s
	kubectl wait --for=condition=available deployment/grafana -n monitoring --timeout=120s

	@echo ""
	@echo "----------------------------------"
	@echo "Restart Keptn Lifecycle Controller"
	@echo "----------------------------------"
	kubectl rollout restart deployment -n "$(LFC_NAMESPACE)" klc-controller-manager
	kubectl rollout status deployment -n "$(LFC_NAMESPACE)" klc-controller-manager --watch
	kubectl rollout restart deployment -n "$(LFC_NAMESPACE)" keptn-scheduler -n keptn-lifecycle-controller-system
	kubectl rollout status deployment -n "$(LFC_NAMESPACE)" keptn-scheduler --watch

	@echo ""
	@echo "#######################################################"
	@echo "Observability Demo installed"
	@echo "- Install PodTatoHead: make deploy-podtatohead"
	@echo "- Port-Forward Jaeger: make port-forward-jaeger"
	@echo "- Port-Forward Prometheus: make port-forward-prometheus"
	@echo "#######################################################"

.PHONY: port-forward-jaeger
port-forward-jaeger:
	@echo ""
	@echo "Open Jaeger in your Browser: http://localhost:16686"
	@echo "CTRL-c to stop port-forward"

	@echo ""
	kubectl port-forward -n "$(LFC_NAMESPACE)" svc/jaeger-query 16686


.PHONY: port-forward-prometheus
port-forward-prometheus:
	@echo ""
	@echo "Open Prometheus in your Browser: http://localhost:9090"
	@echo "CTRL-c to stop port-forward"
	@echo ""
	kubectl -n monitoring port-forward svc/prometheus-k8s 9090



.PHONY: uninstall
uninstall:
	kubectl delete -f config/jaeger.yaml -n "$(LFC_NAMESPACE)" --ignore-not-found=true
	kubectl delete -f https://github.com/jaegertracing/jaeger-operator/releases/download/$(JAEGER_VERSION)/jaeger-operator.yaml -n observability --ignore-not-found=true
	kubectl delete -f config/prometheus/ --ignore-not-found=true
	kubectl delete -f config/prometheus/setup --ignore-not-found=true
	kubectl delete ns observability --ignore-not-found=true
	kubectl delete ns monitoring --ignore-not-found=true
	@echo ""
	@echo "##########################"
	@echo "Observability Demo deleted"
	@echo "##########################"

.PHONY: deploy-podtatohead
deploy-podtatohead:
	kubectl create namespace "$(PODTATO_NAMESPACE)" --dry-run=client -o yaml | kubectl apply -f -
	kubectl annotate namespace "$(PODTATO_NAMESPACE)" keptn.sh/enabled=true --overwrite
	kubectl apply -f ../podtatohead-deployment/.
	kubectl wait --for=condition=available deployment/podtato-head-entry -n podtato-kubectl --timeout=120s
	@echo "Watch Workload Deployment Progress: kubectl get keptnworkloadinstances -n podtato-kubectl -w"

.PHONY: undeploy-podtatohead
undeploy-podtatohead:
	kubectl delete ns "$(PODTATO_NAMESPACE)" --ignore-not-found=true

	@echo "######################"
	@echo "PodTatoHead undeployed"
	@echo "######################"