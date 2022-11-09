package api_test

import (
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeptnApp(t *testing.T) {
	app := &v1alpha1.KeptnApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "namespace",
		},
		Spec: v1alpha1.KeptnAppSpec{
			Version: "version",
		},
	}

	appVersionName := app.GetAppVersionName()
	require.Equal(t, "app-version", appVersionName)

	appVersion := app.GenerateAppVersion("prev", map[string]string{})
	require.Equal(t, v1alpha1.KeptnAppVersion{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{},
			Name:        "app-version",
			Namespace:   "namespace",
		},
		Spec: v1alpha1.KeptnAppVersionSpec{
			KeptnAppSpec: v1alpha1.KeptnAppSpec{
				Version: "version",
			},
			AppName:         "app",
			PreviousVersion: "prev",
		},
	}, appVersion)
}
