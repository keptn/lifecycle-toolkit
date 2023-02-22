package interfaces

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/interfaces/fake"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
)

func TestActiveMetricsObjectWrapper(t *testing.T) {
	appVersion := v1alpha3.KeptnAppVersion{
		Status: v1alpha3.KeptnAppVersionStatus{
			Status:       apicommon.StateFailed,
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
