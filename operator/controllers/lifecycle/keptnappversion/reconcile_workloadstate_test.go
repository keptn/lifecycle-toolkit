package keptnappversion

import (
	"context"
	"testing"

	lfcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

//nolint:dogsled
func TestKeptnAppVersionReconciler_reconcileWorkloads_noWorkloads(t *testing.T) {
	appVersion := &lfcv1alpha3.KeptnAppVersion{
		ObjectMeta: v1.ObjectMeta{
			Name:      "appversion",
			Namespace: "default",
		},
		Spec: lfcv1alpha3.KeptnAppVersionSpec{
			AppName: "app",
		},
	}
	r, _, _, _ := setupReconciler(appVersion)

	state, err := r.reconcileWorkloads(context.TODO(), appVersion)
	require.Nil(t, err)
	require.Equal(t, apicommon.StateSucceeded, state)

	err = r.Client.Get(context.TODO(), types.NamespacedName{Namespace: appVersion.Namespace, Name: appVersion.Name}, appVersion)
	require.Nil(t, err)
	require.Equal(t, apicommon.StateSucceeded, appVersion.Status.WorkloadOverallStatus)
	require.Len(t, appVersion.Status.WorkloadStatus, 0)
}

//nolint:dogsled
func TestKeptnAppVersionReconciler_reconcileWorkloads(t *testing.T) {
	appVersion := &lfcv1alpha3.KeptnAppVersion{
		ObjectMeta: v1.ObjectMeta{
			Name:      "appversion",
			Namespace: "default",
		},
		Spec: lfcv1alpha3.KeptnAppVersionSpec{
			KeptnAppSpec: lfcv1alpha3.KeptnAppSpec{
				Workloads: []lfcv1alpha3.KeptnWorkloadRef{
					{
						Name:    "workload",
						Version: "ver1",
					},
				},
			},
			AppName: "app",
		},
	}
	r, _, _, _ := setupReconciler(appVersion)

	// No workloadInstances are created yet, should stay in Pending state

	state, err := r.reconcileWorkloads(context.TODO(), appVersion)
	require.Nil(t, err)
	require.Equal(t, apicommon.StatePending, state)

	err = r.Client.Get(context.TODO(), types.NamespacedName{Namespace: appVersion.Namespace, Name: appVersion.Name}, appVersion)
	require.Nil(t, err)
	require.Equal(t, apicommon.StatePending, appVersion.Status.WorkloadOverallStatus)
	require.Len(t, appVersion.Status.WorkloadStatus, 1)
	require.Equal(t, []lfcv1alpha3.WorkloadStatus{
		{
			Workload: lfcv1alpha3.KeptnWorkloadRef{
				Name:    "workload",
				Version: "ver1",
			},
			Status: apicommon.StatePending,
		},
	}, appVersion.Status.WorkloadStatus)

	// Creating WorkloadInstace that is not part of the App -> should stay Pending

	wi1 := &lfcv1alpha3.KeptnWorkloadInstance{
		ObjectMeta: v1.ObjectMeta{
			Name:      "workload",
			Namespace: "default",
		},
		Spec: lfcv1alpha3.KeptnWorkloadInstanceSpec{
			KeptnWorkloadSpec: lfcv1alpha3.KeptnWorkloadSpec{
				AppName: "app2",
			},
		},
	}

	err = r.Client.Create(context.TODO(), wi1)
	require.Nil(t, err)

	state, err = r.reconcileWorkloads(context.TODO(), appVersion)
	require.Nil(t, err)
	require.Equal(t, apicommon.StatePending, state)

	err = r.Client.Get(context.TODO(), types.NamespacedName{Namespace: appVersion.Namespace, Name: appVersion.Name}, appVersion)
	require.Nil(t, err)
	require.Equal(t, apicommon.StatePending, appVersion.Status.WorkloadOverallStatus)
	require.Len(t, appVersion.Status.WorkloadStatus, 1)
	require.Equal(t, []lfcv1alpha3.WorkloadStatus{
		{
			Workload: lfcv1alpha3.KeptnWorkloadRef{
				Name:    "workload",
				Version: "ver1",
			},
			Status: apicommon.StatePending,
		},
	}, appVersion.Status.WorkloadStatus)

	// Creating WorkloadInstance of App with progressing state -> appVersion should be Progressing

	wi2 := &lfcv1alpha3.KeptnWorkloadInstance{
		ObjectMeta: v1.ObjectMeta{
			Name:      "app-workload-ver1",
			Namespace: "default",
		},
		Spec: lfcv1alpha3.KeptnWorkloadInstanceSpec{
			KeptnWorkloadSpec: lfcv1alpha3.KeptnWorkloadSpec{
				AppName: "app",
			},
		},
	}

	err = r.Client.Create(context.TODO(), wi2)
	require.Nil(t, err)

	err = r.Client.Get(context.TODO(), types.NamespacedName{Namespace: wi2.Namespace, Name: wi2.Name}, wi2)
	require.Nil(t, err)

	wi2.Status.Status = apicommon.StateProgressing
	err = r.Client.Update(context.TODO(), wi2)
	require.Nil(t, err)

	state, err = r.reconcileWorkloads(context.TODO(), appVersion)
	require.Nil(t, err)
	require.Equal(t, apicommon.StateProgressing, state)

	err = r.Client.Get(context.TODO(), types.NamespacedName{Namespace: appVersion.Namespace, Name: appVersion.Name}, appVersion)
	require.Nil(t, err)
	require.Equal(t, apicommon.StateProgressing, appVersion.Status.WorkloadOverallStatus)
	require.Len(t, appVersion.Status.WorkloadStatus, 1)
	require.Equal(t, []lfcv1alpha3.WorkloadStatus{
		{
			Workload: lfcv1alpha3.KeptnWorkloadRef{
				Name:    "workload",
				Version: "ver1",
			},
			Status: apicommon.StateProgressing,
		},
	}, appVersion.Status.WorkloadStatus)

	// Updating WorkloadInstance of App with succeeded state -> appVersion should be Succeeded

	err = r.Client.Get(context.TODO(), types.NamespacedName{Namespace: wi2.Namespace, Name: wi2.Name}, wi2)
	require.Nil(t, err)

	wi2.Status.Status = apicommon.StateSucceeded
	err = r.Client.Update(context.TODO(), wi2)
	require.Nil(t, err)

	state, err = r.reconcileWorkloads(context.TODO(), appVersion)
	require.Nil(t, err)
	require.Equal(t, apicommon.StateSucceeded, state)

	err = r.Client.Get(context.TODO(), types.NamespacedName{Namespace: appVersion.Namespace, Name: appVersion.Name}, appVersion)
	require.Nil(t, err)
	require.Equal(t, apicommon.StateSucceeded, appVersion.Status.WorkloadOverallStatus)
	require.Len(t, appVersion.Status.WorkloadStatus, 1)
	require.Equal(t, []lfcv1alpha3.WorkloadStatus{
		{
			Workload: lfcv1alpha3.KeptnWorkloadRef{
				Name:    "workload",
				Version: "ver1",
			},
			Status: apicommon.StateSucceeded,
		},
	}, appVersion.Status.WorkloadStatus)
}
