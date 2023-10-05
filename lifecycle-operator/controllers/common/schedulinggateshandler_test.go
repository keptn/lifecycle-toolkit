package common

import (
	"context"
	"fmt"
	"testing"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	controllererrors "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/errors"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func Test_RemovePodGates(t *testing.T) {
	tests := []struct {
		name        string
		podName     string
		pod         *corev1.Pod
		wantError   bool
		annotations map[string]string
	}{
		{
			name:    "pod does not exist",
			podName: "pod",
			pod: &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Name:      "pod2",
					Namespace: "default",
				},
			},
			wantError: true,
		},
		{
			name:    "scheduling gates already removed",
			podName: "pod",
			pod: &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Name:      "pod",
					Namespace: "default",
					Annotations: map[string]string{
						apicommon.SchedulingGateRemoved: "true",
					},
				},
			},
			wantError: false,
			annotations: map[string]string{
				apicommon.SchedulingGateRemoved: "true",
			},
		},
		{
			name:    "scheduling gates removed - empty annotations",
			podName: "pod",
			pod: &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Name:      "pod",
					Namespace: "default",
				},
				Spec: corev1.PodSpec{
					SchedulingGates: []corev1.PodSchedulingGate{
						{
							Name: "gate",
						},
					},
				},
			},
			wantError: false,
			annotations: map[string]string{
				apicommon.SchedulingGateRemoved: "true",
			},
		},
		{
			name:    "scheduling gates removed - not empty annotations",
			podName: "pod",
			pod: &corev1.Pod{
				ObjectMeta: v1.ObjectMeta{
					Name:      "pod",
					Namespace: "default",
					Annotations: map[string]string{
						"test": "test",
					},
				},
				Spec: corev1.PodSpec{
					SchedulingGates: []corev1.PodSchedulingGate{
						{
							Name: "gate",
						},
					},
				},
			},
			wantError: false,
			annotations: map[string]string{
				apicommon.SchedulingGateRemoved: "true",
				"test":                          "test",
			},
		},
	}

	err := klcv1alpha3.AddToScheme(scheme.Scheme)
	require.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := fake.NewClientBuilder().WithObjects(tt.pod).WithStatusSubresource(tt.pod).Build()
			err := removePodGates(context.TODO(), client, tt.podName, tt.pod.Namespace)
			if tt.wantError != (err != nil) {
				t.Errorf("want error: %t, got: %v", tt.wantError, err)
			}
			if !tt.wantError {
				pod := &corev1.Pod{}
				err := client.Get(context.TODO(), types.NamespacedName{Namespace: tt.pod.Namespace, Name: tt.podName}, pod)
				require.Nil(t, err)
				require.Equal(t, []corev1.PodSchedulingGate(nil), pod.Spec.SchedulingGates)
				require.Equal(t, tt.annotations, pod.Annotations)
			}
		})
	}
}

func Test_GetPodsOfOwner(t *testing.T) {
	namespace := "default"
	tests := []struct {
		name   string
		uid    types.UID
		kind   string
		pods   *corev1.PodList
		result []string
	}{
		{
			name:   "pod list empty",
			pods:   &corev1.PodList{},
			result: nil,
		},
		{
			name: "pod list not matching kind or uid",
			pods: &corev1.PodList{
				Items: []corev1.Pod{
					{
						ObjectMeta: v1.ObjectMeta{
							Name:      "pod1",
							Namespace: "default",
							OwnerReferences: []v1.OwnerReference{
								{
									Kind: "unknown",
									UID:  types.UID("uid"),
								},
							},
						},
					},
				},
			},
			kind:   "unknown2",
			uid:    types.UID("uid2"),
			result: nil,
		},
		{
			name: "pod list matches one pod of list",
			pods: &corev1.PodList{
				Items: []corev1.Pod{
					{
						ObjectMeta: v1.ObjectMeta{
							Name:      "pod1",
							Namespace: "default",
							OwnerReferences: []v1.OwnerReference{
								{
									Kind: "unknown",
									UID:  types.UID("uid"),
								},
							},
						},
					},
					{
						ObjectMeta: v1.ObjectMeta{
							Name:      "pod2",
							Namespace: "default",
							OwnerReferences: []v1.OwnerReference{
								{
									Kind: "unknown2",
									UID:  types.UID("uid2"),
								},
							},
						},
					},
				},
			},
			kind:   "unknown",
			uid:    types.UID("uid"),
			result: []string{"pod1"},
		},
	}

	err := klcv1alpha3.AddToScheme(scheme.Scheme)
	require.Nil(t, err)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := fake.NewClientBuilder().WithLists(tt.pods).Build()
			res, err := getPodsOfOwner(context.TODO(), client, tt.uid, tt.kind, namespace)
			require.Nil(t, err)
			require.Equal(t, tt.result, res)
		})
	}
}

func Test_SchedulingGatesHandler_IsSchedulingGatesEnabled(t *testing.T) {
	h := SchedulingGatesHandler{
		enabled: true,
	}

	require.True(t, h.Enabled())

	h.enabled = false

	require.False(t, h.Enabled())
}

func Test_SchedulingGatesHandler_IsSchedulingGatesEnabledRemoveGates(t *testing.T) {
	tests := []struct {
		name    string
		handler SchedulingGatesHandler
		wi      *klcv1alpha3.KeptnWorkloadVersion
		wantErr error
	}{
		{
			name:    "unsuported resource ref",
			handler: SchedulingGatesHandler{},
			wi: &klcv1alpha3.KeptnWorkloadVersion{
				Spec: klcv1alpha3.KeptnWorkloadVersionSpec{
					KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
						ResourceReference: klcv1alpha3.ResourceReference{
							Kind: "unsupported",
						},
					},
				},
			},
			wantErr: controllererrors.ErrUnsupportedWorkloadVersionResourceReference,
		},
		{
			name: "pod - happy path",
			handler: SchedulingGatesHandler{
				removeGates: func(ctx context.Context, c client.Client, podName, podNamespace string) error {
					return nil
				},
			},
			wi: &klcv1alpha3.KeptnWorkloadVersion{
				Spec: klcv1alpha3.KeptnWorkloadVersionSpec{
					KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
						ResourceReference: klcv1alpha3.ResourceReference{
							Kind: "Pod",
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "pod - fail path",
			handler: SchedulingGatesHandler{
				removeGates: func(ctx context.Context, c client.Client, podName, podNamespace string) error {
					return fmt.Errorf("pod")
				},
			},
			wi: &klcv1alpha3.KeptnWorkloadVersion{
				Spec: klcv1alpha3.KeptnWorkloadVersionSpec{
					KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
						ResourceReference: klcv1alpha3.ResourceReference{
							Kind: "Pod",
						},
					},
				},
			},
			wantErr: fmt.Errorf("pod"),
		},
		{
			name: "ReplicaSet, StatefulSet, DaemonSet - happy path",
			handler: SchedulingGatesHandler{
				removeGates: func(ctx context.Context, c client.Client, podName, podNamespace string) error {
					return nil
				},
				getPods: func(ctx context.Context, c client.Client, ownerUID types.UID, ownerKind, namespace string) ([]string, error) {
					return []string{"podName"}, nil
				},
			},
			wi: &klcv1alpha3.KeptnWorkloadVersion{
				Spec: klcv1alpha3.KeptnWorkloadVersionSpec{
					KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
						ResourceReference: klcv1alpha3.ResourceReference{
							Kind: "ReplicaSet",
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "ReplicaSet, StatefulSet, DaemonSet - happy path - no pods found",
			handler: SchedulingGatesHandler{
				getPods: func(ctx context.Context, c client.Client, ownerUID types.UID, ownerKind, namespace string) ([]string, error) {
					return []string{}, nil
				},
			},
			wi: &klcv1alpha3.KeptnWorkloadVersion{
				Spec: klcv1alpha3.KeptnWorkloadVersionSpec{
					KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
						ResourceReference: klcv1alpha3.ResourceReference{
							Kind: "ReplicaSet",
						},
					},
				},
			},
			wantErr: nil,
		},
		{
			name: "ReplicaSet, StatefulSet, DaemonSet - happy path - err getPods",
			handler: SchedulingGatesHandler{
				getPods: func(ctx context.Context, c client.Client, ownerUID types.UID, ownerKind, namespace string) ([]string, error) {
					return []string{}, fmt.Errorf("err")
				},
			},
			wi: &klcv1alpha3.KeptnWorkloadVersion{
				Spec: klcv1alpha3.KeptnWorkloadVersionSpec{
					KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
						ResourceReference: klcv1alpha3.ResourceReference{
							Kind: "ReplicaSet",
						},
					},
				},
			},
			wantErr: fmt.Errorf("err"),
		},
		{
			name: "ReplicaSet, StatefulSet, DaemonSet - err removeGates",
			handler: SchedulingGatesHandler{
				removeGates: func(ctx context.Context, c client.Client, podName, podNamespace string) error {
					return fmt.Errorf("err")
				},
				getPods: func(ctx context.Context, c client.Client, ownerUID types.UID, ownerKind, namespace string) ([]string, error) {
					return []string{"podName"}, nil
				},
			},
			wi: &klcv1alpha3.KeptnWorkloadVersion{
				Spec: klcv1alpha3.KeptnWorkloadVersionSpec{
					KeptnWorkloadSpec: klcv1alpha3.KeptnWorkloadSpec{
						ResourceReference: klcv1alpha3.ResourceReference{
							Kind: "ReplicaSet",
						},
					},
				},
			},
			wantErr: fmt.Errorf("err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.handler.RemoveGates(context.TODO(), tt.wi)
			require.Equal(t, tt.wantErr, err)
		})
	}
}
