package keptnworkloadversion

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/config"
	keptncontext "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/context"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/evaluation"
	evaluationfake "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/evaluation/fake"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/eventsender"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/phase"
	phasefake "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/phase/fake"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/telemetry"
	telemetryfake "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/telemetry/fake"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/testcommon"
	controllererrors "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/errors"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/interfaces"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8sfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestKeptnWorkloadVersionReconciler_reconcileDeployment_FailedReplicaSet(t *testing.T) {

	rep := int32(1)
	replicasetFail := makeReplicaSet("myrep", "default", &rep, 0)
	workloadVersion := makeWorkloadVersionWithRef(replicasetFail.ObjectMeta, "ReplicaSet")

	fakeClient := testcommon.NewTestClient(replicasetFail, workloadVersion)

	r := &KeptnWorkloadVersionReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadVersion)
	require.Nil(t, err)
	require.Equal(t, apicommon.StateProgressing, keptnState)
	require.False(t, workloadVersion.Status.DeploymentStartTime.IsZero())
}

func TestKeptnWorkloadVersionReconciler_reconcileDeployment_UnavailableReplicaSet(t *testing.T) {

	rep := int32(1)
	replicasetFail := makeReplicaSet("myrep", "default", &rep, 0)
	workloadVersion := makeWorkloadVersionWithRef(replicasetFail.ObjectMeta, "ReplicaSet")

	// do not put the ReplicaSet into the cluster
	fakeClient := testcommon.NewTestClient(workloadVersion)

	r := &KeptnWorkloadVersionReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadVersion)
	require.NotNil(t, err)
	require.Equal(t, apicommon.StateUnknown, keptnState)
	require.True(t, workloadVersion.Status.DeploymentStartTime.IsZero())
}

func TestKeptnWorkloadVersionReconciler_reconcileDeployment_WorkloadDeploymentTimedOut(t *testing.T) {

	rep := int32(1)
	replicaset := makeReplicaSet("myrep", "default", &rep, 0)
	workloadVersion := makeWorkloadVersionWithRef(replicaset.ObjectMeta, "ReplicaSet")

	fakeClient := testcommon.NewTestClient(replicaset, workloadVersion)

	fakeRecorder := record.NewFakeRecorder(100)

	r := &KeptnWorkloadVersionReconciler{
		Client:      fakeClient,
		Config:      config.Instance(),
		EventSender: eventsender.NewK8sSender(fakeRecorder),
	}

	r.Config.SetObservabilityTimeout(metav1.Duration{
		Duration: 5 * time.Second,
	})

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadVersion)
	require.Nil(t, err)
	require.Equal(t, apicommon.StateProgressing, keptnState)
	require.False(t, workloadVersion.Status.DeploymentStartTime.IsZero())

	//revert the start time parameter backwards to check the timer
	workloadVersion.Status.DeploymentStartTime = metav1.Time{
		Time: workloadVersion.Status.DeploymentStartTime.Add(-10 * time.Second),
	}

	err = r.Client.Status().Update(context.TODO(), workloadVersion)
	require.Nil(t, err)

	keptnState, err = r.reconcileDeployment(context.TODO(), workloadVersion)
	require.Nil(t, err)
	require.Equal(t, apicommon.StateFailed, keptnState)
	require.False(t, workloadVersion.Status.DeploymentStartTime.IsZero())

	event := <-fakeRecorder.Events
	require.Equal(t, strings.Contains(event, workloadVersion.GetName()), true, "wrong workloadVersion")
	require.Equal(t, strings.Contains(event, workloadVersion.GetNamespace()), true, "wrong namespace")
	require.Equal(t, strings.Contains(event, "has reached timeout"), true, "wrong message")
}

func TestKeptnWorkloadVersionReconciler_reconcileDeployment_FailedStatefulSet(t *testing.T) {

	rep := int32(1)
	statefulsetFail := makeStatefulSet("mystat", "default", &rep, 0)
	workloadVersion := makeWorkloadVersionWithRef(statefulsetFail.ObjectMeta, "StatefulSet")

	fakeClient := testcommon.NewTestClient(statefulsetFail, workloadVersion)
	r := &KeptnWorkloadVersionReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadVersion)
	require.Nil(t, err)
	require.Equal(t, apicommon.StateProgressing, keptnState)
	require.False(t, workloadVersion.Status.DeploymentStartTime.IsZero())
}

func TestKeptnWorkloadVersionReconciler_reconcileDeployment_UnavailableStatefulSet(t *testing.T) {

	rep := int32(1)
	statefulSetFail := makeStatefulSet("mystat", "default", &rep, 0)
	workloadVersion := makeWorkloadVersionWithRef(statefulSetFail.ObjectMeta, "StatefulSet")

	// do not put the StatefulSet into the cluster
	fakeClient := testcommon.NewTestClient(workloadVersion)

	r := &KeptnWorkloadVersionReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadVersion)
	require.NotNil(t, err)
	require.Equal(t, apicommon.StateUnknown, keptnState)
	require.True(t, workloadVersion.Status.DeploymentStartTime.IsZero())
}

func TestKeptnWorkloadVersionReconciler_reconcileDeployment_FailedDaemonSet(t *testing.T) {

	daemonSetFail := makeDaemonSet("mystat", "default", 1, 0)
	workloadVersion := makeWorkloadVersionWithRef(daemonSetFail.ObjectMeta, "DaemonSet")

	fakeClient := testcommon.NewTestClient(daemonSetFail, workloadVersion)

	r := &KeptnWorkloadVersionReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadVersion)
	require.Nil(t, err)
	require.Equal(t, apicommon.StateProgressing, keptnState)
	require.False(t, workloadVersion.Status.DeploymentStartTime.IsZero())
}

func TestKeptnWorkloadVersionReconciler_reconcileDeployment_UnavailableDaemonSet(t *testing.T) {
	daemonSetFail := makeDaemonSet("mystat", "default", 1, 0)
	workloadVersion := makeWorkloadVersionWithRef(daemonSetFail.ObjectMeta, "DaemonSet")

	// do not put the DaemonSet into the cluster
	fakeClient := testcommon.NewTestClient(workloadVersion)

	r := &KeptnWorkloadVersionReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadVersion)
	require.NotNil(t, err)
	require.Equal(t, apicommon.StateUnknown, keptnState)
	require.True(t, workloadVersion.Status.DeploymentStartTime.IsZero())
}

func TestKeptnWorkloadVersionReconciler_reconcileDeployment_ReadyReplicaSet(t *testing.T) {

	rep := int32(1)
	replicaSet := makeReplicaSet("myrep", "default", &rep, 1)
	workloadVersion := makeWorkloadVersionWithRef(replicaSet.ObjectMeta, "ReplicaSet")

	fakeClient := testcommon.NewTestClient(replicaSet, workloadVersion)

	r := &KeptnWorkloadVersionReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadVersion)
	require.Nil(t, err)
	require.Equal(t, apicommon.StateSucceeded, keptnState)
	require.False(t, workloadVersion.Status.DeploymentStartTime.IsZero())
}

func TestKeptnWorkloadVersionReconciler_reconcileDeployment_ReadyStatefulSet(t *testing.T) {

	rep := int32(1)
	statefulSet := makeStatefulSet("mystat", "default", &rep, 1)
	workloadVersion := makeWorkloadVersionWithRef(statefulSet.ObjectMeta, "StatefulSet")

	fakeClient := testcommon.NewTestClient(statefulSet, workloadVersion)

	r := &KeptnWorkloadVersionReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadVersion)
	require.Nil(t, err)
	require.Equal(t, apicommon.StateSucceeded, keptnState)
	require.False(t, workloadVersion.Status.DeploymentStartTime.IsZero())
}

func TestKeptnWorkloadVersionReconciler_reconcileDeployment_ReadyDaemonSet(t *testing.T) {

	daemonSet := makeDaemonSet("mystat", "default", 1, 1)
	workloadVersion := makeWorkloadVersionWithRef(daemonSet.ObjectMeta, "DaemonSet")

	fakeClient := testcommon.NewTestClient(daemonSet, workloadVersion)

	r := &KeptnWorkloadVersionReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadVersion)
	require.Nil(t, err)
	require.Equal(t, apicommon.StateSucceeded, keptnState)
	require.False(t, workloadVersion.Status.DeploymentStartTime.IsZero())
}

func TestKeptnWorkloadVersionReconciler_reconcileDeployment_UnsupportedReferenceKind(t *testing.T) {

	workloadVersion := makeWorkloadVersionWithRef(metav1.ObjectMeta{}, "Unknown")
	fakeClient := testcommon.NewTestClient(workloadVersion)
	r := &KeptnWorkloadVersionReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadVersion)
	require.ErrorIs(t, err, controllererrors.ErrUnsupportedWorkloadVersionResourceReference)
	require.Equal(t, apicommon.StateUnknown, keptnState)
	require.True(t, workloadVersion.Status.DeploymentStartTime.IsZero())
}

func makeReplicaSet(name string, namespace string, wanted *int32, available int32) *appsv1.ReplicaSet {

	return &appsv1.ReplicaSet{
		TypeMeta: metav1.TypeMeta{
			Kind: "ReplicaSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			UID:       types.UID(name),
		},
		Spec: appsv1.ReplicaSetSpec{
			Replicas: wanted,
		},
		Status: appsv1.ReplicaSetStatus{
			AvailableReplicas: available,
		},
	}

}

func makeStatefulSet(name string, namespace string, wanted *int32, available int32) *appsv1.StatefulSet {

	return &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind: "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			UID:       types.UID(name),
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: wanted,
		},
		Status: appsv1.StatefulSetStatus{
			AvailableReplicas: available,
		},
	}

}

func makeDaemonSet(name string, namespace string, wanted int32, available int32) *appsv1.DaemonSet {

	return &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind: "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			UID:       types.UID(name),
		},
		Spec: appsv1.DaemonSetSpec{},
		Status: appsv1.DaemonSetStatus{
			DesiredNumberScheduled: wanted,
			NumberReady:            available,
		},
	}

}

func Test_getAppVersionForWorkloadVersion(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name           string
		wli            *apilifecycle.KeptnWorkloadVersion
		list           *apilifecycle.KeptnAppVersionList
		wantFound      bool
		wantAppVersion apilifecycle.KeptnAppVersion
		wantErr        bool
	}{
		{
			name: "no appVersions",
			wli: &apilifecycle.KeptnWorkloadVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-workloadVersion",
					Namespace: "default",
				},
				Spec: apilifecycle.KeptnWorkloadVersionSpec{
					KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
						AppName: "my-app",
						Version: "1.0",
					},
					WorkloadName: "my-app-my-workload",
				},
			},
			list:           &apilifecycle.KeptnAppVersionList{},
			wantFound:      false,
			wantAppVersion: apilifecycle.KeptnAppVersion{},
			wantErr:        false,
		},
		{
			name: "appVersion found",
			wli: &apilifecycle.KeptnWorkloadVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-workloadVersion",
					Namespace: "default",
				},
				Spec: apilifecycle.KeptnWorkloadVersionSpec{
					KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
						AppName: "my-app",
						Version: "1.0",
					},
					WorkloadName: "my-app-my-workload",
				},
			},
			list: &apilifecycle.KeptnAppVersionList{
				Items: []apilifecycle.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "my-app",
							Namespace:         "default",
							CreationTimestamp: metav1.Time{Time: now},
						},
						Spec: apilifecycle.KeptnAppVersionSpec{
							KeptnAppSpec: apilifecycle.KeptnAppSpec{
								Version: "1.0",
								Workloads: []apilifecycle.KeptnWorkloadRef{
									{
										Name:    "my-workload",
										Version: "1.0",
									},
								},
							},
							AppName: "my-app",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "my-app2",
							Namespace:         "default",
							CreationTimestamp: metav1.Time{Time: now.Add(5 * time.Second)},
						},
						Spec: apilifecycle.KeptnAppVersionSpec{
							KeptnAppSpec: apilifecycle.KeptnAppSpec{
								Version: "2.0",
								Workloads: []apilifecycle.KeptnWorkloadRef{
									{
										Name:    "my-workload",
										Version: "1.0",
									},
								},
							},
							AppName: "my-app",
						},
					},
				},
			},
			wantFound: true,
			wantAppVersion: apilifecycle.KeptnAppVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "my-app2",
					Namespace:       "default",
					ResourceVersion: "999",
				},
				Spec: apilifecycle.KeptnAppVersionSpec{
					KeptnAppSpec: apilifecycle.KeptnAppSpec{
						Version: "2.0",
						Workloads: []apilifecycle.KeptnWorkloadRef{
							{
								Name:    "my-workload",
								Version: "1.0",
							},
						},
					},
					AppName: "my-app",
				},
			},
			wantErr: false,
		},
		{
			name: "appVersion deprecated",
			wli: &apilifecycle.KeptnWorkloadVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-workloadVersion",
					Namespace: "default",
				},
				Spec: apilifecycle.KeptnWorkloadVersionSpec{
					KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
						AppName: "my-app",
						Version: "1.0",
					},
					WorkloadName: "my-app-my-workload",
				},
			},
			list: &apilifecycle.KeptnAppVersionList{
				Items: []apilifecycle.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "my-app",
							Namespace: "default",
						},
						Spec: apilifecycle.KeptnAppVersionSpec{
							KeptnAppSpec: apilifecycle.KeptnAppSpec{
								Version: "1.0",
								Workloads: []apilifecycle.KeptnWorkloadRef{
									{
										Name:    "my-workload",
										Version: "1.0",
									},
								},
							},
							AppName: "my-app",
						},
						Status: apilifecycle.KeptnAppVersionStatus{
							Status: apicommon.StateDeprecated,
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "my-app2",
							Namespace: "default",
						},
						Spec: apilifecycle.KeptnAppVersionSpec{
							KeptnAppSpec: apilifecycle.KeptnAppSpec{
								Version: "2.0",
								Workloads: []apilifecycle.KeptnWorkloadRef{
									{
										Name:    "my-workload",
										Version: "1.0",
									},
								},
							},
							AppName: "my-app",
						},
						Status: apilifecycle.KeptnAppVersionStatus{
							Status: apicommon.StateDeprecated,
						},
					},
				},
			},
			wantFound:      false,
			wantAppVersion: apilifecycle.KeptnAppVersion{},
			wantErr:        false,
		},
		{
			name: "no workload for appversion",
			wli: &apilifecycle.KeptnWorkloadVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-workloadVersion3",
					Namespace: "default",
				},
				Spec: apilifecycle.KeptnWorkloadVersionSpec{
					KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
						AppName: "my-app333",
						Version: "1.0.0",
					},
					WorkloadName: "my-app-my-workload",
				},
			},
			list: &apilifecycle.KeptnAppVersionList{
				Items: []apilifecycle.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "my-app",
							Namespace: "default",
						},
						Spec: apilifecycle.KeptnAppVersionSpec{
							KeptnAppSpec: apilifecycle.KeptnAppSpec{
								Version: "1.0",
								Workloads: []apilifecycle.KeptnWorkloadRef{
									{
										Name:    "my-workload",
										Version: "1.0",
									},
								},
							},
							AppName: "my-app",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "my-app2",
							Namespace: "default",
						},
						Spec: apilifecycle.KeptnAppVersionSpec{
							KeptnAppSpec: apilifecycle.KeptnAppSpec{
								Version: "2.0",
								Workloads: []apilifecycle.KeptnWorkloadRef{
									{
										Name:    "my-workload",
										Version: "1.0",
									},
								},
							},
							AppName: "my-app",
						},
					},
				},
			},
			wantFound:      false,
			wantAppVersion: apilifecycle.KeptnAppVersion{},
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apilifecycle.AddToScheme(scheme.Scheme)
			require.Nil(t, err)
			r := &KeptnWorkloadVersionReconciler{
				Client: k8sfake.NewClientBuilder().WithLists(tt.list).Build(),
			}
			found, gotAppVersion, err := r.getAppVersionForWorkloadVersion(context.TODO(), tt.wli)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.wantFound, found)
			if tt.wantFound {
				// set the creation timestamp of the returned appVersion to the time zero value because this is
				// set internally by the fake client
				gotAppVersion.ObjectMeta.CreationTimestamp = metav1.Time{Time: time.Time{}}
			}
			require.Equal(t, tt.wantAppVersion, gotAppVersion)
		})
	}
}

func Test_getLatestAppVersion(t *testing.T) {

	now := time.Now()
	type args struct {
		apps *apilifecycle.KeptnAppVersionList
		wli  *apilifecycle.KeptnWorkloadVersion
	}
	tests := []struct {
		name           string
		args           args
		wantFound      bool
		wantAppVersion apilifecycle.KeptnAppVersion
		wantErr        bool
	}{
		{
			name: "app version found",
			args: args{
				apps: &apilifecycle.KeptnAppVersionList{
					Items: []apilifecycle.KeptnAppVersion{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:              "my-app",
								Namespace:         "default",
								CreationTimestamp: metav1.Time{Time: now},
							},
							Spec: apilifecycle.KeptnAppVersionSpec{
								KeptnAppSpec: apilifecycle.KeptnAppSpec{
									Version: "1.0",
									Workloads: []apilifecycle.KeptnWorkloadRef{
										{
											Name:    "my-workload",
											Version: "1.0",
										},
									},
								},
								AppName: "my-app",
							},
						},
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:              "my-app",
								Namespace:         "default",
								CreationTimestamp: metav1.Time{Time: now.Add(5 * time.Second)},
							},
							Spec: apilifecycle.KeptnAppVersionSpec{
								KeptnAppSpec: apilifecycle.KeptnAppSpec{
									Version: "2.0",
									Workloads: []apilifecycle.KeptnWorkloadRef{
										{
											Name:    "my-workload",
											Version: "1.0",
										},
									},
								},
								AppName: "my-app",
							},
						},
					},
				},
				wli: &apilifecycle.KeptnWorkloadVersion{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-workloadVersion",
						Namespace: "default",
					},
					Spec: apilifecycle.KeptnWorkloadVersionSpec{
						KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
							AppName: "my-app",
							Version: "1.0",
						},
						WorkloadName: "my-app-my-workload",
					},
				},
			},
			wantFound: true,
			wantAppVersion: apilifecycle.KeptnAppVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "my-app",
					Namespace:         "default",
					CreationTimestamp: metav1.Time{Time: now.Add(5 * time.Second)},
				},
				Spec: apilifecycle.KeptnAppVersionSpec{
					KeptnAppSpec: apilifecycle.KeptnAppSpec{
						Version: "2.0",
						Workloads: []apilifecycle.KeptnWorkloadRef{
							{
								Name:    "my-workload",
								Version: "1.0",
							},
						},
					},
					AppName: "my-app",
				},
			},
			wantErr: false,
		},
		{
			name: "app version not found",
			args: args{
				apps: &apilifecycle.KeptnAppVersionList{
					Items: []apilifecycle.KeptnAppVersion{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "my-app",
								Namespace: "default",
							},
							Spec: apilifecycle.KeptnAppVersionSpec{
								KeptnAppSpec: apilifecycle.KeptnAppSpec{
									Version: "1.0",
									Workloads: []apilifecycle.KeptnWorkloadRef{
										{
											Name:    "my-other-workload",
											Version: "1.0",
										},
									},
								},
								AppName: "my-app",
							},
						},
					},
				},
				wli: &apilifecycle.KeptnWorkloadVersion{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-workloadVersion",
						Namespace: "default",
					},
					Spec: apilifecycle.KeptnWorkloadVersionSpec{
						KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
							AppName: "my-app",
							Version: "1.0",
						},
						WorkloadName: "my-app-my-workload",
					},
				},
			},
			wantFound:      false,
			wantAppVersion: apilifecycle.KeptnAppVersion{},
			wantErr:        false,
		},
		{
			name: "app version list empty",
			args: args{
				apps: &apilifecycle.KeptnAppVersionList{
					Items: []apilifecycle.KeptnAppVersion{},
				},
				wli: &apilifecycle.KeptnWorkloadVersion{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-workloadVersion",
						Namespace: "default",
					},
					Spec: apilifecycle.KeptnWorkloadVersionSpec{
						KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
							AppName: "my-app",
							Version: "1.0",
						},
						WorkloadName: "my-app-my-workload",
					},
				},
			},
			wantFound:      false,
			wantAppVersion: apilifecycle.KeptnAppVersion{},
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, gotAppVersion, err := getLatestAppVersion(tt.args.apps, tt.args.wli)
			require.Equal(t, tt.wantErr, err != nil)
			require.Equal(t, tt.wantFound, found)
			require.Equal(t, tt.wantAppVersion, gotAppVersion)
		})
	}
}

func TestKeptnWorkloadVersionReconciler_ReconcileNoActionRequired(t *testing.T) {
	r, _, _ := setupReconciler()

	result, err := r.Reconcile(context.TODO(), ctrl.Request{})

	require.Nil(t, err)
	require.NotNil(t, result)
}

func TestKeptnWorkloadVersionReconciler_ReconcileReachCompletion(t *testing.T) {

	testNamespace := "some-ns"

	wi := &apilifecycle.KeptnWorkloadVersion{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "some-wi",
			Namespace: testNamespace,
		},
		Spec: apilifecycle.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
				AppName: "some-app",
				Version: "1.0.0",
				Metadata: map[string]string{
					"foo": "bar",
				},
			},
			WorkloadName:    "some-app-some-workload",
			PreviousVersion: "",
			TraceId:         nil,
		},
		Status: apilifecycle.KeptnWorkloadVersionStatus{
			DeploymentStatus:               apicommon.StateSucceeded,
			PreDeploymentStatus:            apicommon.StateSucceeded,
			PostDeploymentStatus:           apicommon.StateSucceeded,
			PreDeploymentEvaluationStatus:  apicommon.StateSucceeded,
			PostDeploymentEvaluationStatus: apicommon.StateSucceeded,
			CurrentPhase:                   apicommon.PhaseWorkloadPostEvaluation.ShortName,
			Status:                         apicommon.StateSucceeded,
			AppContextMetadata: map[string]string{
				"testy": "test",
			},
			StartTime: metav1.Time{},
			EndTime:   metav1.Time{},
		},
	}

	app := testcommon.ReturnAppVersion(
		testNamespace,
		"some-app",
		"1.0.0",
		[]apilifecycle.KeptnWorkloadRef{
			{
				Name:    "some-workload",
				Version: "1.0.0",
			},
		},
		apilifecycle.KeptnAppVersionStatus{
			PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
		},
	)
	r, eventChannel, _ := setupReconciler(wi, app)

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Namespace: testNamespace,
			Name:      "some-wi",
		},
	}

	result, err := r.Reconcile(context.TODO(), req)

	require.Nil(t, err)

	// do not requeue since we reached completion
	require.False(t, result.Requeue)

	// here we do not expect an event about the application preEvaluation being finished since that  will have been sent in
	// one of the previous reconciliation loops that lead to the first phase being reached
	expectedEvents := []string{
		"CompletedFinished",
	}

	for _, e := range expectedEvents {
		event := <-eventChannel
		assert.Equal(t, strings.Contains(event, req.Name), true, "wrong workloadVersion")
		assert.Equal(t, strings.Contains(event, req.Namespace), true, "wrong namespace")
		assert.Equal(t, strings.Contains(event, e), true, fmt.Sprintf("no %s found in %s", e, event))
	}

	spanHandlerMock := r.SpanHandler.(*telemetryfake.ISpanHandlerMock)

	require.Len(t, spanHandlerMock.GetSpanCalls(), 1)
	require.Len(t, spanHandlerMock.UnbindSpanCalls(), 1)

	// verify the propagation of the context attributes to the span handler
	metadata, b := keptncontext.GetAppMetadataFromContext(spanHandlerMock.GetSpanCalls()[0].Ctx)

	require.True(t, b)
	require.Equal(t, "bar", metadata["foo"])
	require.Equal(t, "test", metadata["testy"])
}

func TestKeptnWorkloadVersionReconciler_ReconcileFailed(t *testing.T) {

	testNamespace := "some-ns"

	wi := &apilifecycle.KeptnWorkloadVersion{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "some-wi",
			Namespace: testNamespace,
		},
		Spec: apilifecycle.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
				AppName: "some-app",
				Version: "1.0.0",
				PreDeploymentTasks: []string{
					"task",
				},
			},
			WorkloadName:    "some-app-some-workload",
			PreviousVersion: "",
			TraceId:         nil,
		},
		Status: apilifecycle.KeptnWorkloadVersionStatus{
			DeploymentStatus:               apicommon.StatePending,
			PreDeploymentStatus:            apicommon.StateProgressing,
			PostDeploymentStatus:           apicommon.StatePending,
			PreDeploymentEvaluationStatus:  apicommon.StatePending,
			PostDeploymentEvaluationStatus: apicommon.StatePending,
			CurrentPhase:                   apicommon.PhaseWorkloadPreDeployment.ShortName,
			Status:                         apicommon.StateProgressing,
			PreDeploymentTaskStatus: []apilifecycle.ItemStatus{
				{
					Name:           "pre-task",
					DefinitionName: "task",
					Status:         apicommon.StateFailed,
				},
			},
		},
	}

	app := testcommon.ReturnAppVersion(

		testNamespace,
		"some-app",
		"1.0.0",
		[]apilifecycle.KeptnWorkloadRef{
			{
				Name:    "some-workload",
				Version: "1.0.0",
			},
		},
		apilifecycle.KeptnAppVersionStatus{
			PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
		},
	)

	r, eventChannel, _ := setupReconciler(app, wi)

	r.PhaseHandler = &phasefake.MockHandler{HandlePhaseFunc: func(ctx context.Context, ctxTrace context.Context, tracer telemetry.ITracer, reconcileObject client.Object, phaseMoqParam apicommon.KeptnPhaseType, reconcilePhase func(phaseCtx context.Context) (apicommon.KeptnState, error)) (phase.PhaseResult, error) {
		piWrapper, _ := interfaces.NewPhaseItemWrapperFromClientObject(reconcileObject)
		piWrapper.SetState(apicommon.StateFailed)
		return phase.PhaseResult{Continue: false, Result: ctrl.Result{Requeue: false}}, nil
	}}

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Namespace: testNamespace,
			Name:      "some-wi",
		},
	}

	result, err := r.Reconcile(context.TODO(), req)

	require.Nil(t, err)

	// here we do not expect an event about the application preEvaluation being finished since that  will have been sent in
	// one of the previous reconciliation loops that lead to the first phase being reached
	expectedEvents := []string{
		"CompletedFailed",
	}

	for _, e := range expectedEvents {
		event := <-eventChannel
		assert.Equal(t, strings.Contains(event, req.Name), true, "wrong workloadVersion")
		assert.Equal(t, strings.Contains(event, req.Namespace), true, "wrong namespace")
		assert.Equal(t, strings.Contains(event, e), true, fmt.Sprintf("no %s found in %s", e, event))
	}

	spanHandlerMock := r.SpanHandler.(*telemetryfake.ISpanHandlerMock)

	require.Len(t, spanHandlerMock.GetSpanCalls(), 1)
	require.Len(t, spanHandlerMock.UnbindSpanCalls(), 1)

	// do not requeue since we reached completion
	require.False(t, result.Requeue)
}

func TestKeptnWorkloadVersionReconciler_ReconcileCouldNotRetrieveObject(t *testing.T) {

	testNamespace := "some-ns"

	r, _, _ := setupReconciler()

	fakeClient := k8sfake.NewClientBuilder().WithScheme(scheme.Scheme).WithInterceptorFuncs(interceptor.Funcs{
		Get: func(ctx context.Context, client client.WithWatch, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
			return errors.New("unexpected error")
		},
	}).Build()

	r.Client = fakeClient

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Namespace: testNamespace,
			Name:      "some-wi",
		},
	}

	result, err := r.Reconcile(context.TODO(), req)

	require.NotNil(t, err)
	require.False(t, result.Requeue)
}

func TestKeptnWorkloadVersionReconciler_ReconcilePreDeploymentEvaluationUnexpectedError(t *testing.T) {

	testNamespace := "some-ns"

	wi := &apilifecycle.KeptnWorkloadVersion{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "some-wi",
			Namespace: testNamespace,
		},
		Spec: apilifecycle.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
				AppName:                  "some-app",
				Version:                  "1.0.0",
				PreDeploymentEvaluations: []string{"my-pre-evaluation"},
			},
			WorkloadName:    "some-app-some-workload",
			PreviousVersion: "",
			TraceId:         nil,
		},
		Status: apilifecycle.KeptnWorkloadVersionStatus{
			CurrentPhase: apicommon.PhaseWorkloadPreDeployment.ShortName,
			StartTime:    metav1.Time{},
			EndTime:      metav1.Time{},
		},
	}

	app := testcommon.ReturnAppVersion(
		testNamespace,
		"some-app",
		"1.0.0",
		[]apilifecycle.KeptnWorkloadRef{
			{
				Name:    "some-workload",
				Version: "1.0.0",
			},
		},
		apilifecycle.KeptnAppVersionStatus{
			PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
		},
	)

	r, _, _ := setupReconciler(wi, app)

	mockEvaluationHandler := r.EvaluationHandler.(*evaluationfake.MockEvaluationHandler)

	mockEvaluationHandler.ReconcileEvaluationsFunc = func(ctx context.Context, phaseCtx context.Context, reconcileObject client.Object, evaluationCreateAttributes evaluation.CreateEvaluationAttributes) ([]apilifecycle.ItemStatus, apicommon.StatusSummary, error) {
		return nil, apicommon.StatusSummary{}, errors.New("unexpected error")
	}

	mockPhaseHandler := &phasefake.MockHandler{
		HandlePhaseFunc: func(ctx context.Context, ctxTrace context.Context, tracer telemetry.ITracer, reconcileObject client.Object, phaseMoqParam apicommon.KeptnPhaseType, reconcilePhase func(phaseCtx context.Context) (apicommon.KeptnState, error)) (phase.PhaseResult, error) {
			return phase.PhaseResult{Continue: false, Result: ctrl.Result{Requeue: true}}, errors.New("unexpected error")
		},
	}

	r.PhaseHandler = mockPhaseHandler

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Namespace: testNamespace,
			Name:      "some-wi",
		},
	}

	result, err := r.Reconcile(context.TODO(), req)

	require.NotNil(t, err)
	require.True(t, result.Requeue)
}

func setupReconciler(objs ...client.Object) (*KeptnWorkloadVersionReconciler, chan string, *telemetryfake.ITracerMock) {
	// setup logger
	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	// fake a tracer
	tr := &telemetryfake.ITracerMock{StartFunc: func(ctx context.Context, spanName string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
		return ctx, trace.SpanFromContext(ctx)
	}}

	tf := &telemetryfake.TracerFactoryMock{GetTracerFunc: func(name string) telemetry.ITracer {
		return tr
	}}

	fakeClient := testcommon.NewTestClient(objs...)

	recorder := record.NewFakeRecorder(100)

	spanHandlerMock := &telemetryfake.ISpanHandlerMock{
		GetSpanFunc: func(ctx context.Context, tracer telemetry.ITracer, reconcileObject client.Object, phase string, links ...trace.Link) (context.Context, trace.Span, error) {
			ctx, span := tracer.Start(ctx, phase, trace.WithSpanKind(trace.SpanKindConsumer))
			return ctx, span, nil
		},
		UnbindSpanFunc: func(_ client.Object, _ string) error {
			return nil
		},
	}

	r := &KeptnWorkloadVersionReconciler{
		Client:        fakeClient,
		Scheme:        scheme.Scheme,
		EventSender:   eventsender.NewK8sSender(recorder),
		Log:           ctrl.Log.WithName("test-appController"),
		Meters:        testcommon.InitAppMeters(),
		SpanHandler:   spanHandlerMock,
		TracerFactory: tf,
		Config:        config.Instance(),
		EvaluationHandler: &evaluationfake.MockEvaluationHandler{
			ReconcileEvaluationsFunc: func(ctx context.Context, phaseCtx context.Context, reconcileObject client.Object, evaluationCreateAttributes evaluation.CreateEvaluationAttributes) ([]apilifecycle.ItemStatus, apicommon.StatusSummary, error) {
				return []apilifecycle.ItemStatus{}, apicommon.StatusSummary{}, nil
			},
		},
	}
	return r, recorder.Events, tr
}

func makeWorkloadVersionWithRef(objectMeta metav1.ObjectMeta, refKind string) *apilifecycle.KeptnWorkloadVersion {
	workloadVersion := &apilifecycle.KeptnWorkloadVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-wli",
			Namespace: "default",
		},
		Spec: apilifecycle.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
				ResourceReference: apilifecycle.ResourceReference{
					UID:  objectMeta.UID,
					Name: objectMeta.Name,
					Kind: refKind,
				},
			},
		},
	}
	return workloadVersion
}

func TestKeptnWorkloadVersionReconciler_checkPreEvaluationStatusOfAppCannotRetrieveApp(t *testing.T) {
	r, _, _ := setupReconciler()

	wv := &apilifecycle.KeptnWorkloadVersion{
		Spec: apilifecycle.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{},
		},
	}

	fakeClient := k8sfake.NewClientBuilder().WithScheme(scheme.Scheme).WithInterceptorFuncs(interceptor.Funcs{
		List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
			return errors.New("unexpected error")
		},
	}).Build()

	r.Client = fakeClient

	requeue, err := r.checkPreEvaluationStatusOfApp(context.TODO(), wv)

	require.True(t, requeue)
	require.Error(t, err)
}

func TestKeptnWorkloadVersionReconciler_checkPreEvaluationStatusOfAppAppNotFound(t *testing.T) {
	r, _, _ := setupReconciler()

	wv := &apilifecycle.KeptnWorkloadVersion{
		Spec: apilifecycle.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
				AppName: "my-unknown-app",
			},
		},
	}

	requeue, err := r.checkPreEvaluationStatusOfApp(context.TODO(), wv)

	require.True(t, requeue)
	require.ErrorIs(t, err, controllererrors.ErrNoMatchingAppVersionFound)
}

func TestKeptnWorkloadVersionReconciler_checkPreEvaluationStatusOfAppAppPreTasksNotFinished(t *testing.T) {
	appVersion := &apilifecycle.KeptnAppVersion{
		ObjectMeta: metav1.ObjectMeta{},
		Spec: apilifecycle.KeptnAppVersionSpec{
			AppName: "my-app",
			KeptnAppSpec: apilifecycle.KeptnAppSpec{
				Version: "1.0",
				Workloads: []apilifecycle.KeptnWorkloadRef{
					{
						Name:    "my-workload",
						Version: "1.0",
					},
				},
			},
		},
	}

	r, _, _ := setupReconciler(appVersion)

	wv := &apilifecycle.KeptnWorkloadVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app-my-workload-1.0",
		},
		Spec: apilifecycle.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
				AppName: "my-app",
				Version: "1.0",
			},
			WorkloadName: "my-app-my-workload",
		},
	}

	requeue, err := r.checkPreEvaluationStatusOfApp(context.TODO(), wv)

	require.True(t, requeue)
	require.Nil(t, err)
}

func TestKeptnWorkloadVersionReconciler_checkPreEvaluationStatusOfAppUpdateTraceID(t *testing.T) {
	appVersion := &apilifecycle.KeptnAppVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app-version",
		},
		Spec: apilifecycle.KeptnAppVersionSpec{
			AppName: "my-app",
			KeptnAppSpec: apilifecycle.KeptnAppSpec{
				Version: "1.0",
				Workloads: []apilifecycle.KeptnWorkloadRef{
					{
						Name:    "my-workload",
						Version: "1.0",
					},
				},
			},
		},
		Status: apilifecycle.KeptnAppVersionStatus{
			PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
			PhaseTraceIDs: map[string]propagation.MapCarrier{
				apicommon.PhaseAppDeployment.ShortName: map[string]string{"traceparent": "parent-id"},
			},
		},
	}

	wv := &apilifecycle.KeptnWorkloadVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app-my-workload-1.0",
		},
		Spec: apilifecycle.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
				AppName: "my-app",
				Version: "1.0",
			},
			WorkloadName: "my-app-my-workload",
		},
	}

	r, _, _ := setupReconciler(appVersion, wv)

	appVersion.Status = apilifecycle.KeptnAppVersionStatus{
		PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
		PhaseTraceIDs: map[string]propagation.MapCarrier{
			apicommon.PhaseAppDeployment.ShortName: map[string]string{"traceparent": "parent-id"},
		},
	}

	err := r.Client.Status().Update(context.TODO(), appVersion)

	require.Nil(t, err)

	requeue, err := r.checkPreEvaluationStatusOfApp(context.TODO(), wv)

	require.False(t, requeue)
	require.Nil(t, err)

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: wv.Name}, wv)

	require.Nil(t, err)

	require.Equal(t, map[string]string{"traceparent": "parent-id"}, wv.Spec.TraceId)
}

func TestKeptnWorkloadVersionReconciler_checkPreEvaluationStatusOfAppUpdateTraceIDWithAppVersionSpecTraceID(t *testing.T) {
	appVersion := &apilifecycle.KeptnAppVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app-version",
		},
		Spec: apilifecycle.KeptnAppVersionSpec{
			AppName: "my-app",
			TraceId: map[string]string{"traceparent": "parent-id"},
			KeptnAppSpec: apilifecycle.KeptnAppSpec{
				Version: "1.0",
				Workloads: []apilifecycle.KeptnWorkloadRef{
					{
						Name:    "my-workload",
						Version: "1.0",
					},
				},
			},
		},
	}

	wv := &apilifecycle.KeptnWorkloadVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app-my-workload-1.0",
		},
		Spec: apilifecycle.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
				AppName: "my-app",
				Version: "1.0",
			},
			WorkloadName: "my-app-my-workload",
		},
	}

	r, _, _ := setupReconciler(appVersion, wv)

	appVersion.Status = apilifecycle.KeptnAppVersionStatus{
		PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
	}

	err := r.Client.Status().Update(context.TODO(), appVersion)

	require.Nil(t, err)

	requeue, err := r.checkPreEvaluationStatusOfApp(context.TODO(), wv)

	require.False(t, requeue)
	require.Nil(t, err)

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: wv.Name}, wv)

	require.Nil(t, err)

	require.Equal(t, map[string]string{"traceparent": "parent-id"}, wv.Spec.TraceId)
}

func TestKeptnWorkloadVersionReconciler_checkPreEvaluationStatusOfAppErrorWhenUpdatingWorkloadVersion(t *testing.T) {
	appVersion := &apilifecycle.KeptnAppVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app-version",
		},
		Spec: apilifecycle.KeptnAppVersionSpec{
			AppName: "my-app",
			TraceId: map[string]string{"traceparent": "parent-id"},
			KeptnAppSpec: apilifecycle.KeptnAppSpec{
				Version: "1.0",
				Workloads: []apilifecycle.KeptnWorkloadRef{
					{
						Name:    "my-workload",
						Version: "1.0",
					},
				},
			},
		},
	}

	wv := &apilifecycle.KeptnWorkloadVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app-my-workload-1.0",
		},
		Spec: apilifecycle.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
				AppName: "my-app",
				Version: "1.0",
			},
			WorkloadName: "my-app-my-workload",
		},
	}

	r, _, _ := setupReconciler(appVersion, wv)

	// inject an error into the fake client to return an error on update
	fakeClient := k8sfake.NewClientBuilder().WithScheme(scheme.Scheme).WithInterceptorFuncs(interceptor.Funcs{
		Update: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
			return errors.New("unexpected error")
		},
	}).WithObjects(wv, appVersion).WithStatusSubresource(appVersion).Build()

	r.Client = fakeClient

	appVersion.Status = apilifecycle.KeptnAppVersionStatus{
		PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
	}

	err := r.Client.Status().Update(context.TODO(), appVersion)

	require.Nil(t, err)

	requeue, err := r.checkPreEvaluationStatusOfApp(context.TODO(), wv)

	require.NotNil(t, err)
	require.True(t, requeue)
}

func TestKeptnWorkloadVersionReconciler_checkPreEvaluationStatusOfAppUpdateMetadata(t *testing.T) {
	appVersion := &apilifecycle.KeptnAppVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app-version",
		},
		Spec: apilifecycle.KeptnAppVersionSpec{
			AppName: "my-app",
			TraceId: map[string]string{"traceparent": "parent-id"},
			KeptnAppSpec: apilifecycle.KeptnAppSpec{
				Version: "1.0",
				Workloads: []apilifecycle.KeptnWorkloadRef{
					{
						Name:    "my-workload",
						Version: "1.0",
					},
				},
			},
			KeptnAppContextSpec: apilifecycle.KeptnAppContextSpec{
				Metadata: map[string]string{
					"test": "testy",
				},
			},
		},
	}

	wv := &apilifecycle.KeptnWorkloadVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app-my-workload-1.0",
		},
		Spec: apilifecycle.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
				AppName: "my-app",
				Version: "1.0",
			},
			TraceId:      map[string]string{"traceparent": "parent-id"},
			WorkloadName: "my-app-my-workload",
		},
	}

	r, _, _ := setupReconciler(appVersion, wv)

	appVersion.Status = apilifecycle.KeptnAppVersionStatus{
		PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
	}

	err := r.Client.Status().Update(context.TODO(), appVersion)

	require.Nil(t, err)

	requeue, err := r.checkPreEvaluationStatusOfApp(context.TODO(), wv)

	require.False(t, requeue)
	require.Nil(t, err)

	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: wv.Name}, wv)

	require.Nil(t, err)

	require.Equal(t, map[string]string{"test": "testy"}, wv.Status.AppContextMetadata)
}

func TestKeptnWorkloadVersionReconciler_checkPreEvaluationStatusOfAppErrorWhenUpdatingWorkloadVersionStatus(t *testing.T) {
	appVersion := &apilifecycle.KeptnAppVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app-version",
		},
		Spec: apilifecycle.KeptnAppVersionSpec{
			AppName: "my-app",
			TraceId: map[string]string{"traceparent": "parent-id"},
			KeptnAppSpec: apilifecycle.KeptnAppSpec{
				Version: "1.0",
				Workloads: []apilifecycle.KeptnWorkloadRef{
					{
						Name:    "my-workload",
						Version: "1.0",
					},
				},
			},
			KeptnAppContextSpec: apilifecycle.KeptnAppContextSpec{
				Metadata: map[string]string{
					"test": "testy",
				},
			},
		},
	}

	wv := &apilifecycle.KeptnWorkloadVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app-my-workload-1.0",
		},
		Spec: apilifecycle.KeptnWorkloadVersionSpec{
			TraceId: map[string]string{"traceparent": "parent-id"},
			KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
				AppName: "my-app",
				Version: "1.0",
			},
			WorkloadName: "my-app-my-workload",
		},
	}

	r, _, _ := setupReconciler(appVersion, wv)

	// inject an error into the fake client to return an error on update
	fakeClient := k8sfake.NewClientBuilder().WithScheme(scheme.Scheme).WithInterceptorFuncs(interceptor.Funcs{
		Update: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
			return errors.New("unexpected error")
		},
	}).WithObjects(wv, appVersion).WithStatusSubresource(appVersion).Build()

	r.Client = fakeClient

	appVersion.Status = apilifecycle.KeptnAppVersionStatus{
		PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
	}

	err := r.Client.Status().Update(context.TODO(), appVersion)

	require.Nil(t, err)

	requeue, err := r.checkPreEvaluationStatusOfApp(context.TODO(), wv)

	require.NotNil(t, err)
	require.True(t, requeue)
}
