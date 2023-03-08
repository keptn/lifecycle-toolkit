package provider

import "github.com/pkg/errors"

const (
	metricsGroup        = "metrics.keptn.sh"
	metricsResource     = "keptnmetrics"
	defaultMetricsValue = "0.0"
)

var ErrMetricNotFound = errors.New("no metric value found")
