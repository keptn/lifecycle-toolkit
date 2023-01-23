package common

import (
	"github.com/pkg/errors"
)

const KeptnMetricProviderName = "keptn-metric"

var ErrForbiddenProvider = errors.New("Forbidden! KeptnMetrics should define a provider different from keptn-metric")
