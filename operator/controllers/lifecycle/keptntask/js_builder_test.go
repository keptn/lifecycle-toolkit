package keptntask

import (
	"github.com/go-logr/logr/testr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"testing"
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
	tests := []struct {
		name    string
		options BuilderOptions
		params  FunctionExecutionParams
		wantErr bool
		err     string
	}{{
		name: "no definition",
		options: BuilderOptions{
			Client:   fake.NewClient(),
			recorder: &record.FakeRecorder{},
			req: ctrl.Request{
				NamespacedName: types.NamespacedName{Namespace: "default"},
			},
			Log:     testr.New(t),
			taskDef: def,
			task:    makeTask("myt", "default", "mytaskdef"),
		},
		params:  FunctionExecutionParams{},
		wantErr: true,
		err:     "not found",
	}, {
		name: "definition exists",
		options: BuilderOptions{
			Client:   fake.NewClient(def),
			recorder: &record.FakeRecorder{},
			req: ctrl.Request{
				NamespacedName: types.NamespacedName{Namespace: "default"},
			},
			Log:     testr.New(t),
			taskDef: def,
			task:    makeTask("myt", "default", "mytaskdef"),
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
			err := js.handleParent(context.TODO(), tt.params)
			if !tt.wantErr {
				require.Nil(t, err)
			} else {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err)
			}

		})
	}
}
