package v1alpha3

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestKeptnTaskDefinition_ValidateFields(t *testing.T) {
	tests := []struct {
		name string
		spec KeptnTaskDefinitionSpec
		want *field.Error
	}{
		{
			name: "with-no-function-or-container",
			spec: KeptnTaskDefinitionSpec{},
			want: field.Invalid(
				field.NewPath("spec"),
				KeptnTaskDefinitionSpec{},
				errors.New("Forbidden! Either Function or Container field must be defined").Error(),
			),
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
		},
		{
			name: "with-function-only",
			spec: KeptnTaskDefinitionSpec{
				Function: &FunctionSpec{},
			},
		},
		{
			name: "with-container-only",
			spec: KeptnTaskDefinitionSpec{
				Container: &ContainerSpec{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ktd := &KeptnTaskDefinition{
				ObjectMeta: metav1.ObjectMeta{Name: tt.name},
				Spec:       tt.spec,
			}
			if got := ktd.validateFields(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("validateFields() = %v, want %v", got, tt.want)
			}
		})
	}
}
