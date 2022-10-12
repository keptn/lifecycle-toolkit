JAEGER_VERSION ?= v1.38.0

.PHONY: observability
observability:
	# Create Namespace and install Jaeger
	kubectl create namespace observability --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply -f https://github.com/jaegertracing/jaeger-operator/releases/download/$(JAEGER_VERSION)/jaeger-operator.yaml -n observability
	kubectl wait --for=condition=available deployment/jaeger-operator -n observability --timeout=120s
	kubectl apply -f observability/config/jaeger.yaml

	# Install Prometheus
	kubectl create namespace monitoring --dry-run=client -o yaml | kubectl apply -f -
	kubectl apply --server-side -f observability/config/prometheus/setup
	kubectl apply -f observability/config/prometheus/

	# Install OpenTelemetry Collector
	kubectl apply -f observability/config/otel-collector.yaml

	# Wait for Resources

	kubectl wait --for=condition=available deployment/jaeger -n keptn-lifecycle-controller-system --timeout=120s
	kubectl wait --for=condition=available deployment/prometheus-operator -n monitoring --timeout=120s
	kubectl wait --for=condition=available deployment/prometheus-adapter -n monitoring --timeout=120s
	kubectl wait --for=condition=available deployment/kube-state-metrics -n monitoring --timeout=120s
	kubectl wait --for=condition=available deployment/otel-collector -n keptn-lifecycle-controller-system --timeout=120s
	kubectl wait --for=condition=available deployment/grafana -n monitoring --timeout=120s

	# Restart keptn-lifecycle-controller
	kubectl rollout restart deployment klc-controller-manager -n keptn-lifecycle-controller-system
	kubectl rollout status deployment -n keptn-lifecycle-controller-system klc-controller-manager --watch
	kubectl rollout restart deployment keptn-scheduler -n keptn-lifecycle-controller-system
	kubectl rollout status deployment -n keptn-lifecycle-controller-system keptn-scheduler --watch

.PHONY: observability-port-forward-jaeger
observability-port-forward-jaeger:
	echo "Open Jaeger in your Browser: http://localhost:16686"
	kubectl port-forward -n keptn-lifecycle-controller-system svc/jaeger-query 16686

.PHONY: observability-port-forward-prometheus
observability-port-forward-prometheus:
	echo "Open Prometheus in your Browser: http://localhost:9090"
	kubectl -n monitoring port-forward svc/prometheus-k8s 9090


.PHONY: undeploy-observability
undeploy-observability:
	kubectl delete -f observability/config/jaeger.yaml -n keptn-lifecycle-controller-system --ignore-not-found=true
	kubectl delete -f https://github.com/jaegertracing/jaeger-operator/releases/download/$(JAEGER_VERSION)/jaeger-operator.yaml -n observability --ignore-not-found=true
	kubectl delete -f observability/config/prometheus/ --ignore-not-found=true
	kubectl delete -f observability/config/prometheus/setup --ignore-not-found=true
	kubectl delete ns observability --ignore-not-found=true
	kubectl delete ns monitoring --ignore-not-found=true



.PHONY: deploy-podtatohead
deploy-podtatohead:
	kubectl apply -f podtatohead-deployment/.
	kubectl wait --for=condition=available deployment/podtato-head-entry -n podtato-kubectl --timeout=120s
	echo "Watch Workload Deployment Progress: kubectl get keptnworkloadinstances -n podtato-kubectl -w"
