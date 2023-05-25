package keptntask

import (
	"testing"

	"github.com/go-logr/logr/testr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
)

func TestJSBuilder_handleParent(t *testing.T) {

	def := &klcv1alpha3.KeptnTaskDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mytaskdef",
			Namespace: "default",
		},
		Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
			Function: klcv1alpha3.FunctionSpec{
				FunctionReference: klcv1alpha3.FunctionReference{
					Name: "mytaskdef",
				}}},
	}
	paramDef := &klcv1alpha3.KeptnTaskDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mytd",
			Namespace: "default",
		},
		Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
			Function: klcv1alpha3.FunctionSpec{
				FunctionReference: klcv1alpha3.FunctionReference{
					Name: "mytd"},
				Parameters: klcv1alpha3.TaskParameters{
					Inline: map[string]string{"DATA": "mydata"},
				},
				SecureParameters: klcv1alpha3.SecureParameters{
					Secret: "mysecret",
				},
			},
		},
	}

	tests := []struct {
		name    string
		options BuilderOptions
		params  FunctionExecutionParams
		wantErr bool
		err     string
	}{
		{
			name: "no definition",
			options: BuilderOptions{
				Client:   fake.NewClient(),
				recorder: &record.FakeRecorder{},
				req: ctrl.Request{
					NamespacedName: types.NamespacedName{Namespace: "default"},
				},
				Log:     testr.New(t),
				taskDef: def,
				task:    makeTask("myt", "default", def.Name),
			},
			params:  FunctionExecutionParams{},
			wantErr: true,
			err:     "not found",
		},
		{
			name: "definition exists, recursive",
			options: BuilderOptions{
				Client:   fake.NewClient(def),
				recorder: &record.FakeRecorder{},
				req: ctrl.Request{
					NamespacedName: types.NamespacedName{Namespace: "default"},
				},
				Log:     testr.New(t),
				taskDef: def,
				task:    makeTask("myt2", "default", def.Name),
			},
			params:  FunctionExecutionParams{},
			wantErr: false,
		},
		{
			name: "definition exists, with parameters and secrets",
			options: BuilderOptions{
				Client:   fake.NewClient(paramDef, def),
				recorder: &record.FakeRecorder{},
				req: ctrl.Request{
					NamespacedName: types.NamespacedName{Namespace: "default"},
				},
				Log:     testr.New(t),
				taskDef: paramDef,
				task:    makeTask("myt3", "default", paramDef.Name),
			},
			params:  FunctionExecutionParams{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			js := &JSBuilder{
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
func TestJSBuilder_hasParams(t *testing.T) {

	def := &klcv1alpha3.KeptnTaskDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mytaskdef",
			Namespace: "default",
		},
		Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
			Function: klcv1alpha3.FunctionSpec{
				HttpReference: klcv1alpha3.HttpReference{Url: "donothing"},
				Parameters: klcv1alpha3.TaskParameters{
					Inline: map[string]string{"DATA2": "mydata2"},
				},
				SecureParameters: klcv1alpha3.SecureParameters{
					Secret: "mysecret2",
				},
			}},
	}
	paramDef := &klcv1alpha3.KeptnTaskDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "mytd",
			Namespace: "default",
		},
		Spec: klcv1alpha3.KeptnTaskDefinitionSpec{
			Function: klcv1alpha3.FunctionSpec{
				HttpReference: klcv1alpha3.HttpReference{Url: "something"},
				FunctionReference: klcv1alpha3.FunctionReference{
					Name: "mytaskdef"},
				Parameters: klcv1alpha3.TaskParameters{
					Inline: map[string]string{"DATA1": "user"},
				},
				SecureParameters: klcv1alpha3.SecureParameters{
					Secret: "pw",
				},
			},
		},
	}

	tests := []struct {
		name    string
		options BuilderOptions
		params  *FunctionExecutionParams
		wantErr bool
		err     string
	}{
		{
			name: "definition exists, no parent",
			options: BuilderOptions{
				Client:   fake.NewClient(def),
				recorder: &record.FakeRecorder{},
				req: ctrl.Request{
					NamespacedName: types.NamespacedName{Namespace: "default"},
				},
				Log:     testr.New(t),
				taskDef: def,
				task:    makeTask("myt2", "default", def.Name),
			},
			params: &FunctionExecutionParams{
				ConfigMap: "",
				Parameters: map[string]string{
					"DATA2": "mydata2",
				},
				SecureParameters: "mysecret2",
				URL:              "donothing",
				Context: klcv1alpha3.TaskContext{
					WorkloadName: "my-workload",
					AppName:      "my-app",
					ObjectType:   "Workload"},
			},
			wantErr: false,
		},
		{
			name: "definition exists, parent with parameters and secrets",
			options: BuilderOptions{
				Client:   fake.NewClient(paramDef, def),
				recorder: &record.FakeRecorder{},
				req: ctrl.Request{
					NamespacedName: types.NamespacedName{Namespace: "default"},
				},
				Log:     testr.New(t),
				taskDef: paramDef,
				task:    makeTask("myt3", "default", paramDef.Name),
			},
			params: &FunctionExecutionParams{
				ConfigMap: "",
				Parameters: map[string]string{ //maps should be merged
					"DATA2": "mydata2",
					"DATA1": "user",
				},
				URL:              "something", //we support a single URL so the original should be taken not the parent one
				SecureParameters: "pw",        //we support a single secret so the original task secret should be taken not the parent one
				Context: klcv1alpha3.TaskContext{
					WorkloadName: "my-workload",
					AppName:      "my-app",
					ObjectType:   "Workload"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			js := &JSBuilder{
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
