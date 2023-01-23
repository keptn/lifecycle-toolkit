package common

import (
	"github.com/pkg/errors"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const KeptnMetricProviderName = "keptn-metric"

var ErrForbiddenProvider = errors.New("Forbidden! KeptnMetrics should define a provider different from keptn-metric")
