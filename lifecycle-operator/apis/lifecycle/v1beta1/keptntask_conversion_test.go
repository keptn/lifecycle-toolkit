//nolint:dupl
package v1beta1

import (
	"testing"
	"time"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	v1alpha3common "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1/common"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v2 "sigs.k8s.io/controller-runtime/pkg/webhook/conversion/testdata/api/v2"
)

func TestKeptnTask_ConvertFrom(t *testing.T) {
	tests := []struct {
		name    string
		srcObj  *v1alpha3.KeptnTask
		wantErr bool
		wantObj *KeptnTask
	}{
		{
			name: "Test that conversion from v1beta1 to v1alpha3 works",
			srcObj: &v1alpha3.KeptnTask{
				TypeMeta: v1.TypeMeta{
					Kind:       "KeptnTask",
					APIVersion: "lifecycle.keptn.sh/v1beta1",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      "some-keptn-app-name",
					Namespace: "",
					Labels: map[string]string{
						"some-label": "some-label-value",
					},
					Annotations: map[string]string{
						"some-annotation": "some-annotation-value",
					},
				},
				Spec: v1alpha3.KeptnTaskSpec{
					TaskDefinition: "sample-TaskDefinition",
					Context: v1alpha3.TaskContext{
						WorkloadName:    "sample-Workload",
						AppName:         "sample-App",
						AppVersion:      "1.2.3",
						WorkloadVersion: "4.5.6",
						TaskType:        "sample-task",
						ObjectType:      "sample-type",
					},
					Parameters: v1alpha3.TaskParameters{
						Inline: map[string]string{
							"key1": "sample-inline",
							"key2": "sample-inline2",
						},
					},
					SecureParameters: v1alpha3.SecureParameters{
						Secret: "new-secret",
					},
					Type: v1alpha3common.PostDeploymentCheckType,
					Timeout: v1.Duration{
						Duration: time.Duration(5 * time.Minute),
					},
				},
				Status: v1alpha3.KeptnTaskStatus{
					JobName:   "sample-Taskdefinition",
					Status:    v1alpha3common.StateFailed,
					Message:   "sample-message",
					StartTime: v1.Time{},
					EndTime:   v1.Time{},
					Reason:    "reason",
				},
			},
			wantErr: false,
			wantObj: &KeptnTask{
				ObjectMeta: v1.ObjectMeta{
					Name:      "some-keptn-app-name",
					Namespace: "",
					Labels: map[string]string{
						"some-label": "some-label-value",
					},
					Annotations: map[string]string{
						"some-annotation": "some-annotation-value",
					},
				},
				Spec: KeptnTaskSpec{
					TaskDefinition: "sample-TaskDefinition",
					Context: TaskContext{
						WorkloadName:    "sample-Workload",
						AppName:         "sample-App",
						AppVersion:      "1.2.3",
						WorkloadVersion: "4.5.6",
						TaskType:        "sample-task",
						ObjectType:      "sample-type",
					},
					Parameters: TaskParameters{
						Inline: map[string]string{
							"key1": "sample-inline",
							"key2": "sample-inline2",
						},
					},
					SecureParameters: SecureParameters{
						Secret: "new-secret",
					},
					Type: common.PostDeploymentCheckType,
					Timeout: v1.Duration{
						Duration: time.Duration(5 * time.Minute),
					},
				},
				Status: KeptnTaskStatus{
					JobName:   "sample-Taskdefinition",
					Status:    common.StateFailed,
					Message:   "sample-message",
					StartTime: v1.Time{},
					EndTime:   v1.Time{},
					Reason:    "reason",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := &KeptnTask{
				TypeMeta:   v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{},
				Spec:       KeptnTaskSpec{},
				Status:     KeptnTaskStatus{},
			}
			if err := dst.ConvertFrom(tt.srcObj); (err != nil) != tt.wantErr {
				t.Errorf("ConvertFrom() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantObj != nil {
				require.Equal(t, tt.wantObj, dst, "Object was not converted correctly")
			}
		})
	}
}

func TestKeptnTask_ConvertTo(t *testing.T) {
	tests := []struct {
		name    string
		src     *KeptnTask
		wantErr bool
		wantObj *v1alpha3.KeptnTask
	}{
		{
			name: "Test that conversion from v1beta1 to v1alpha3 works",
			src: &KeptnTask{
				TypeMeta: v1.TypeMeta{
					Kind:       "KeptnTask",
					APIVersion: "lifecycle.keptn.sh/v1alpha3",
				},
				ObjectMeta: v1.ObjectMeta{
					Name:      "some-keptn-app-name",
					Namespace: "",
					Labels: map[string]string{
						"some-label": "some-label-value",
					},
					Annotations: map[string]string{
						"some-annotation": "some-annotation-value",
					},
				},
				Spec: KeptnTaskSpec{
					TaskDefinition: "sample-TaskDefinition",
					Context: TaskContext{
						WorkloadName:    "sample-Workload",
						AppName:         "sample-App",
						AppVersion:      "1.2.3",
						WorkloadVersion: "4.5.6",
						TaskType:        "sample-task",
						ObjectType:      "sample-type",
					},
					Parameters: TaskParameters{
						Inline: map[string]string{
							"key1": "sample-inline",
							"key2": "sample-inline2",
						},
					},
					SecureParameters: SecureParameters{
						Secret: "new-secret",
					},
					Type: common.PostDeploymentCheckType,
					Timeout: v1.Duration{
						Duration: time.Duration(5 * time.Minute),
					},
				},
				Status: KeptnTaskStatus{
					JobName:   "sample-Taskdefinition",
					Status:    common.StateSucceeded,
					Message:   "sample-message",
					StartTime: v1.Time{},
					EndTime:   v1.Time{},
					Reason:    "reason",
				},
			},
			wantErr: false,
			wantObj: &v1alpha3.KeptnTask{
				ObjectMeta: v1.ObjectMeta{
					Name:      "some-keptn-app-name",
					Namespace: "",
					Labels: map[string]string{
						"some-label": "some-label-value",
					},
					Annotations: map[string]string{
						"some-annotation": "some-annotation-value",
					},
				},
				Spec: v1alpha3.KeptnTaskSpec{
					TaskDefinition: "sample-TaskDefinition",
					Context: v1alpha3.TaskContext{
						WorkloadName:    "sample-Workload",
						AppName:         "sample-App",
						AppVersion:      "1.2.3",
						WorkloadVersion: "4.5.6",
						TaskType:        "sample-task",
						ObjectType:      "sample-type",
					},
					Parameters: v1alpha3.TaskParameters{
						Inline: map[string]string{
							"key1": "sample-inline",
							"key2": "sample-inline2",
						},
					},
					SecureParameters: v1alpha3.SecureParameters{
						Secret: "new-secret",
					},
					Type: v1alpha3common.PostDeploymentCheckType,
					Timeout: v1.Duration{
						Duration: time.Duration(5 * time.Minute),
					},
				},
				Status: v1alpha3.KeptnTaskStatus{
					JobName:   "sample-Taskdefinition",
					Status:    v1alpha3common.StateSucceeded,
					Message:   "sample-message",
					StartTime: v1.Time{},
					EndTime:   v1.Time{},
					Reason:    "reason",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dst := v1alpha3.KeptnTask{
				TypeMeta:   v1.TypeMeta{},
				ObjectMeta: v1.ObjectMeta{},
				Spec:       v1alpha3.KeptnTaskSpec{},
				Status:     v1alpha3.KeptnTaskStatus{},
			}
			if err := tt.src.ConvertTo(&dst); (err != nil) != tt.wantErr {
				t.Errorf("ConvertTo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantObj != nil {
				require.Equal(t, tt.wantObj, &dst, "Object was not converted correctly")
			}
		})
	}
}

func TestKeptnTask_ConvertFrom_Errorcase(t *testing.T) {
	// A random different object is used here to simulate a different API version
	testObj := v2.ExternalJob{}

	dst := &KeptnTask{
		TypeMeta:   v1.TypeMeta{},
		ObjectMeta: v1.ObjectMeta{},
		Spec:       KeptnTaskSpec{},
		Status:     KeptnTaskStatus{},
	}

	if err := dst.ConvertFrom(&testObj); err == nil {
		t.Errorf("ConvertFrom() error = %v", err)
	} else {
		require.ErrorIs(t, err, common.ErrCannotCastKeptnAppVersion)
	}
}

func TestKeptnTask_ConvertTo_Errorcase(t *testing.T) {
	testObj := KeptnTask{}

	// A random different object is used here to simulate a different API version
	dst := v2.ExternalJob{}

	if err := testObj.ConvertTo(&dst); err == nil {
		t.Errorf("ConvertTo() error = %v", err)
	} else {
		require.ErrorIs(t, err, common.ErrCannotCastKeptnAppVersion)
	}
}
