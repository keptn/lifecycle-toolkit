package handlers

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	fakeclient "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	admissionv1 "k8s.io/api/admission/v1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func Test_copyAnnotationsIfParentAnnotated(t *testing.T) {
	testNamespace := "test-namespace"
	rsUidWithDpOwner := types.UID("this-is-the-replicaset-with-dp-owner")
	rsUidWithNoOwner := types.UID("this-is-the-replicaset-with-no-owner")
	testStsUid := types.UID("this-is-the-stateful-set-uid")
	tstStsName := "test-stateful-set"
	testDsUid := types.UID("this-is-the-daemon-set-uid")
	testDsName := "test-daemon-set"

	rsWithDpOwner := &appsv1.ReplicaSet{
		TypeMeta: metav1.TypeMeta{
			Kind: "ReplicaSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-replicaset1",
			UID:       rsUidWithDpOwner,
			Namespace: testNamespace,
			OwnerReferences: []metav1.OwnerReference{
				{
					Kind: "Deployment",
					Name: "this-is-the-deployment",
					UID:  "this-is-the-deployment-uid",
				},
			},
		},
	}
	// TODO: fix tests where an RS has a STS or DS as owner. they should not have a RS in between
	rsWithNoOwner := &appsv1.ReplicaSet{
		TypeMeta: metav1.TypeMeta{
			Kind: "ReplicaSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-replicaset4",
			UID:       rsUidWithNoOwner,
			Namespace: testNamespace,
		},
	}
	testDp := &appsv1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind: "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-deployment",
			UID:       "this-is-the-deployment-uid",
			Namespace: testNamespace,
		},
	}
	testSts := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind: "StatefulSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      tstStsName,
			UID:       testStsUid,
			Namespace: testNamespace,
		},
	}
	testDs := &appsv1.DaemonSet{
		TypeMeta: metav1.TypeMeta{
			Kind: "DaemonSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      testDsName,
			UID:       testDsUid,
			Namespace: testNamespace,
		},
	}

	fakeClient := fakeclient.NewClient(rsWithDpOwner, rsWithNoOwner, testDp, testSts, testDs)

	type fields struct {
		Client client.Client
		Log    logr.Logger
	}
	type args struct {
		ctx context.Context
		req *admission.Request
		pod *corev1.Pod
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "Test that nothing happens if owner UID is pod UID",
			fields: fields{
				Log:    testr.New(t),
				Client: fakeClient,
			},
			args: args{
				ctx: context.TODO(),
				req: &admission.Request{
					AdmissionRequest: admissionv1.AdmissionRequest{
						Namespace: testNamespace,
					},
				},
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						UID: "some-uid",
					},
				},
			},
			want: false,
		},
		{
			name: "Test fetching of replicaset owner of pod and deployment owner of replicaset",
			fields: fields{
				Log:    testr.New(t),
				Client: fakeClient,
			},
			args: args{
				ctx: context.TODO(),
				req: &admission.Request{
					AdmissionRequest: admissionv1.AdmissionRequest{
						Namespace: testNamespace,
					},
				},
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						UID: "this-is-the-pod-uid",
						OwnerReferences: []metav1.OwnerReference{
							{
								Name: rsWithDpOwner.Name,
								UID:  rsUidWithDpOwner,
								Kind: "ReplicaSet",
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "Test fetching of statefulset owner of pod",
			fields: fields{
				Log:    testr.New(t),
				Client: fakeClient,
			},
			args: args{
				ctx: context.TODO(),
				req: &admission.Request{
					AdmissionRequest: admissionv1.AdmissionRequest{
						Namespace: testNamespace,
					},
				},
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						UID: "this-is-the-pod-uid",
						OwnerReferences: []metav1.OwnerReference{
							{
								Name: testSts.Name,
								UID:  testSts.UID,
								Kind: testSts.Kind,
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "Test fetching of daemonset owner of pod",
			fields: fields{
				Log:    testr.New(t),
				Client: fakeClient,
			},
			args: args{
				ctx: context.TODO(),
				req: &admission.Request{
					AdmissionRequest: admissionv1.AdmissionRequest{
						Namespace: testNamespace,
					},
				},
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						UID: "this-is-the-pod-uid",
						OwnerReferences: []metav1.OwnerReference{
							{
								Name: testDs.Name,
								UID:  testDs.UID,
								Kind: testDs.Kind,
							},
						},
					},
				},
			},
			want: false,
		},
		{
			name: "Test that method returns without doing anything when we get a pod with replicaset without owner",
			fields: fields{
				Log:    testr.New(t),
				Client: fakeClient,
			},
			args: args{
				ctx: context.TODO(),
				req: &admission.Request{
					AdmissionRequest: admissionv1.AdmissionRequest{
						Namespace: testNamespace,
					},
				},
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						UID: "this-is-the-pod-uid",
						OwnerReferences: []metav1.OwnerReference{
							{
								Name: rsWithNoOwner.Name,
								UID:  rsUidWithNoOwner,
								Kind: "ReplicaSet",
							},
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &PodAnnotationHandler{
				Client: tt.fields.Client,
				Log:    tt.fields.Log,
			}
			got := a.copyAnnotationsIfParentAnnotated(tt.args.ctx, tt.args.req, tt.args.pod)
			if got != tt.want {
				t.Errorf("copyAnnotationsIfParentAnnotated() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_copyResourceLabelsIfPresent(t *testing.T) {

	type args struct {
		sourceResource *metav1.ObjectMeta
		targetPod      *corev1.Pod
	}
	tests := []struct {
		name      string
		args      args
		want      bool
		wantedPod *corev1.Pod
	}{
		{
			name: "Test that annotations get copied from source to target",
			args: args{
				sourceResource: &metav1.ObjectMeta{
					Name: "testSourceObject",
					Annotations: map[string]string{
						apicommon.WorkloadAnnotation:                 "some-workload-name",
						apicommon.AppAnnotation:                      "some-app-name",
						apicommon.VersionAnnotation:                  "v1.0.0",
						apicommon.PreDeploymentTaskAnnotation:        "some-pre-deployment-task",
						apicommon.PostDeploymentTaskAnnotation:       "some-post-deployment-task",
						apicommon.PreDeploymentEvaluationAnnotation:  "some-pre-deployment-evaluation",
						apicommon.PostDeploymentEvaluationAnnotation: "some-post-deployment-evaluation",
					},
				},
				targetPod: &corev1.Pod{
					TypeMeta:   metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{},
					Spec:       corev1.PodSpec{},
					Status:     corev1.PodStatus{},
				},
			},
			want: true,
			wantedPod: &corev1.Pod{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						apicommon.WorkloadAnnotation:                 "some-workload-name",
						apicommon.AppAnnotation:                      "some-app-name",
						apicommon.VersionAnnotation:                  "v1.0.0",
						apicommon.PreDeploymentTaskAnnotation:        "some-pre-deployment-task",
						apicommon.PostDeploymentTaskAnnotation:       "some-post-deployment-task",
						apicommon.PreDeploymentEvaluationAnnotation:  "some-pre-deployment-evaluation",
						apicommon.PostDeploymentEvaluationAnnotation: "some-post-deployment-evaluation",
					},
				},
			},
		},
		{
			name: "Test that source labels get copied to target annotations",
			args: args{
				sourceResource: &metav1.ObjectMeta{
					Name: "testSourceObject",
					Labels: map[string]string{
						apicommon.WorkloadAnnotation:                 "some-workload-name",
						apicommon.AppAnnotation:                      "some-app-name",
						apicommon.VersionAnnotation:                  "v1.0.0",
						apicommon.PreDeploymentTaskAnnotation:        "some-pre-deployment-task",
						apicommon.PostDeploymentTaskAnnotation:       "some-post-deployment-task",
						apicommon.PreDeploymentEvaluationAnnotation:  "some-pre-deployment-evaluation",
						apicommon.PostDeploymentEvaluationAnnotation: "some-post-deployment-evaluation",
					},
				},
				targetPod: &corev1.Pod{},
			},
			want: true,
			wantedPod: &corev1.Pod{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						apicommon.WorkloadAnnotation:                 "some-workload-name",
						apicommon.AppAnnotation:                      "some-app-name",
						apicommon.VersionAnnotation:                  "v1.0.0",
						apicommon.PreDeploymentTaskAnnotation:        "some-pre-deployment-task",
						apicommon.PostDeploymentTaskAnnotation:       "some-post-deployment-task",
						apicommon.PreDeploymentEvaluationAnnotation:  "some-pre-deployment-evaluation",
						apicommon.PostDeploymentEvaluationAnnotation: "some-post-deployment-evaluation",
					},
				},
			},
		},
		{
			name: "Test that version label is generated correctly and rest is copied",
			args: args{
				sourceResource: &metav1.ObjectMeta{
					Name: "testSourceObject",
					Labels: map[string]string{
						apicommon.WorkloadAnnotation:                 "some-workload-name",
						apicommon.AppAnnotation:                      "some-app-name",
						apicommon.PreDeploymentTaskAnnotation:        "some-pre-deployment-task",
						apicommon.PostDeploymentTaskAnnotation:       "some-post-deployment-task",
						apicommon.PreDeploymentEvaluationAnnotation:  "some-pre-deployment-evaluation",
						apicommon.PostDeploymentEvaluationAnnotation: "some-post-deployment-evaluation",
					},
				},
				targetPod: &corev1.Pod{
					TypeMeta:   metav1.TypeMeta{},
					ObjectMeta: metav1.ObjectMeta{},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Image: "some-image:v1.0.0",
							},
						},
					},
					Status: corev1.PodStatus{},
				},
			},
			want: true,
			wantedPod: &corev1.Pod{
				TypeMeta: metav1.TypeMeta{},
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						apicommon.WorkloadAnnotation:                 "some-workload-name",
						apicommon.AppAnnotation:                      "some-app-name",
						apicommon.VersionAnnotation:                  "v1.0.0",
						apicommon.PreDeploymentTaskAnnotation:        "some-pre-deployment-task",
						apicommon.PostDeploymentTaskAnnotation:       "some-post-deployment-task",
						apicommon.PreDeploymentEvaluationAnnotation:  "some-pre-deployment-evaluation",
						apicommon.PostDeploymentEvaluationAnnotation: "some-post-deployment-evaluation",
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "some-image:v1.0.0",
						},
					},
				},
				Status: corev1.PodStatus{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := copyResourceLabelsIfPresent(tt.args.sourceResource, tt.args.targetPod)
			if got != tt.want {
				t.Errorf("copyResourceLabelsIfPresent() got = %v, want %v", got, tt.want)
			}
			if tt.wantedPod != nil {
				require.Equal(t, tt.wantedPod, tt.args.targetPod)
			}
		})
	}
}

func Test_IsPodAnnotated(t *testing.T) {
	type args struct {
		pod *corev1.Pod
	}
	tests := []struct {
		name      string
		args      args
		want      bool
		wantedPod *corev1.Pod
	}{
		{
			name: "Test return true when pod has workload annotation",
			args: args{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							apicommon.WorkloadAnnotation: "some-workload-name",
						},
					},
				},
			},
			want: true,
		},
		{
			name: "Test return true and initialize annotations when labels are set",
			args: args{
				pod: &corev1.Pod{
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{
							{
								Image: "some-image:v1",
							},
						},
					},
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							apicommon.WorkloadAnnotation: "some-workload-name",
						},
					},
				},
			},
			want: true,
			wantedPod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "some-image:v1",
						},
					},
				},
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						apicommon.WorkloadAnnotation: "some-workload-name",
					},
					Annotations: map[string]string{
						apicommon.VersionAnnotation: "v1",
					},
				},
			},
		},
		{
			name: "Test return false when annotations and labels are not set",
			args: args{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"some-other-label": "some-value",
						},
					},
				},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := isPodAnnotated(tt.args.pod)
			if got != tt.want {
				t.Errorf("isPodAnnotated() got = %v, want %v", got, tt.want)
			}
			if tt.wantedPod != nil {
				require.Equal(t, tt.wantedPod, tt.args.pod)
			}
		})
	}
}
