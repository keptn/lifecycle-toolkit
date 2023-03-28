.PHONY: get-raw-metric
get-raw-metric:
	kubectl get --raw "/apis/custom.metrics.k8s.io/v1beta2/namespaces/podtato-kubectl/keptnmetrics.metrics.sh/cpu-throttling/cpu-throttling" | jq .

.PHONY: get-hpa-status
get-hpa-status:
	kubectl describe  horizontalpodautoscalers.autoscaling -n podtato-kubectl podtato-mint
