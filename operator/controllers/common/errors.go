package common

import "fmt"

var ErrCannotWrapToPhaseItem = fmt.Errorf("provided object does not implement PhaseItem interface")
var ErrCannotWrapToListItem = fmt.Errorf("provided object does not implement ListItem interface")
var ErrRetryCountExceeded = fmt.Errorf("retryCount for evaluation exceeded")
var ErrNoValues = fmt.Errorf("no values")
var ErrInvalidOperator = fmt.Errorf("invalid operator")
var ErrCannotMarshalParams = fmt.Errorf("could not marshal parameters")

var ErrCannotRetrieveInstancesMsg = "could not retrieve instances: %w"
var ErrCannotFetchAppMsg = "could not retrieve KeptnApp: %w"
var ErrCannotFetchAppVersionMsg = "could not retrieve KeptnappVersion: %w"
var ErrCannotRetrievieWorkloadInstancesMsg = "could not retrieve KeptnWorkloadInstance: %w"
var ErrCannotRetrievieWorkloadMsg = "could not retrieve KeptnWorkload: %w"
var ErrNoLabelsFoundTask = "no labels found for task: %s"
var ErrNoConfigMapMsg = "No ConfigMap specified or HTTP source specified in TaskDefinition) / Namespace: %s, Name: %s"
var ErrCannotGetFunctionConfigMap = "could not get function configMap: %w"
var ErrCannotFetchAppVersionForWorkloadInstanceMsg = "could not fetch AppVersion for KeptnWorkloadInstance: "
