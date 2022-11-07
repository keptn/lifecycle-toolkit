package common

import (
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPhaseItemWrapper_GetState(t *testing.T) {
	appVersion := &v1alpha1.KeptnAppVersion{
		Status: v1alpha1.KeptnAppVersionStatus{
			Status:       common.StateFailed,
			CurrentPhase: "test",
		},
	}

	object, err := NewPhaseItemWrapperFromClientObject(appVersion)
	require.Nil(t, err)

	require.Equal(t, "test", object.GetCurrentPhase())

	object.Complete()

	require.NotZero(t, appVersion.Status.EndTime)
}
