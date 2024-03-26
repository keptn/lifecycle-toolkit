package keptntask

import (
	"testing"

	"github.com/go-logr/logr/testr"
	klcv1beta1 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/eventsender"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/taskdefinition"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/testcommon"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
)

func TestJSBuilder_handleParent(t *testing.T) {

	def := &klcv1beta1.KeptnTaskDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mytaskdef",
			Namespace: "default",
		},
		Spec: klcv1beta1.KeptnTaskDefinitionSpec{
			Deno: &klcv1beta1.RuntimeSpec{
				FunctionReference: klcv1beta1.FunctionReference{
					Name: "mytaskdef",
				}}},
	}
	paramDef := &klcv1beta1.KeptnTaskDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mytd",
			Namespace: "default",
		},
		Spec: klcv1beta1.KeptnTaskDefinitionSpec{
			Deno: &klcv1beta1.RuntimeSpec{
				FunctionReference: klcv1beta1.FunctionReference{
					Name: "mytd"},
				Parameters: klcv1beta1.TaskParameters{
					Inline: map[string]string{Data: "mydata"},
				},
				SecureParameters: klcv1beta1.SecureParameters{
					Secret: "mysecret",
				},
			},
		},
	}

	tests := []struct {
		name    string
		options BuilderOptions
		params  RuntimeExecutionParams
		wantErr bool
		err     string
	}{
		{
			name: "no definition",
			options: BuilderOptions{
				Client:      testcommon.NewTestClient(),
				eventSender: eventsender.NewK8sSender(record.NewFakeRecorder(100)),
				req: ctrl.Request{
					NamespacedName: types.NamespacedName{Namespace: "default"},
				},
				Log:      testr.New(t),
				funcSpec: taskdefinition.GetRuntimeSpec(def),
				task:     makeTask("myt", "default", def.Name),
			},
			params:  RuntimeExecutionParams{},
			wantErr: true,
			err:     "not found",
		},
		{
			name: "definition exists, recursive",
			options: BuilderOptions{
				Client:      testcommon.NewTestClient(def),
				eventSender: eventsender.NewK8sSender(record.NewFakeRecorder(100)),
				req: ctrl.Request{
					NamespacedName: types.NamespacedName{Namespace: "default"},
				},
				Log:      testr.New(t),
				funcSpec: taskdefinition.GetRuntimeSpec(def),
				task:     makeTask("myt2", "default", def.Name),
			},
			params:  RuntimeExecutionParams{},
			wantErr: false,
		},
		{
			name: "definition exists, with parameters and secrets",
			options: BuilderOptions{
				Client:      testcommon.NewTestClient(paramDef, def),
				eventSender: eventsender.NewK8sSender(record.NewFakeRecorder(100)),
				req: ctrl.Request{
					NamespacedName: types.NamespacedName{Namespace: "default"},
				},
				Log:      testr.New(t),
				funcSpec: taskdefinition.GetRuntimeSpec(paramDef),
				task:     makeTask("myt3", "default", paramDef.Name),
			},
			params:  RuntimeExecutionParams{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			js := &RuntimeBuilder{
				options: tt.options,
			}
			err := js.handleParent(context.TODO(), &tt.params)
			if !tt.wantErr {
				require.Nil(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err)
			}

		})
	}
}
func TestJSBuilder_getParams(t *testing.T) {
	t.Setenv(taskdefinition.FunctionRuntimeImageKey, taskdefinition.FunctionScriptKey)
	t.Setenv(taskdefinition.PythonRuntimeImageKey, taskdefinition.PythonScriptKey)

	def := &klcv1beta1.KeptnTaskDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mytaskdef",
			Namespace: "default",
		},
		Spec: klcv1beta1.KeptnTaskDefinitionSpec{
			Deno: &klcv1beta1.RuntimeSpec{
				Parameters: klcv1beta1.TaskParameters{
					Inline: map[string]string{"DATA2": "parent_data"},
				},
				SecureParameters: klcv1beta1.SecureParameters{
					Secret: "parent_secret",
				},
			},
		},
		Status: klcv1beta1.KeptnTaskDefinitionStatus{
			Function: klcv1beta1.FunctionStatus{
				ConfigMap: "mymap",
			},
		},
	}
	paramDef := &klcv1beta1.KeptnTaskDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mytd",
			Namespace: "default",
		},
		Spec: klcv1beta1.KeptnTaskDefinitionSpec{
			Deno: &klcv1beta1.RuntimeSpec{
				FunctionReference: klcv1beta1.FunctionReference{
					Name: def.Name},
				Parameters: klcv1beta1.TaskParameters{
					Inline: map[string]string{"DATA1": "child_data"},
				},
				SecureParameters: klcv1beta1.SecureParameters{
					Secret: "child_pw",
				},
			},
		},
		Status: klcv1beta1.KeptnTaskDefinitionStatus{
			Function: klcv1beta1.FunctionStatus{
				ConfigMap: "mychildmap",
			},
		},
	}

	parentPy := &klcv1beta1.KeptnTaskDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "parentPy",
			Namespace: "default",
		},
		Spec: klcv1beta1.KeptnTaskDefinitionSpec{
			Python: &klcv1beta1.RuntimeSpec{
				HttpReference: klcv1beta1.HttpReference{Url: "donothing"},
			}},
	}
	defJS := &klcv1beta1.KeptnTaskDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "myJS",
			Namespace: "default",
		},
		Spec: klcv1beta1.KeptnTaskDefinitionSpec{
			Deno: &klcv1beta1.RuntimeSpec{
				FunctionReference: klcv1beta1.FunctionReference{
					Name: parentPy.Name},
			},
		},
		Status: klcv1beta1.KeptnTaskDefinitionStatus{
			Function: klcv1beta1.FunctionStatus{
				ConfigMap: "myJSChildmap",
			},
		},
	}

	tests := []struct {
		name    string
		options BuilderOptions
		params  *RuntimeExecutionParams
		wantErr bool
		err     string
	}{
		{
			name: "definition exists, no parent",
			options: BuilderOptions{
				Client:      testcommon.NewTestClient(def),
				eventSender: eventsender.NewK8sSender(record.NewFakeRecorder(100)),
				req: ctrl.Request{
					NamespacedName: types.NamespacedName{Namespace: "default"},
				},
				Log:       testr.New(t),
				funcSpec:  taskdefinition.GetRuntimeSpec(def),
				task:      makeTask("myt2", "default", def.Name),
				Image:     taskdefinition.FunctionScriptKey,
				MountPath: taskdefinition.FunctionScriptMountPath,
				ConfigMap: def.Status.Function.ConfigMap,
			},
			params: &RuntimeExecutionParams{
				ConfigMap:        def.Status.Function.ConfigMap,
				Parameters:       def.Spec.Deno.Parameters.Inline,
				SecureParameters: def.Spec.Deno.SecureParameters.Secret,
				URL:              def.Spec.Deno.HttpReference.Url,
				Context: klcv1beta1.TaskContext{
					WorkloadName: "my-workload",
					AppName:      "my-app",
					AppVersion:   "0.1.0",
					ObjectType:   "Workload",
					TaskType:     string(apicommon.PostDeploymentCheckType),
				},
				Image:     taskdefinition.FunctionScriptKey,
				MountPath: taskdefinition.FunctionScriptMountPath,
			},
			wantErr: false,
		},
		{
			name: "definition exists, parent with parameters and secrets",
			options: BuilderOptions{
				Client:      testcommon.NewTestClient(paramDef, def),
				eventSender: eventsender.NewK8sSender(record.NewFakeRecorder(100)),
				req: ctrl.Request{
					NamespacedName: types.NamespacedName{Namespace: "default"},
				},
				Log:       testr.New(t),
				funcSpec:  taskdefinition.GetRuntimeSpec(paramDef),
				task:      makeTask("myt3", "default", paramDef.Name),
				ConfigMap: def.Status.Function.ConfigMap,
			},
			params: &RuntimeExecutionParams{
				ConfigMap: def.Status.Function.ConfigMap,
				Parameters: map[string]string{ // maps should be merged
					"DATA2": "parent_data",
					"DATA1": "child_data",
				},
				SecureParameters: paramDef.Spec.Deno.SecureParameters.Secret, // uses child
				URL:              def.Spec.Deno.HttpReference.Url,            // uses parent
				Context: klcv1beta1.TaskContext{
					WorkloadName: "my-workload",
					AppName:      "my-app",
					AppVersion:   "0.1.0",
					ObjectType:   "Workload",
					TaskType:     string(apicommon.PostDeploymentCheckType),
				},
				Image:     taskdefinition.FunctionScriptKey,
				MountPath: taskdefinition.FunctionScriptMountPath,
			},
			wantErr: false,
		},
		{
			name: "definition exists, parent is of a different runtime",
			options: BuilderOptions{
				Client:      testcommon.NewTestClient(parentPy, defJS),
				eventSender: eventsender.NewK8sSender(record.NewFakeRecorder(100)),
				req: ctrl.Request{
					NamespacedName: types.NamespacedName{Namespace: "default"},
				},
				Log:       testr.New(t),
				funcSpec:  taskdefinition.GetRuntimeSpec(defJS),
				task:      makeTask("myt4", "default", defJS.Name),
				ConfigMap: defJS.Status.Function.ConfigMap,
			},
			params: &RuntimeExecutionParams{
				ConfigMap: parentPy.Status.Function.ConfigMap,
				URL:       parentPy.Spec.Python.HttpReference.Url, // we support a single URL so the original should be taken not the parent one
				Context: klcv1beta1.TaskContext{
					WorkloadName: "my-workload",
					AppName:      "my-app",
					AppVersion:   "0.1.0",
					ObjectType:   "Workload",
					TaskType:     string(apicommon.PostDeploymentCheckType),
				},
				Image:     taskdefinition.PythonScriptKey,
				MountPath: taskdefinition.PythonScriptMountPath,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			js := &RuntimeBuilder{
				options: tt.options,
			}
			params, err := js.getParams(context.TODO())
			if !tt.wantErr {
				require.Nil(t, err)
				require.Equal(t, tt.params, params)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err)
			}

		})
	}
}
