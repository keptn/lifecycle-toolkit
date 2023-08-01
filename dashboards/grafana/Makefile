GRAFANA_SCHEME ?= http
GRAFANA_HOSTNAME ?= localhost
GRAFANA_USERNAME ?= admin
GRAFANA_PASSWORD ?= admin
GRAFANA_PORT ?= 3000

.PHONY: generate
generate: generate-configmaps generate-import

.PHONY: generate-import
generate-import:
	mkdir import | true
	./utils/generate-import.sh ./import

.PHONY: generate-configmap
generate-configmaps:
	mkdir configmap | true
	kubectl create configmap grafana-dashboard-keptn-overview -n monitoring --from-file=./grafana_dashboard_overview.json -o yaml --dry-run=client > configmap/grafana-dashboard-keptn-overview.yaml
	kubectl create configmap grafana-dashboard-keptn-applications -n monitoring --from-file=./grafana_dashboard_applications.json -o yaml --dry-run=client > configmap/grafana-dashboard-keptn-applications.yaml
	kubectl create configmap grafana-dashboard-keptn-workloads -n monitoring --from-file=./grafana_dashboard_workloads.json -o yaml --dry-run=client > configmap/grafana-dashboard-keptn-workloads.yaml
	kubectl apply -f /configmap
	#cp configmap/* ../../examples/support/observability/config/prometheus

.PHONY: import
import:
	@curl -X POST -H "Content-Type: application/json" -d @import/grafana_dashboard_overview.json ${GRAFANA_SCHEME}://${GRAFANA_USERNAME}:${GRAFANA_PASSWORD}@${GRAFANA_HOSTNAME}:$(GRAFANA_PORT)/api/dashboards/db
	@curl -X POST -H "Content-Type: application/json" -d @import/grafana_dashboard_applications.json ${GRAFANA_SCHEME}://${GRAFANA_USERNAME}:${GRAFANA_PASSWORD}@${GRAFANA_HOSTNAME}:$(GRAFANA_PORT)/api/dashboards/db
	@curl -X POST -H "Content-Type: application/json" -d @import/grafana_dashboard_workloads.json ${GRAFANA_SCHEME}://${GRAFANA_USERNAME}:${GRAFANA_PASSWORD}@${GRAFANA_HOSTNAME}:$(GRAFANA_PORT)/api/dashboards/db

