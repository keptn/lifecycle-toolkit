package common

import (
	"context"
	"fmt"

	lfcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2/common"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/metric"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func AddApp(c client.Client, name string) error {
	app := &lfcv1alpha2.KeptnApp{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:       name,
			Namespace:  "default",
			Generation: 1,
		},
		Spec: lfcv1alpha2.KeptnAppSpec{
			Version:  "1.0.0",
			Revision: 1,
		},
		Status: lfcv1alpha2.KeptnAppStatus{},
	}
	return c.Create(context.TODO(), app)

}

func AddAppVersion(c client.Client, namespace string, appName string, version string, workloads []lfcv1alpha2.KeptnWorkloadRef, status lfcv1alpha2.KeptnAppVersionStatus) error {
	appVersionName := fmt.Sprintf("%s-%s", appName, version)
	app := &lfcv1alpha2.KeptnAppVersion{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:       appVersionName,
			Namespace:  namespace,
			Generation: 1,
		},
		Spec: lfcv1alpha2.KeptnAppVersionSpec{
			KeptnAppSpec: lfcv1alpha2.KeptnAppSpec{
				Version:   version,
				Workloads: workloads,
				Revision:  1,
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

func AddWorkloadInstance(c client.Client, name string, namespace string) error {
	wi := &lfcv1alpha2.KeptnWorkloadInstance{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: lfcv1alpha2.KeptnWorkloadInstanceSpec{
			KeptnWorkloadSpec: lfcv1alpha2.KeptnWorkloadSpec{
				AppName: "some-app",
				Version: "1.0.0",
			},
			WorkloadName:    "some-app-some-workload",
			PreviousVersion: "",
			TraceId:         nil,
		},
		Status: lfcv1alpha2.KeptnWorkloadInstanceStatus{
			DeploymentStatus:                   apicommon.StateSucceeded,
			PreDeploymentStatus:                apicommon.StateSucceeded,
			PostDeploymentStatus:               apicommon.StateSucceeded,
			PreDeploymentEvaluationStatus:      apicommon.StateSucceeded,
			PostDeploymentEvaluationStatus:     apicommon.StateSucceeded,
			CurrentPhase:                       apicommon.PhaseWorkloadPostEvaluation.ShortName,
			PreDeploymentTaskStatus:            nil,
			PostDeploymentTaskStatus:           nil,
			PreDeploymentEvaluationTaskStatus:  nil,
			PostDeploymentEvaluationTaskStatus: nil,
			Status:                             apicommon.StateSucceeded,
			StartTime:                          metav1.Time{},
			EndTime:                            metav1.Time{},
		},
	}
	return c.Create(context.TODO(), wi)
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
