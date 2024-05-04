// Package v1alpha1 contains API Schema definitions for the lifecycle v1alpha1 API group
// +groupName=lifecycle.keptn.sh
// +versionName=v1alpha1
package common

import (
	"errors"

	"go.opentelemetry.io/otel/propagation"
)

var ErrCannotCastKeptnAppVersion = errors.New("cannot cast KeptnAppVersion to v1")
var ErrCannotCastKeptnApp = errors.New("cannot cast KeptnApp to v1")

// KeptnState  is a string containing current Phase state  (Progressing/Succeeded/Failed/Unknown/Pending/Cancelled)
type KeptnState string

type CheckType string

// PhaseTraceID is a map storing TraceIDs of OpenTelemetry spans in lifecycle phases
type PhaseTraceID map[string]propagation.MapCarrier
