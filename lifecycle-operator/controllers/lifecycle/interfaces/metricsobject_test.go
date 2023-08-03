package interfaces

import (
	"testing"
	"time"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/interfaces/fake"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
)

func TestMetricsObjectWrapper(t *testing.T) {
	appVersion := v1alpha3.KeptnAppVersion{
		Status: v1alpha3.KeptnAppVersionStatus{
			Status:       apicommon.StateFailed,
			CurrentPhase: "test",
		},
	}

	object, err := NewMetricsObjectWrapperFromClientObject(&appVersion)
	require.Nil(t, err)

	require.False(t, object.IsEndTimeSet())
}

func TestMetricsObject(t *testing.T) {
	metricsObjectMock := fake.MetricsObjectMock{
		GetDurationMetricsAttributesFunc: func() []attribute.KeyValue {
			return nil
		},
		GetMetricsAttributesFunc: func() []attribute.KeyValue {
			return nil
		},
		GetEndTimeFunc: func() time.Time {
			return time.Now().UTC()
		},
		GetStartTimeFunc: func() time.Time {
			return time.Now().UTC()
		},
		IsEndTimeSetFunc: func() bool {
			return true
		},
		GetPreviousVersionFunc: func() string {
			return "version"
		},
		GetParentNameFunc: func() string {
			return "parent"
		},
		GetNamespaceFunc: func() string {
			return "namespace"
		},
	}

	wrapper := MetricsObjectWrapper{Obj: &metricsObjectMock}

	_ = wrapper.GetDurationMetricsAttributes()
	require.Len(t, metricsObjectMock.GetDurationMetricsAttributesCalls(), 1)

	_ = wrapper.GetMetricsAttributes()
	require.Len(t, metricsObjectMock.GetMetricsAttributesCalls(), 1)

	_ = wrapper.GetEndTime()
	require.Len(t, metricsObjectMock.GetEndTimeCalls(), 1)

	_ = wrapper.GetStartTime()
	require.Len(t, metricsObjectMock.GetStartTimeCalls(), 1)

	_ = wrapper.IsEndTimeSet()
	require.Len(t, metricsObjectMock.IsEndTimeSetCalls(), 1)

	_ = wrapper.GetPreviousVersion()
	require.Len(t, metricsObjectMock.GetPreviousVersionCalls(), 1)

	_ = wrapper.GetParentName()
	require.Len(t, metricsObjectMock.GetParentNameCalls(), 1)

	_ = wrapper.GetNamespace()
	require.Len(t, metricsObjectMock.GetNamespaceCalls(), 1)
}
