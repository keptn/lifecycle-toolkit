package keptnworkloadversion

import (
	"context"
	"errors"
	"fmt"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/phase"
	phasefake "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/phase/fake"
	"strings"
	"testing"
	"time"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	klcv1alpha4 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha4"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/evaluation"
	evaluationfake "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/evaluation/fake"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/eventsender"
	schedulinggatesfake "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/schedulinggates/fake"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/telemetry"
	telemetryfake "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/telemetry/fake"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/testcommon"
	controllererrors "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/errors"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	testrequire "github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
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
	testrequire.Nil(t, err)
	testrequire.Equal(t, apicommon.StateProgressing, keptnState)
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
	testrequire.NotNil(t, err)
	testrequire.Equal(t, apicommon.StateUnknown, keptnState)
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
	testrequire.Nil(t, err)
	testrequire.Equal(t, apicommon.StateProgressing, keptnState)
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
	testrequire.NotNil(t, err)
	testrequire.Equal(t, apicommon.StateUnknown, keptnState)
}

func TestKeptnWorkloadVersionReconciler_reconcileDeployment_FailedDaemonSet(t *testing.T) {

	daemonSetFail := makeDaemonSet("mystat", "default", 1, 0)
	workloadVersion := makeWorkloadVersionWithRef(daemonSetFail.ObjectMeta, "DaemonSet")

	fakeClient := testcommon.NewTestClient(daemonSetFail, workloadVersion)

	r := &KeptnWorkloadVersionReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadVersion)
	testrequire.Nil(t, err)
	testrequire.Equal(t, apicommon.StateProgressing, keptnState)
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
	testrequire.NotNil(t, err)
	testrequire.Equal(t, apicommon.StateUnknown, keptnState)
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
	testrequire.Nil(t, err)
	testrequire.Equal(t, apicommon.StateSucceeded, keptnState)
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
	testrequire.Nil(t, err)
	testrequire.Equal(t, apicommon.StateSucceeded, keptnState)
}

func TestKeptnWorkloadVersionReconciler_reconcileDeployment_ReadyDaemonSet(t *testing.T) {

	daemonSet := makeDaemonSet("mystat", "default", 1, 1)
	workloadVersion := makeWorkloadVersionWithRef(daemonSet.ObjectMeta, "DaemonSet")

	fakeClient := testcommon.NewTestClient(daemonSet, workloadVersion)

	r := &KeptnWorkloadVersionReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadVersion)
	testrequire.Nil(t, err)
	testrequire.Equal(t, apicommon.StateSucceeded, keptnState)
}

func TestKeptnWorkloadVersionReconciler_reconcileDeployment_UnsupportedReferenceKind(t *testing.T) {

	workloadVersion := makeWorkloadVersionWithRef(metav1.ObjectMeta{}, "Unknown")
	fakeClient := testcommon.NewTestClient(workloadVersion)
	r := &KeptnWorkloadVersionReconciler{
		Client: fakeClient,
	}

	keptnState, err := r.reconcileDeployment(context.TODO(), workloadVersion)
	testrequire.ErrorIs(t, err, controllererrors.ErrUnsupportedWorkloadVersionResourceReference)
	testrequire.Equal(t, apicommon.StateUnknown, keptnState)
}

func TestKeptnWorkloadVersionReconciler_IsPodRunning(t *testing.T) {
	p1 := makeNominatedPod("pod1", "node1", v1.PodRunning)
	p2 := makeNominatedPod("pod2", "node1", v1.PodPending)
	podList := &v1.PodList{Items: []v1.Pod{p1, p2}}
	podList2 := &v1.PodList{Items: []v1.Pod{p2}}
	r := &KeptnWorkloadVersionReconciler{
		Client: k8sfake.NewClientBuilder().WithLists(podList).Build(),
	}
	isPodRunning, err := r.isPodRunning(context.TODO(), klcv1alpha3.ResourceReference{UID: "pod1"}, "node1")
	testrequire.Nil(t, err)
	if !isPodRunning {
		t.Errorf("Wrong!")
	}

	r2 := &KeptnWorkloadVersionReconciler{
		Client: k8sfake.NewClientBuilder().WithLists(podList2).Build(),
	}
	isPodRunning, err = r2.isPodRunning(context.TODO(), klcv1alpha3.ResourceReference{UID: "pod1"}, "node1")
	testrequire.Nil(t, err)
	if isPodRunning {
		t.Errorf("Wrong!")
	}

}

func makeNominatedPod(podName string, nodeName string, phase v1.PodPhase) v1.Pod {
	return v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: nodeName,
			Name:      podName,
			UID:       types.UID(podName),
		},
		Status: v1.PodStatus{
			Phase:             phase,
			NominatedNodeName: nodeName,
		},
	}
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
		wli            *klcv1alpha4.KeptnWorkloadVersion
		list           *klcv1alpha3.KeptnAppVersionList
		wantFound      bool
		wantAppVersion klcv1alpha3.KeptnAppVersion
		wantErr        bool
	}{
		{
			name: "no appVersions",
			wli: &klcv1alpha4.KeptnWorkloadVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-workloadVersion",
					Namespace: "default",
				},
				Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
					KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
						AppName: "my-app",
						Version: "1.0",
					},
					WorkloadName: "my-app-my-workload",
				},
			},
			list:           &klcv1alpha3.KeptnAppVersionList{},
			wantFound:      false,
			wantAppVersion: klcv1alpha3.KeptnAppVersion{},
			wantErr:        false,
		},
		{
			name: "appVersion found",
			wli: &klcv1alpha4.KeptnWorkloadVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-workloadVersion",
					Namespace: "default",
				},
				Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
					KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
						AppName: "my-app",
						Version: "1.0",
					},
					WorkloadName: "my-app-my-workload",
				},
			},
			list: &klcv1alpha3.KeptnAppVersionList{
				Items: []klcv1alpha3.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:              "my-app",
							Namespace:         "default",
							CreationTimestamp: metav1.Time{Time: now},
						},
						Spec: klcv1alpha3.KeptnAppVersionSpec{
							KeptnAppSpec: klcv1alpha3.KeptnAppSpec{
								Version: "1.0",
								Workloads: []klcv1alpha3.KeptnWorkloadRef{
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
						Spec: klcv1alpha3.KeptnAppVersionSpec{
							KeptnAppSpec: klcv1alpha3.KeptnAppSpec{
								Version: "2.0",
								Workloads: []klcv1alpha3.KeptnWorkloadRef{
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
			wantAppVersion: klcv1alpha3.KeptnAppVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name:            "my-app2",
					Namespace:       "default",
					ResourceVersion: "999",
				},
				Spec: klcv1alpha3.KeptnAppVersionSpec{
					KeptnAppSpec: klcv1alpha3.KeptnAppSpec{
						Version: "2.0",
						Workloads: []klcv1alpha3.KeptnWorkloadRef{
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
			wli: &klcv1alpha4.KeptnWorkloadVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-workloadVersion",
					Namespace: "default",
				},
				Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
					KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
						AppName: "my-app",
						Version: "1.0",
					},
					WorkloadName: "my-app-my-workload",
				},
			},
			list: &klcv1alpha3.KeptnAppVersionList{
				Items: []klcv1alpha3.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "my-app",
							Namespace: "default",
						},
						Spec: klcv1alpha3.KeptnAppVersionSpec{
							KeptnAppSpec: klcv1alpha3.KeptnAppSpec{
								Version: "1.0",
								Workloads: []klcv1alpha3.KeptnWorkloadRef{
									{
										Name:    "my-workload",
										Version: "1.0",
									},
								},
							},
							AppName: "my-app",
						},
						Status: klcv1alpha3.KeptnAppVersionStatus{
							Status: apicommon.StateDeprecated,
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "my-app2",
							Namespace: "default",
						},
						Spec: klcv1alpha3.KeptnAppVersionSpec{
							KeptnAppSpec: klcv1alpha3.KeptnAppSpec{
								Version: "2.0",
								Workloads: []klcv1alpha3.KeptnWorkloadRef{
									{
										Name:    "my-workload",
										Version: "1.0",
									},
								},
							},
							AppName: "my-app",
						},
						Status: klcv1alpha3.KeptnAppVersionStatus{
							Status: apicommon.StateDeprecated,
						},
					},
				},
			},
			wantFound:      false,
			wantAppVersion: klcv1alpha3.KeptnAppVersion{},
			wantErr:        false,
		},
		{
			name: "no workload for appversion",
			wli: &klcv1alpha4.KeptnWorkloadVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-workloadVersion3",
					Namespace: "default",
				},
				Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
					KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
						AppName: "my-app333",
						Version: "1.0.0",
					},
					WorkloadName: "my-app-my-workload",
				},
			},
			list: &klcv1alpha3.KeptnAppVersionList{
				Items: []klcv1alpha3.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "my-app",
							Namespace: "default",
						},
						Spec: klcv1alpha3.KeptnAppVersionSpec{
							KeptnAppSpec: klcv1alpha3.KeptnAppSpec{
								Version: "1.0",
								Workloads: []klcv1alpha3.KeptnWorkloadRef{
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
						Spec: klcv1alpha3.KeptnAppVersionSpec{
							KeptnAppSpec: klcv1alpha3.KeptnAppSpec{
								Version: "2.0",
								Workloads: []klcv1alpha3.KeptnWorkloadRef{
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
			wantAppVersion: klcv1alpha3.KeptnAppVersion{},
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := klcv1alpha3.AddToScheme(scheme.Scheme)
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
		apps *klcv1alpha3.KeptnAppVersionList
		wli  *klcv1alpha4.KeptnWorkloadVersion
	}
	tests := []struct {
		name           string
		args           args
		wantFound      bool
		wantAppVersion klcv1alpha3.KeptnAppVersion
		wantErr        bool
	}{
		{
			name: "app version found",
			args: args{
				apps: &klcv1alpha3.KeptnAppVersionList{
					Items: []klcv1alpha3.KeptnAppVersion{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:              "my-app",
								Namespace:         "default",
								CreationTimestamp: metav1.Time{Time: now},
							},
							Spec: klcv1alpha3.KeptnAppVersionSpec{
								KeptnAppSpec: klcv1alpha3.KeptnAppSpec{
									Version: "1.0",
									Workloads: []klcv1alpha3.KeptnWorkloadRef{
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
							Spec: klcv1alpha3.KeptnAppVersionSpec{
								KeptnAppSpec: klcv1alpha3.KeptnAppSpec{
									Version: "2.0",
									Workloads: []klcv1alpha3.KeptnWorkloadRef{
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
				wli: &klcv1alpha4.KeptnWorkloadVersion{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-workloadVersion",
						Namespace: "default",
					},
					Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
						KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
							AppName: "my-app",
							Version: "1.0",
						},
						WorkloadName: "my-app-my-workload",
					},
				},
			},
			wantFound: true,
			wantAppVersion: klcv1alpha3.KeptnAppVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name:              "my-app",
					Namespace:         "default",
					CreationTimestamp: metav1.Time{Time: now.Add(5 * time.Second)},
				},
				Spec: klcv1alpha3.KeptnAppVersionSpec{
					KeptnAppSpec: klcv1alpha3.KeptnAppSpec{
						Version: "2.0",
						Workloads: []klcv1alpha3.KeptnWorkloadRef{
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
				apps: &klcv1alpha3.KeptnAppVersionList{
					Items: []klcv1alpha3.KeptnAppVersion{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "my-app",
								Namespace: "default",
							},
							Spec: klcv1alpha3.KeptnAppVersionSpec{
								KeptnAppSpec: klcv1alpha3.KeptnAppSpec{
									Version: "1.0",
									Workloads: []klcv1alpha3.KeptnWorkloadRef{
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
				wli: &klcv1alpha4.KeptnWorkloadVersion{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-workloadVersion",
						Namespace: "default",
					},
					Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
						KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
							AppName: "my-app",
							Version: "1.0",
						},
						WorkloadName: "my-app-my-workload",
					},
				},
			},
			wantFound:      false,
			wantAppVersion: klcv1alpha3.KeptnAppVersion{},
			wantErr:        false,
		},
		{
			name: "app version list empty",
			args: args{
				apps: &klcv1alpha3.KeptnAppVersionList{
					Items: []klcv1alpha3.KeptnAppVersion{},
				},
				wli: &klcv1alpha4.KeptnWorkloadVersion{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-workloadVersion",
						Namespace: "default",
					},
					Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
						KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
							AppName: "my-app",
							Version: "1.0",
						},
						WorkloadName: "my-app-my-workload",
					},
				},
			},
			wantFound:      false,
			wantAppVersion: klcv1alpha3.KeptnAppVersion{},
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

	wi := &klcv1alpha4.KeptnWorkloadVersion{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "some-wi",
			Namespace: testNamespace,
		},
		Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
				AppName: "some-app",
				Version: "1.0.0",
			},
			WorkloadName:    "some-app-some-workload",
			PreviousVersion: "",
			TraceId:         nil,
		},
		Status: klcv1alpha4.KeptnWorkloadVersionStatus{
			DeploymentStatus:               apicommon.StateSucceeded,
			PreDeploymentStatus:            apicommon.StateSucceeded,
			PostDeploymentStatus:           apicommon.StateSucceeded,
			PreDeploymentEvaluationStatus:  apicommon.StateSucceeded,
			PostDeploymentEvaluationStatus: apicommon.StateSucceeded,
			CurrentPhase:                   apicommon.PhaseWorkloadPostEvaluation.ShortName,
			Status:                         apicommon.StateSucceeded,
			StartTime:                      metav1.Time{},
			EndTime:                        metav1.Time{},
		},
	}

	app := testcommon.ReturnAppVersion(
		testNamespace,
		"some-app",
		"1.0.0",
		[]klcv1alpha3.KeptnWorkloadRef{
			{
				Name:    "some-workload",
				Version: "1.0.0",
			},
		},
		klcv1alpha3.KeptnAppVersionStatus{
			PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
		},
	)
	r, eventChannel, _ := setupReconciler(wi, app)
	r.SchedulingGatesHandler = &schedulinggatesfake.ISchedulingGatesHandlerMock{
		EnabledFunc: func() bool {
			return false
		},
	}

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
}

func TestKeptnWorkloadVersionReconciler_ReconcileReachCompletion_SchedulingGates(t *testing.T) {

	testNamespace := "some-ns"

	wi := &klcv1alpha4.KeptnWorkloadVersion{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "some-wi",
			Namespace: testNamespace,
		},
		Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
				AppName: "some-app",
				Version: "1.0.0",
			},
			WorkloadName:    "some-app-some-workload",
			PreviousVersion: "",
			TraceId:         nil,
		},
		Status: klcv1alpha4.KeptnWorkloadVersionStatus{
			DeploymentStatus:               apicommon.StateSucceeded,
			PreDeploymentStatus:            apicommon.StateSucceeded,
			PostDeploymentStatus:           apicommon.StateSucceeded,
			PreDeploymentEvaluationStatus:  apicommon.StateSucceeded,
			PostDeploymentEvaluationStatus: apicommon.StateSucceeded,
			CurrentPhase:                   apicommon.PhaseWorkloadPostEvaluation.ShortName,
			Status:                         apicommon.StateSucceeded,
			StartTime:                      metav1.Time{},
			EndTime:                        metav1.Time{},
		},
	}

	app := testcommon.ReturnAppVersion(
		testNamespace,
		"some-app",
		"1.0.0",
		[]klcv1alpha3.KeptnWorkloadRef{
			{
				Name:    "some-workload",
				Version: "1.0.0",
			},
		},
		klcv1alpha3.KeptnAppVersionStatus{
			PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
		},
	)

	schedulingGatesMock := &schedulinggatesfake.ISchedulingGatesHandlerMock{
		RemoveGatesFunc: func(ctx context.Context, workloadVersion *klcv1alpha4.KeptnWorkloadVersion) error {
			return nil
		},
		EnabledFunc: func() bool {
			return true
		},
	}
	r, eventChannel, _ := setupReconciler(wi, app)
	r.SchedulingGatesHandler = schedulingGatesMock

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Namespace: testNamespace,
			Name:      "some-wi",
		},
	}

	result, err := r.Reconcile(context.TODO(), req)

	require.Len(t, schedulingGatesMock.RemoveGatesCalls(), 1)
	require.Nil(t, err)

	// do not requeue since we reached completion
	require.False(t, result.Requeue)

	// here we do not expect an event about the application preEvaluation being finished since that  will have been sent in
	// one of the previous reconciliation loops that lead to the first phase being reached
	expectedEvents := []string{
		"CompletedFinished",
	}

	for _, e := range expectedEvents {
		select {
		case event := <-eventChannel:
			assert.Equal(t, strings.Contains(event, req.Name), true, "wrong workloadVersion")
			assert.Equal(t, strings.Contains(event, req.Namespace), true, "wrong namespace")
			assert.Equal(t, strings.Contains(event, e), true, fmt.Sprintf("no %s found in %s", e, event))
		case <-time.After(5 * time.Second):
			t.Error("Didn't receive the cloud event")
		}
	}
}

func TestKeptnWorkloadVersionReconciler_RemoveGates_fail(t *testing.T) {

	testNamespace := "some-ns"

	wi := &klcv1alpha4.KeptnWorkloadVersion{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "some-wi",
			Namespace: testNamespace,
		},
		Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
				AppName: "some-app",
				Version: "1.0.0",
			},
			WorkloadName:    "some-app-some-workload",
			PreviousVersion: "",
			TraceId:         nil,
		},
		Status: klcv1alpha4.KeptnWorkloadVersionStatus{
			DeploymentStatus:               apicommon.StateSucceeded,
			PreDeploymentStatus:            apicommon.StateSucceeded,
			PostDeploymentStatus:           apicommon.StateSucceeded,
			PreDeploymentEvaluationStatus:  apicommon.StateSucceeded,
			PostDeploymentEvaluationStatus: apicommon.StateSucceeded,
			CurrentPhase:                   apicommon.PhaseWorkloadPostEvaluation.ShortName,
			Status:                         apicommon.StateSucceeded,
			StartTime:                      metav1.Time{},
			EndTime:                        metav1.Time{},
		},
	}

	app := testcommon.ReturnAppVersion(
		testNamespace,
		"some-app",
		"1.0.0",
		[]klcv1alpha3.KeptnWorkloadRef{
			{
				Name:    "some-workload",
				Version: "1.0.0",
			},
		},
		klcv1alpha3.KeptnAppVersionStatus{
			PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
		},
	)
	r, _, _ := setupReconciler(wi, app)
	r.SchedulingGatesHandler = &schedulinggatesfake.ISchedulingGatesHandlerMock{
		RemoveGatesFunc: func(ctx context.Context, workloadVersion *klcv1alpha4.KeptnWorkloadVersion) error {
			return fmt.Errorf("err")
		},
		EnabledFunc: func() bool {
			return true
		},
	}

	req := ctrl.Request{
		NamespacedName: types.NamespacedName{
			Namespace: testNamespace,
			Name:      "some-wi",
		},
	}

	result, err := r.Reconcile(context.TODO(), req)

	require.NotNil(t, err)

	// do not requeue since we reached completion
	require.True(t, result.Requeue)
}

func TestKeptnWorkloadVersionReconciler_ReconcileFailed(t *testing.T) {

	testNamespace := "some-ns"

	wi := &klcv1alpha4.KeptnWorkloadVersion{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "some-wi",
			Namespace: testNamespace,
		},
		Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
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
		Status: klcv1alpha4.KeptnWorkloadVersionStatus{
			DeploymentStatus:               apicommon.StatePending,
			PreDeploymentStatus:            apicommon.StateProgressing,
			PostDeploymentStatus:           apicommon.StatePending,
			PreDeploymentEvaluationStatus:  apicommon.StatePending,
			PostDeploymentEvaluationStatus: apicommon.StatePending,
			CurrentPhase:                   apicommon.PhaseWorkloadPreDeployment.ShortName,
			Status:                         apicommon.StateProgressing,
			PreDeploymentTaskStatus: []klcv1alpha3.ItemStatus{
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
		[]klcv1alpha3.KeptnWorkloadRef{
			{
				Name:    "some-workload",
				Version: "1.0.0",
			},
		},
		klcv1alpha3.KeptnAppVersionStatus{
			PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
		},
	)

	r, _, _ := setupReconciler(app, wi)

	r.PhaseHandler = &phasefake.MockHandler{HandlePhaseFunc: func(ctx context.Context, ctxTrace context.Context, tracer trace.Tracer, reconcileObject client.Object, phaseMoqParam apicommon.KeptnPhaseType, reconcilePhase func(phaseCtx context.Context) (apicommon.KeptnState, error)) (phase.PhaseResult, error) {
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

	wi := &klcv1alpha4.KeptnWorkloadVersion{
		TypeMeta: metav1.TypeMeta{},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "some-wi",
			Namespace: testNamespace,
		},
		Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
				AppName:                  "some-app",
				Version:                  "1.0.0",
				PreDeploymentEvaluations: []string{"my-pre-evaluation"},
			},
			WorkloadName:    "some-app-some-workload",
			PreviousVersion: "",
			TraceId:         nil,
		},
		Status: klcv1alpha4.KeptnWorkloadVersionStatus{
			CurrentPhase: apicommon.PhaseWorkloadPreDeployment.ShortName,
			StartTime:    metav1.Time{},
			EndTime:      metav1.Time{},
		},
	}

	app := testcommon.ReturnAppVersion(
		testNamespace,
		"some-app",
		"1.0.0",
		[]klcv1alpha3.KeptnWorkloadRef{
			{
				Name:    "some-workload",
				Version: "1.0.0",
			},
		},
		klcv1alpha3.KeptnAppVersionStatus{
			PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
		},
	)

	r, _, _ := setupReconciler(wi, app)

	mockEvaluationHandler := r.EvaluationHandler.(*evaluationfake.MockEvaluationHandler)

	mockEvaluationHandler.ReconcileEvaluationsFunc = func(ctx context.Context, phaseCtx context.Context, reconcileObject client.Object, evaluationCreateAttributes evaluation.CreateEvaluationAttributes) ([]klcv1alpha3.ItemStatus, apicommon.StatusSummary, error) {
		return nil, apicommon.StatusSummary{}, errors.New("unexpected error")
	}

	mockPhaseHandler := &phasefake.MockHandler{
		HandlePhaseFunc: func(ctx context.Context, ctxTrace context.Context, tracer trace.Tracer, reconcileObject client.Object, phaseMoqParam apicommon.KeptnPhaseType, reconcilePhase func(phaseCtx context.Context) (apicommon.KeptnState, error)) (phase.PhaseResult, error) {
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

	tf := &telemetryfake.TracerFactoryMock{GetTracerFunc: func(name string) trace.Tracer {
		return tr
	}}

	fakeClient := testcommon.NewTestClient(objs...)

	recorder := record.NewFakeRecorder(100)

	r := &KeptnWorkloadVersionReconciler{
		Client:        fakeClient,
		Scheme:        scheme.Scheme,
		EventSender:   eventsender.NewK8sSender(recorder),
		Log:           ctrl.Log.WithName("test-appController"),
		Meters:        testcommon.InitAppMeters(),
		SpanHandler:   &telemetry.Handler{},
		TracerFactory: tf,
		EvaluationHandler: &evaluationfake.MockEvaluationHandler{
			ReconcileEvaluationsFunc: func(ctx context.Context, phaseCtx context.Context, reconcileObject client.Object, evaluationCreateAttributes evaluation.CreateEvaluationAttributes) ([]klcv1alpha3.ItemStatus, apicommon.StatusSummary, error) {
				return []klcv1alpha3.ItemStatus{}, apicommon.StatusSummary{}, nil
			},
		},
	}
	return r, recorder.Events, tr
}

func makeWorkloadVersionWithRef(objectMeta metav1.ObjectMeta, refKind string) *klcv1alpha4.KeptnWorkloadVersion {
	workloadVersion := &klcv1alpha4.KeptnWorkloadVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-wli",
			Namespace: "default",
		},
		Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
				ResourceReference: klcv1alpha3.ResourceReference{
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

	wv := &klcv1alpha4.KeptnWorkloadVersion{
		Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{},
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

	wv := &klcv1alpha4.KeptnWorkloadVersion{
		Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
				AppName: "my-unknown-app",
			},
		},
	}

	requeue, err := r.checkPreEvaluationStatusOfApp(context.TODO(), wv)

	require.True(t, requeue)
	require.ErrorIs(t, err, controllererrors.ErrNoMatchingAppVersionFound)
}

func TestKeptnWorkloadVersionReconciler_checkPreEvaluationStatusOfAppAppPreTasksNotFinished(t *testing.T) {
	appVersion := &klcv1alpha3.KeptnAppVersion{
		ObjectMeta: metav1.ObjectMeta{},
		Spec: klcv1alpha3.KeptnAppVersionSpec{
			AppName: "my-app",
			KeptnAppSpec: klcv1alpha3.KeptnAppSpec{
				Version: "1.0",
				Workloads: []klcv1alpha3.KeptnWorkloadRef{
					{
						Name:    "my-workload",
						Version: "1.0",
					},
				},
			},
		},
	}

	r, _, _ := setupReconciler(appVersion)

	wv := &klcv1alpha4.KeptnWorkloadVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app-my-workload-1.0",
		},
		Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
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
	appVersion := &klcv1alpha3.KeptnAppVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app-version",
		},
		Spec: klcv1alpha3.KeptnAppVersionSpec{
			AppName: "my-app",
			KeptnAppSpec: klcv1alpha3.KeptnAppSpec{
				Version: "1.0",
				Workloads: []klcv1alpha3.KeptnWorkloadRef{
					{
						Name:    "my-workload",
						Version: "1.0",
					},
				},
			},
		},
		Status: klcv1alpha3.KeptnAppVersionStatus{
			PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
			PhaseTraceIDs: map[string]propagation.MapCarrier{
				apicommon.PhaseAppDeployment.ShortName: map[string]string{"traceparent": "parent-id"},
			},
		},
	}

	wv := &klcv1alpha4.KeptnWorkloadVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app-my-workload-1.0",
		},
		Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
				AppName: "my-app",
				Version: "1.0",
			},
			WorkloadName: "my-app-my-workload",
		},
	}

	r, _, _ := setupReconciler(appVersion, wv)

	appVersion.Status = klcv1alpha3.KeptnAppVersionStatus{
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
	appVersion := &klcv1alpha3.KeptnAppVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app-version",
		},
		Spec: klcv1alpha3.KeptnAppVersionSpec{
			AppName: "my-app",
			TraceId: map[string]string{"traceparent": "parent-id"},
			KeptnAppSpec: klcv1alpha3.KeptnAppSpec{
				Version: "1.0",
				Workloads: []klcv1alpha3.KeptnWorkloadRef{
					{
						Name:    "my-workload",
						Version: "1.0",
					},
				},
			},
		},
	}

	wv := &klcv1alpha4.KeptnWorkloadVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app-my-workload-1.0",
		},
		Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
				AppName: "my-app",
				Version: "1.0",
			},
			WorkloadName: "my-app-my-workload",
		},
	}

	r, _, _ := setupReconciler(appVersion, wv)

	appVersion.Status = klcv1alpha3.KeptnAppVersionStatus{
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
	appVersion := &klcv1alpha3.KeptnAppVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app-version",
		},
		Spec: klcv1alpha3.KeptnAppVersionSpec{
			AppName: "my-app",
			TraceId: map[string]string{"traceparent": "parent-id"},
			KeptnAppSpec: klcv1alpha3.KeptnAppSpec{
				Version: "1.0",
				Workloads: []klcv1alpha3.KeptnWorkloadRef{
					{
						Name:    "my-workload",
						Version: "1.0",
					},
				},
			},
		},
	}

	wv := &klcv1alpha4.KeptnWorkloadVersion{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-app-my-workload-1.0",
		},
		Spec: klcv1alpha4.KeptnWorkloadVersionSpec{
			KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
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

	appVersion.Status = klcv1alpha3.KeptnAppVersionStatus{
		PreDeploymentEvaluationStatus: apicommon.StateSucceeded,
	}

	err := r.Client.Status().Update(context.TODO(), appVersion)

	require.Nil(t, err)

	requeue, err := r.checkPreEvaluationStatusOfApp(context.TODO(), wv)

	require.NotNil(t, err)
	require.True(t, requeue)
}
