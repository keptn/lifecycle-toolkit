package common

import (
	"context"
	"fmt"
	lfcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	keptncommon "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1/common"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func AddApp(c client.Client, name string) error {
	app := &lfcv1alpha1.KeptnApp{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: "default",
		},
		Spec: lfcv1alpha1.KeptnAppSpec{
			Version: "1.0.0",
		},
		Status: lfcv1alpha1.KeptnAppStatus{},
	}
	return c.Create(context.TODO(), app)

}

func AddAppVersion(c client.Client, namespace string, appName string, version string, workloads []lfcv1alpha1.KeptnWorkloadRef, status lfcv1alpha1.KeptnAppVersionStatus) error {
	appVersionName := fmt.Sprintf("%s-%s", appName, version)
	app := &lfcv1alpha1.KeptnAppVersion{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      appVersionName,
			Namespace: namespace,
		},
		Spec: lfcv1alpha1.KeptnAppVersionSpec{
			KeptnAppSpec: lfcv1alpha1.KeptnAppSpec{
				Version:   version,
				Workloads: workloads,
			},
			AppName: appName,
		},
		Status: status,
	}
	return c.Create(context.TODO(), app)
}

func AddWorkloadInstance(c client.Client, name string, namespace string) error {
	wi := &lfcv1alpha1.KeptnWorkloadInstance{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: lfcv1alpha1.KeptnWorkloadInstanceSpec{
			KeptnWorkloadSpec: lfcv1alpha1.KeptnWorkloadSpec{
				AppName: "some-app",
				Version: "1.0.0",
			},
			WorkloadName:    "some-app-some-workload",
			PreviousVersion: "",
			TraceId:         nil,
		},
		Status: lfcv1alpha1.KeptnWorkloadInstanceStatus{
			DeploymentStatus:                   keptncommon.StateSucceeded,
			PreDeploymentStatus:                keptncommon.StateSucceeded,
			PostDeploymentStatus:               keptncommon.StateSucceeded,
			PreDeploymentEvaluationStatus:      keptncommon.StateSucceeded,
			PostDeploymentEvaluationStatus:     keptncommon.StateSucceeded,
			CurrentPhase:                       keptncommon.PhaseWorkloadPostEvaluation.ShortName,
			PreDeploymentTaskStatus:            nil,
			PostDeploymentTaskStatus:           nil,
			PreDeploymentEvaluationTaskStatus:  nil,
			PostDeploymentEvaluationTaskStatus: nil,
			Status:                             keptncommon.StateSucceeded,
			StartTime:                          metav1.Time{},
			EndTime:                            metav1.Time{},
		},
	}
	return c.Create(context.TODO(), wi)

}
