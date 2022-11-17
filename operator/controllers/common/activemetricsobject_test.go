package common

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
)

func TestActiveMetricsObjectWrapper(t *testing.T) {
	appVersion := v1alpha1.KeptnAppVersion{
		Status: v1alpha1.KeptnAppVersionStatus{
			Status:       common.StateFailed,
			CurrentPhase: "test",
		},
	}

	object, err := NewActiveMetricsObjectWrapperFromClientObject(&appVersion)
	require.Nil(t, err)

	require.False(t, object.IsEndTimeSet())
}

func TestActiveMetricsObject(t *testing.T) {
	activeMetricsObjectMock := fake.ActiveMetricsObjectMock{
		GetActiveMetricsAttributesFunc: func() []attribute.KeyValue {
			return nil
		},
		IsEndTimeSetFunc: func() bool {
			return true
		},
	}

	wrapper := ActiveMetricsObjectWrapper{Obj: &activeMetricsObjectMock}

	_ = wrapper.GetActiveMetricsAttributes()
	require.Len(t, activeMetricsObjectMock.GetActiveMetricsAttributesCalls(), 1)

	_ = wrapper.IsEndTimeSet()
	require.Len(t, activeMetricsObjectMock.IsEndTimeSetCalls(), 1)
}
