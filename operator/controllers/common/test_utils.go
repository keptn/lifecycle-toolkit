package common

import (
	"context"
	"fmt"

	lfcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/metric"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func AddApp(c client.Client, name string) error {
	app := &lfcv1alpha3.KeptnApp{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:       name,
			Namespace:  "default",
			Generation: 1,
		},
		Spec: lfcv1alpha3.KeptnAppSpec{
			Version: "1.0.0",
		},
		Status: lfcv1alpha3.KeptnAppStatus{},
	}
	return c.Create(context.TODO(), app)
}

func UpdateAppRevision(c client.Client, name string, revision uint) error {
	app := &lfcv1alpha3.KeptnApp{}
	err := c.Get(context.TODO(), types.NamespacedName{Namespace: "default", Name: name}, app)
	if err != nil {
		return err
	}
	app.Spec.Revision = revision
	app.Generation = int64(revision)
	return c.Update(context.TODO(), app)
}

func AddAppVersion(c client.Client, namespace string, appName string, version string, workloads []lfcv1alpha3.KeptnWorkloadRef, status lfcv1alpha3.KeptnAppVersionStatus) error {
	appVersionName := fmt.Sprintf("%s-%s", appName, version)
	app := &lfcv1alpha3.KeptnAppVersion{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:       appVersionName,
			Namespace:  namespace,
			Generation: 1,
		},
		Spec: lfcv1alpha3.KeptnAppVersionSpec{
			KeptnAppSpec: lfcv1alpha3.KeptnAppSpec{
				Version:   version,
				Workloads: workloads,
			},
			AppName: appName,
			TraceId: map[string]string{
				"traceparent": "parent-trace",
			},
		},
		Status: status,
	}
	return c.Create(context.TODO(), app)
}

func InitAppMeters() apicommon.KeptnMeters {
	provider := metric.NewMeterProvider()
	meter := provider.Meter("keptn/task")
	appCount, _ := meter.SyncInt64().Counter("keptn.app.count", instrument.WithDescription("a simple counter for Keptn Apps"))
	appDuration, _ := meter.SyncFloat64().Histogram("keptn.app.duration", instrument.WithDescription("a histogram of duration for Keptn Apps"), instrument.WithUnit("s"))
	deploymentCount, _ := meter.SyncInt64().Counter("keptn.deployment.count", instrument.WithDescription("a simple counter for Keptn Deployments"))
	deploymentDuration, _ := meter.SyncFloat64().Histogram("keptn.deployment.duration", instrument.WithDescription("a histogram of duration for Keptn Deployments"), instrument.WithUnit(unit.Unit("s")))

	meters := apicommon.KeptnMeters{
		AppCount:           appCount,
		AppDuration:        appDuration,
		DeploymentCount:    deploymentCount,
		DeploymentDuration: deploymentDuration,
	}
	return meters
}
