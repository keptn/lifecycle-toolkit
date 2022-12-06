package v1alpha2_test

import (
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2"
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2/common"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestKeptnApp(t *testing.T) {
	app := &v1alpha2.KeptnApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "app",
			Namespace: "namespace",
		},
		Spec: v1alpha2.KeptnAppSpec{
			Version: "version",
		},
	}

	appVersionName := app.GetAppVersionName()
	require.Equal(t, "app-version", appVersionName)

	appVersion := app.GenerateAppVersion("prev", map[string]string{})
	require.Equal(t, v1alpha2.KeptnAppVersion{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{},
			Name:        "app-version",
			Namespace:   "namespace",
		},
		Spec: v1alpha2.KeptnAppVersionSpec{
			KeptnAppSpec: v1alpha2.KeptnAppSpec{
				Version: "version",
			},
			AppName:         "app",
			PreviousVersion: "prev",
		},
	}, appVersion)

	require.Equal(t, []attribute.KeyValue{
		common.AppName.String("app"),
		common.AppVersion.String("version"),
	}, app.GetSpanAttributes())
}
