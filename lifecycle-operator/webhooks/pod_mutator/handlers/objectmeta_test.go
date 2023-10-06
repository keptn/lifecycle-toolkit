package handlers

import (
	"reflect"
	"testing"

	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const appname = "SOME-APP-NAME"

func Test_getWorkloadName(t *testing.T) {

	type args struct {
		pod *corev1.Pod
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Return concatenated app name and workload name in lower case when annotations are set",
			args: args{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							apicommon.AppAnnotation:      appname,
							apicommon.WorkloadAnnotation: "SOME-WORKLOAD-NAME",
						},
					},
				},
			},
			want: "some-app-name-some-workload-name",
		},
		{
			name: "Return concatenated app name and workload name in lower case when labels are set",
			args: args{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							apicommon.AppAnnotation:      appname,
							apicommon.WorkloadAnnotation: "SOME-WORKLOAD-NAME",
						},
					},
				},
			},
			want: "some-app-name-some-workload-name",
		},
		{
			name: "Return concatenated keptn app name and workload name from annotation in lower case when annotations and labels are set",
			args: args{
				pod: &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Annotations: map[string]string{
							apicommon.AppAnnotation:      "SOME-APP-NAME-ANNOTATION",
							apicommon.WorkloadAnnotation: "SOME-WORKLOAD-NAME-ANNOTATION",
						},
						Labels: map[string]string{
							apicommon.AppAnnotation:      "SOME-APP-NAME-LABEL",
							apicommon.WorkloadAnnotation: "SOME-WORKLOAD-NAME-LABEL",
						},
					},
				},
			},
			want: "some-app-name-annotation-some-workload-name-annotation",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := getWorkloadName(&tt.args.pod.ObjectMeta, getAppName(&tt.args.pod.ObjectMeta)); got != tt.want {
				t.Errorf("getWorkloadName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getLabelOrAnnotation(t *testing.T) {
	type args struct {
		resource            *metav1.ObjectMeta
		primaryAnnotation   string
		secondaryAnnotation string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name: "Test if primary annotation is returned from annotations",
			args: args{
				resource: &metav1.ObjectMeta{
					Annotations: map[string]string{
						apicommon.AppAnnotation: appname,
					},
				},
				primaryAnnotation:   apicommon.AppAnnotation,
				secondaryAnnotation: apicommon.K8sRecommendedAppAnnotations,
			},
			want:  appname,
			want1: true,
		},
		{
			name: "Test if secondary annotation is returned from annotations",
			args: args{
				resource: &metav1.ObjectMeta{
					Annotations: map[string]string{
						apicommon.K8sRecommendedAppAnnotations: appname,
					},
				},
				primaryAnnotation:   apicommon.AppAnnotation,
				secondaryAnnotation: apicommon.K8sRecommendedAppAnnotations,
			},
			want:  appname,
			want1: true,
		},
		{
			name: "Test if primary annotation is returned from labels",
			args: args{
				resource: &metav1.ObjectMeta{
					Labels: map[string]string{
						apicommon.AppAnnotation: appname,
					},
				},
				primaryAnnotation:   apicommon.AppAnnotation,
				secondaryAnnotation: apicommon.K8sRecommendedAppAnnotations,
			},
			want:  appname,
			want1: true,
		},
		{
			name: "Test if secondary annotation is returned from labels",
			args: args{
				resource: &metav1.ObjectMeta{
					Labels: map[string]string{
						apicommon.K8sRecommendedAppAnnotations: appname,
					},
				},
				primaryAnnotation:   apicommon.AppAnnotation,
				secondaryAnnotation: apicommon.K8sRecommendedAppAnnotations,
			},
			want:  appname,
			want1: true,
		},
		{
			name: "Test that empty string is returned when no annotations or labels are found",
			args: args{
				resource: &metav1.ObjectMeta{
					Annotations: map[string]string{
						"some-other-annotation": appname,
					},
				},
				primaryAnnotation:   apicommon.AppAnnotation,
				secondaryAnnotation: apicommon.K8sRecommendedAppAnnotations,
			},
			want:  "",
			want1: false,
		},
		{
			name: "Test that empty string is returned when primary annotation cannot be found and secondary annotation is empty",
			args: args{
				resource: &metav1.ObjectMeta{
					Annotations: map[string]string{
						"some-other-annotation": appname,
					},
				},
				primaryAnnotation:   apicommon.AppAnnotation,
				secondaryAnnotation: "",
			},
			want:  "",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetLabelOrAnnotation(tt.args.resource, tt.args.primaryAnnotation, tt.args.secondaryAnnotation)
			if got != tt.want {
				t.Errorf("getLabelOrAnnotation() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getLabelOrAnnotation() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_calculateVersion(t *testing.T) {

	tests := []struct {
		name string
		pod  *corev1.Pod
		want string
	}{
		{
			name: "simple tag",
			pod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Image: "ciao:1.0.0"},
					},
				}},
			want: "1.0.0",
		}, {
			name: "local registry",
			pod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{Image: "localhost:5000/node-web-app:1.0.0"},
					},
				}},
			want: "1.0.0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateVersion(tt.pod); got != tt.want {
				t.Errorf("calculateVersion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAppName(t *testing.T) {

	tests := []struct {
		name string
		pod  *corev1.Pod
		want string
	}{
		{
			name: "Return keptn app name in lower case when annotation is set",

			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						apicommon.AppAnnotation: appname,
					},
				},
			},

			want: appname,
		},
		{
			name: "Return keptn app name in lower case when label is set",

			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						apicommon.AppAnnotation: appname,
					},
				},
			},

			want: appname,
		},
		{
			name: "Return keptn app name from annotation in lower case when annotation and label is set",

			pod: &corev1.Pod{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						apicommon.AppAnnotation: "SOME-APP-NAME-ANNOTATION",
					},
					Labels: map[string]string{
						apicommon.AppAnnotation: "SOME-APP-NAME-LABEL",
					},
				},
			},

			want: "some-app-name-annotation",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAppName(&tt.pod.ObjectMeta); got != tt.want {
				t.Errorf("getAppName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getOwnerReference(t *testing.T) {

	ownerRef := metav1.OwnerReference{
		UID:  "the-replicaset-uid",
		Kind: "ReplicaSet",
		Name: "some-name",
	}

	type args struct {
		resource metav1.ObjectMeta
	}
	tests := []struct {
		name string
		args args
		want metav1.OwnerReference
	}{
		{
			name: "Test simple return when UID and Kind is set",
			args: args{
				resource: metav1.ObjectMeta{
					UID: "the-pod-uid",
					OwnerReferences: []metav1.OwnerReference{
						ownerRef,
					},
				},
			},
			want: ownerRef,
		},
		{
			name: "Test return is input argument if owner is not found",
			args: args{
				resource: metav1.ObjectMeta{
					UID: "the-pod-uid",
					OwnerReferences: []metav1.OwnerReference{
						{
							Kind: "SomeNonExistentType",
							UID:  "the-replicaset-uid",
						},
					},
				},
			},
			want: metav1.OwnerReference{
				UID:  "",
				Kind: "",
				Name: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetOwnerReference(&tt.args.resource); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getOwnerReference() = %v, want %v", got, tt.want)
			}
		})
	}
}
