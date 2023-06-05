package v1alpha3

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestKeptnTaskDefinition_ValidateFields(t *testing.T) {
	tests := []struct {
		name    string
		spec    KeptnTaskDefinitionSpec
		want    *field.Error
		verb    string
		oldSpec runtime.Object
	}{
		{
			name: "with-no-function-or-container",
			spec: KeptnTaskDefinitionSpec{},
			want: field.Invalid(
				field.NewPath("spec"),
				KeptnTaskDefinitionSpec{},
				errors.New("Forbidden! Either Function or Container field must be defined").Error(),
			),
			verb: "create",
		},
		{
			name: "with-both-function-and-container",
			spec: KeptnTaskDefinitionSpec{
				Function:  &FunctionSpec{},
				Container: &ContainerSpec{},
			},
			want: field.Invalid(
				field.NewPath("spec"),
				KeptnTaskDefinitionSpec{
					Function:  &FunctionSpec{},
					Container: &ContainerSpec{},
				},
				errors.New("Forbidden! Both Function and Container fields cannot be defined simultaneously").Error(),
			),
			verb: "create",
		},
		{
			name: "with-function-only",
			spec: KeptnTaskDefinitionSpec{
				Function: &FunctionSpec{},
			},
			verb: "create",
		},
		{
			name: "with-container-only",
			spec: KeptnTaskDefinitionSpec{
				Container: &ContainerSpec{},
			},
			verb: "create",
		},
		{
			name: "update-with-both-function-and-container",
			spec: KeptnTaskDefinitionSpec{
				Function: &FunctionSpec{},
			},
			want: field.Invalid(
				field.NewPath("spec"),
				KeptnTaskDefinitionSpec{Function: &FunctionSpec{}},
				errors.New("Forbidden! Both Function and Container fields cannot be defined simultaneously").Error()),
			oldSpec: KeptnTaskDefinitionSpec{
				Function: &FunctionSpec{},
			},
			verb: "update",
		},
		{
			name: "delete",
			verb: "delete",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ktd := &KeptnTaskDefinition{
				ObjectMeta: metav1.ObjectMeta{Name: tt.name},
				Spec:       tt.spec,
			}

			var got error
			switch tt.verb {
			case "create":
				got = ktd.ValidateCreate()
			case "update":
				got = ktd.ValidateUpdate(tt.oldSpec)
			case "delete":
				got = ktd.ValidateDelete()
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
