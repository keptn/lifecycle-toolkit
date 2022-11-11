package keptnworkloadinstance

import (
	"context"
	klcv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	"github.com/stretchr/testify/require"
	testrequire "github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"testing"
)

func TestKeptnWorkloadInstanceReconciler_IsPodRunning(t *testing.T) {
	p1 := makeNominatedPod("pod1", "node1", v1.PodRunning)
	p2 := makeNominatedPod("pod2", "node1", v1.PodPending)
	podList := &v1.PodList{Items: []v1.Pod{p1, p2}}
	podList2 := &v1.PodList{Items: []v1.Pod{p2}}
	r := &KeptnWorkloadInstanceReconciler{
		Client: fake.NewClientBuilder().WithLists(podList).Build(),
	}
	isPodRunning, err := r.isPodRunning(context.TODO(), klcv1alpha1.ResourceReference{UID: types.UID("pod1")}, "node1")
	testrequire.Nil(t, err)
	if !isPodRunning {
		t.Errorf("Wrong!")
	}

	r2 := &KeptnWorkloadInstanceReconciler{
		Client: fake.NewClientBuilder().WithLists(podList2).Build(),
	}
	isPodRunning, err = r2.isPodRunning(context.TODO(), klcv1alpha1.ResourceReference{UID: types.UID("pod1")}, "node1")
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

func Test_getLatestAppVersion(t *testing.T) {
	type args struct {
		apps *klcv1alpha1.KeptnAppVersionList
		wli  *klcv1alpha1.KeptnWorkloadInstance
	}
	tests := []struct {
		name           string
		args           args
		wantFound      bool
		wantAppVersion klcv1alpha1.KeptnAppVersion
		wantErr        bool
	}{
		{
			name: "app version found",
			args: args{
				apps: &klcv1alpha1.KeptnAppVersionList{
					Items: []klcv1alpha1.KeptnAppVersion{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "my-app",
								Namespace: "default",
							},
							Spec: klcv1alpha1.KeptnAppVersionSpec{
								KeptnAppSpec: klcv1alpha1.KeptnAppSpec{
									Version: "1.0",
									Workloads: []klcv1alpha1.KeptnWorkloadRef{
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
								Name:      "my-app",
								Namespace: "default",
							},
							Spec: klcv1alpha1.KeptnAppVersionSpec{
								KeptnAppSpec: klcv1alpha1.KeptnAppSpec{
									Version: "2.0",
									Workloads: []klcv1alpha1.KeptnWorkloadRef{
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
				wli: &klcv1alpha1.KeptnWorkloadInstance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-workloadinstance",
						Namespace: "default",
					},
					Spec: klcv1alpha1.KeptnWorkloadInstanceSpec{
						KeptnWorkloadSpec: klcv1alpha1.KeptnWorkloadSpec{
							AppName: "my-app",
							Version: "1.0",
						},
						WorkloadName: "my-app-my-workload",
					},
				},
			},
			wantFound: true,
			wantAppVersion: klcv1alpha1.KeptnAppVersion{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "my-app",
					Namespace: "default",
				},
				Spec: klcv1alpha1.KeptnAppVersionSpec{
					KeptnAppSpec: klcv1alpha1.KeptnAppSpec{
						Version: "2.0",
						Workloads: []klcv1alpha1.KeptnWorkloadRef{
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
				apps: &klcv1alpha1.KeptnAppVersionList{
					Items: []klcv1alpha1.KeptnAppVersion{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "my-app",
								Namespace: "default",
							},
							Spec: klcv1alpha1.KeptnAppVersionSpec{
								KeptnAppSpec: klcv1alpha1.KeptnAppSpec{
									Version: "1.0",
									Workloads: []klcv1alpha1.KeptnWorkloadRef{
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
				wli: &klcv1alpha1.KeptnWorkloadInstance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-workloadinstance",
						Namespace: "default",
					},
					Spec: klcv1alpha1.KeptnWorkloadInstanceSpec{
						KeptnWorkloadSpec: klcv1alpha1.KeptnWorkloadSpec{
							AppName: "my-app",
							Version: "1.0",
						},
						WorkloadName: "my-app-my-workload",
					},
				},
			},
			wantFound:      false,
			wantAppVersion: klcv1alpha1.KeptnAppVersion{},
			wantErr:        false,
		},
		{
			name: "app version with invalid version",
			args: args{
				apps: &klcv1alpha1.KeptnAppVersionList{
					Items: []klcv1alpha1.KeptnAppVersion{
						{
							ObjectMeta: metav1.ObjectMeta{
								Name:      "my-app",
								Namespace: "default",
							},
							Spec: klcv1alpha1.KeptnAppVersionSpec{
								KeptnAppSpec: klcv1alpha1.KeptnAppSpec{
									Version: "",
									Workloads: []klcv1alpha1.KeptnWorkloadRef{
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
				wli: &klcv1alpha1.KeptnWorkloadInstance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-workloadinstance",
						Namespace: "default",
					},
					Spec: klcv1alpha1.KeptnWorkloadInstanceSpec{
						KeptnWorkloadSpec: klcv1alpha1.KeptnWorkloadSpec{
							AppName: "my-app",
							Version: "1.0",
						},
						WorkloadName: "my-app-my-workload",
					},
				},
			},
			wantFound:      false,
			wantAppVersion: klcv1alpha1.KeptnAppVersion{},
			wantErr:        true,
		},
		{
			name: "app version list empty",
			args: args{
				apps: &klcv1alpha1.KeptnAppVersionList{
					Items: []klcv1alpha1.KeptnAppVersion{},
				},
				wli: &klcv1alpha1.KeptnWorkloadInstance{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-workloadinstance",
						Namespace: "default",
					},
					Spec: klcv1alpha1.KeptnWorkloadInstanceSpec{
						KeptnWorkloadSpec: klcv1alpha1.KeptnWorkloadSpec{
							AppName: "my-app",
							Version: "1.0",
						},
						WorkloadName: "my-app-my-workload",
					},
				},
			},
			wantFound:      false,
			wantAppVersion: klcv1alpha1.KeptnAppVersion{},
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
