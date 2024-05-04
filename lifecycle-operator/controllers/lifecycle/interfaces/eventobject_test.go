package interfaces

import (
	"testing"

	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/interfaces/fake"
	"github.com/stretchr/testify/require"
)

func TestEventObjectWrapper(t *testing.T) {
	appVersion := apilifecycle.KeptnAppVersion{
		Status: apilifecycle.KeptnAppVersionStatus{
			Status:       apicommon.StateFailed,
			CurrentPhase: "test",
		},
	}

	object, err := NewEventObjectWrapperFromClientObject(&appVersion)
	require.Nil(t, err)

	require.NotEmpty(t, object.GetEventAnnotations())
}

func TestEventObject(t *testing.T) {
	EventObjectMock := fake.EventObjectMock{
		GetEventAnnotationsFunc: func() map[string]string {
			return nil
		},
	}

	wrapper := EventObjectWrapper{Obj: &EventObjectMock}

	_ = wrapper.GetEventAnnotations()
	require.Len(t, EventObjectMock.GetEventAnnotationsCalls(), 1)
}
