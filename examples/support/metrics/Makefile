.PHONY: 1-deploy-app
1-deploy-app:
	kubectl apply -f deployment.yaml

.PHONY: 2-create-metric
2-create-metric:
	kubectl apply -f load-generator-job.yaml

.PHONY: 3-generate-load
3-generate-load:
	kubectl apply -f load-generator-job.yaml

.PHONY: 4-deploy-hpa
4-deploy-hpa:
	kubectl apply -f hpa.yaml

.PHONY: get-raw-metric
get-raw-metric:
	kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta2/namespaces/podtato-metrics/keptnmetrics.metrics.sh/cpu-throttling/cpu-throttling" | jq .

.PHONY: get-hpa-status
get-hpa-status:
	kubectl describe  horizontalpodautoscalers.autoscaling -n podtato-metrics podtato-metrics-hpa

.PHONY: cleanup
cleanup:
	kubectl delete namespace podtato-metrics
