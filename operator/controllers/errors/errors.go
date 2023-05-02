package errors

import "fmt"

var ErrCannotWrapToPhaseItem = fmt.Errorf("provided object does not implement PhaseItem interface")
var ErrCannotWrapToListItem = fmt.Errorf("provided object does not implement ListItem interface")
var ErrCannotWrapToMetricsObject = fmt.Errorf("provided object does not implement MetricsObject interface")
var ErrCannotWrapToActiveMetricsObject = fmt.Errorf("provided object does not implement ActiveMetricsObject interface")
var ErrCannotWrapToEventObject = fmt.Errorf("provided object does not implement EventObject interface")
var ErrCannotWrapToSpanItem = fmt.Errorf("provided object does not implement SpanItem interface")
var ErrRetryCountExceeded = fmt.Errorf("retryCount for evaluation exceeded")
var ErrNoValues = fmt.Errorf("no values")
var ErrInvalidOperator = fmt.Errorf("invalid operator")
var ErrCannotMarshalParams = fmt.Errorf("could not marshal parameters")
var ErrUnsupportedWorkloadVersionResourceReference = fmt.Errorf("unsupported Resource Reference")
var ErrCannotGetKeptnTaskDefinition = fmt.Errorf("cannot retrieve KeptnTaskDefinition")

var ErrCannotRetrieveConfigMsg = "could not retrieve KeptnConfig: %w"
var ErrCannotRetrieveInstancesMsg = "could not retrieve instances: %w"
var ErrCannotFetchAppMsg = "could not retrieve KeptnApp: %w"
var ErrCannotFetchAppVersionMsg = "could not retrieve KeptnappVersion: %w"
var ErrCannotRetrieveWorkloadVersionsMsg = "could not retrieve KeptnWorkloadVersion: %w"
var ErrCannotRetrieveWorkloadMsg = "could not retrieve KeptnWorkload: %w"
var ErrNoLabelsFoundTask = "no labels found for task: %s"
var ErrNoConfigMapMsg = "No ConfigMap specified or HTTP source specified in TaskDefinition) / Namespace: %s, Name: %s"
var ErrCannotGetFunctionConfigMap = "could not get function configMap: %w"
var ErrCannotFetchAppVersionForWorkloadVersionMsg = "could not fetch AppVersion for KeptnWorkloadVersion: "
var ErrCouldNotUnbindSpan = "could not unbind span for %s"
