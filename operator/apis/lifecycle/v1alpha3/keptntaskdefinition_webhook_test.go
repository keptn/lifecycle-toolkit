package v1alpha3

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func TestKeptnTaskDefinition_ValidateFields(t *testing.T) {

	specWithFunctionAndContainer := KeptnTaskDefinitionSpec{
		Function:  &FunctionSpec{},
		Container: &ContainerSpec{},
	}

	emptySpec := KeptnTaskDefinitionSpec{}

	tests := []struct {
		name    string
		spec    KeptnTaskDefinitionSpec
		want    error
		verb    string
		oldSpec runtime.Object
	}{
		{
			name: "with-no-function-or-container",
			spec: emptySpec,
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "lifecycle.keptn.sh", Kind: "KeptnTaskDefinition"},
				"with-no-function-or-container",
				[]*field.Error{field.Invalid(
					field.NewPath("spec"),
					emptySpec,
					errors.New("Forbidden! Either Function or Container field must be defined").Error(),
				)},
			),
			verb: "create",
		},
		{
			name: "with-both-function-and-container",
			spec: specWithFunctionAndContainer,
			verb: "create",
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "lifecycle.keptn.sh", Kind: "KeptnTaskDefinition"},
				"with-both-function-and-container",
				[]*field.Error{field.Invalid(
					field.NewPath("spec"),
					specWithFunctionAndContainer,
					errors.New("Forbidden! Both Function and Container fields cannot be defined simultaneously").Error(),
				)},
			),
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
			spec: specWithFunctionAndContainer,
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "lifecycle.keptn.sh", Kind: "KeptnTaskDefinition"},
				"update-with-both-function-and-container",
				[]*field.Error{field.Invalid(
					field.NewPath("spec"),
					specWithFunctionAndContainer,
					errors.New("Forbidden! Both Function and Container fields cannot be defined simultaneously").Error(),
				)},
			),
			oldSpec: &KeptnTaskDefinition{
				Spec: KeptnTaskDefinitionSpec{},
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

			if tt.want != nil {
				require.NotNil(t, got)
				require.EqualValues(t, tt.want, got)
			} else {
				require.Nil(t, got)
			}
		})
	}
}
