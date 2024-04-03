package schedulinggates

import (
	"context"
	"errors"
	"testing"
	"time"

	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	controllerruntime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	k8sfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/client/interceptor"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestSchedulingGatesReconciler_Reconcile(t *testing.T) {
	podMeta := metav1.ObjectMeta{
		Name:      "my-pod",
		Namespace: "my-namespace",
		OwnerReferences: []metav1.OwnerReference{
			{
				Kind: "ReplicaSet",
				Name: "my-replica-set",
				UID:  "some-uid",
			},
		},
	}

	req := controllerruntime.Request{
		NamespacedName: types.NamespacedName{
			Namespace: podMeta.Namespace,
			Name:      podMeta.Name,
		},
	}

	type args struct {
		ctx context.Context
		req controllerruntime.Request
	}
	tests := []struct {
		name               string
		objects            []client.Object
		pod                client.Object
		args               args
		want               controllerruntime.Result
		lookupError        bool
		updateError        bool
		expectGatesRemoved bool
		wantErr            bool
	}{
		{
			name:    "no pod found",
			objects: []client.Object{},
			args: args{
				ctx: context.TODO(),
				req: req,
			},
			want:    controllerruntime.Result{},
			wantErr: false,
		},
		{
			name: "no owner references",
			objects: []client.Object{
				&v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-pod",
						Namespace: "my-namespace",
					},
					Spec: v1.PodSpec{
						SchedulingGates: []v1.PodSchedulingGate{
							{
								Name: apicommon.KeptnGate,
							},
						},
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				req: req,
			},
			want:    controllerruntime.Result{},
			wantErr: false,
		},
		{
			name: "no related WorkloadVersions",
			objects: []client.Object{
				&v1.Pod{
					ObjectMeta: podMeta,
					Spec: v1.PodSpec{
						SchedulingGates: []v1.PodSchedulingGate{
							{
								Name: apicommon.KeptnGate,
							},
						},
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				req: req,
			},
			want: controllerruntime.Result{
				RequeueAfter: 10 * time.Second,
			},
			wantErr: false,
		},
		{
			name: "error when executing list request",
			objects: []client.Object{
				&v1.Pod{
					ObjectMeta: podMeta,
					Spec: v1.PodSpec{
						SchedulingGates: []v1.PodSchedulingGate{
							{
								Name: apicommon.KeptnGate,
							},
						},
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				req: req,
			},
			lookupError: true,
			want:        controllerruntime.Result{},
			wantErr:     true,
		},
		{
			name: "related WorkloadVersion is completed",
			objects: []client.Object{
				&apilifecycle.KeptnWorkloadVersion{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-wlv",
						Namespace: "my-namespace",
					},
					Spec: apilifecycle.KeptnWorkloadVersionSpec{
						KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
							ResourceReference: apilifecycle.ResourceReference{
								UID: podMeta.OwnerReferences[0].UID,
							},
						},
					},
					Status: apilifecycle.KeptnWorkloadVersionStatus{DeploymentStatus: apicommon.StateSucceeded},
				},
				&apilifecycle.KeptnWorkloadVersion{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-other-wlv",
						Namespace: "my-namespace",
					},
					Spec: apilifecycle.KeptnWorkloadVersionSpec{},
				},
				&apilifecycle.KeptnWorkloadVersion{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-wlv",
						Namespace: "my-other-namespace",
					},
					Spec: apilifecycle.KeptnWorkloadVersionSpec{
						KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
							ResourceReference: apilifecycle.ResourceReference{
								UID: podMeta.OwnerReferences[0].UID,
							},
						},
					},
				},
				&v1.Pod{
					ObjectMeta: podMeta,
					Spec: v1.PodSpec{
						SchedulingGates: []v1.PodSchedulingGate{
							{
								Name: apicommon.KeptnGate,
							},
						},
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				req: req,
			},
			want:               controllerruntime.Result{},
			wantErr:            false,
			expectGatesRemoved: true,
		},
		{
			name: "related WorkloadVersion is completed - error during update",
			objects: []client.Object{
				&apilifecycle.KeptnWorkloadVersion{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-wlv",
						Namespace: "my-namespace",
					},
					Spec: apilifecycle.KeptnWorkloadVersionSpec{
						KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
							ResourceReference: apilifecycle.ResourceReference{
								UID: podMeta.OwnerReferences[0].UID,
							},
						},
					},
					Status: apilifecycle.KeptnWorkloadVersionStatus{DeploymentStatus: apicommon.StateSucceeded},
				},
				&v1.Pod{
					ObjectMeta: podMeta,
					Spec: v1.PodSpec{
						SchedulingGates: []v1.PodSchedulingGate{
							{
								Name: apicommon.KeptnGate,
							},
						},
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				req: req,
			},
			updateError:        true,
			want:               controllerruntime.Result{},
			wantErr:            true,
			expectGatesRemoved: false,
		},
		{
			name: "related WorkloadVersion is not completed",
			objects: []client.Object{
				&apilifecycle.KeptnWorkloadVersion{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "my-wlv",
						Namespace: "my-namespace",
					},
					Spec: apilifecycle.KeptnWorkloadVersionSpec{
						KeptnWorkloadSpec: apilifecycle.KeptnWorkloadSpec{
							ResourceReference: apilifecycle.ResourceReference{
								UID: podMeta.OwnerReferences[0].UID,
							},
						},
					},
					Status: apilifecycle.KeptnWorkloadVersionStatus{DeploymentStatus: apicommon.StatePending},
				},
				&v1.Pod{
					ObjectMeta: podMeta,
					Spec: v1.PodSpec{
						SchedulingGates: []v1.PodSchedulingGate{
							{
								Name: apicommon.KeptnGate,
							},
						},
					},
				},
			},
			args: args{
				ctx: context.TODO(),
				req: req,
			},
			want:               controllerruntime.Result{RequeueAfter: 10 * time.Second},
			wantErr:            false,
			expectGatesRemoved: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apilifecycle.AddToScheme(scheme.Scheme)
			require.Nil(t, err)
			opts := zap.Options{
				Development: true,
			}
			controllerruntime.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
			mockClient := k8sfake.
				NewClientBuilder().
				WithScheme(scheme.Scheme).
				WithObjects(tt.objects...).
				WithStatusSubresource(&apilifecycle.KeptnWorkloadVersion{}).
				WithIndex(&apilifecycle.KeptnWorkloadVersion{}, ".spec.resourceReference.uid", func(object client.Object) []string {
					return common.KeptnWorkloadVersionResourceRefUIDIndexFunc(object)
				}).
				WithInterceptorFuncs(
					interceptor.Funcs{
						List: func(ctx context.Context, client client.WithWatch, list client.ObjectList, opts ...client.ListOption) error {
							if tt.lookupError {
								return errors.New("unexpected error")
							}
							return client.List(ctx, list, opts...)
						},
						Update: func(ctx context.Context, client client.WithWatch, obj client.Object, opts ...client.UpdateOption) error {
							if tt.updateError {
								return errors.New("unexpected error")
							}
							return client.Update(ctx, obj, opts...)
						}}).
				Build()

			for _, obj := range tt.objects {
				kwv, ok := obj.(*apilifecycle.KeptnWorkloadVersion)
				if ok && kwv.Status.DeploymentStatus != "" {
					err := mockClient.Status().Update(context.TODO(), kwv)
					require.Nil(t, err)
				}
			}

			r := &SchedulingGatesReconciler{
				Client: mockClient,
				Scheme: scheme.Scheme,
				Log:    controllerruntime.Log.WithName("test-appController"),
			}

			got, err := r.Reconcile(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.NotNil(t, err)
			} else {
				require.Nil(t, err)
			}
			require.Equal(t, tt.want, got)

			resultingPod := &v1.Pod{}

			if tt.expectGatesRemoved {
				err = mockClient.Get(context.TODO(), types.NamespacedName{
					Namespace: podMeta.Namespace,
					Name:      podMeta.Name,
				}, resultingPod)

				require.Nil(t, err)

				require.Empty(t, resultingPod.Spec.SchedulingGates)
				require.Equal(t, "true", resultingPod.Annotations[apicommon.SchedulingGateRemoved])
			}
		})
	}
}

func TestHasKeptnSchedulingGate(t *testing.T) {
	tests := []struct {
		name    string
		pod     *v1.Pod
		hasGate bool
	}{
		{
			name: "PodWithKeptnSchedulingGate",
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					SchedulingGates: []v1.PodSchedulingGate{
						{
							Name: apicommon.KeptnGate,
						},
					},
				},
			},
			hasGate: true,
		},
		{
			name:    "PodWithoutSchedulingGate",
			pod:     &v1.Pod{},
			hasGate: false,
		},
		{
			name: "PodWithOtherSchedulingGates",
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					SchedulingGates: []v1.PodSchedulingGate{
						{
							Name: "other-gate",
						},
					},
				},
			},
			hasGate: false,
		},
		{
			name: "PodWithKeptnAndOtherSchedulingGates",
			pod: &v1.Pod{
				Spec: v1.PodSpec{
					SchedulingGates: []v1.PodSchedulingGate{
						{
							Name: apicommon.KeptnGate,
						},
						{
							Name: "other-gate",
						},
					},
				},
			},
			hasGate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hasGate := hasKeptnSchedulingGate(tt.pod)
			require.Equal(t, tt.hasGate, hasGate)
		})
	}
}
