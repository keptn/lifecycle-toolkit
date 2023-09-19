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
		Function:  &RuntimeSpec{},
		Container: &ContainerSpec{},
	}

	specWithFunctionAndPython := KeptnTaskDefinitionSpec{
		Function: &RuntimeSpec{},
		Python:   &RuntimeSpec{},
	}

	specWithFunctionAndDeno := KeptnTaskDefinitionSpec{
		Function: &RuntimeSpec{},
		Deno:     &RuntimeSpec{},
	}

	specWithContainerAndPython := KeptnTaskDefinitionSpec{
		Container: &ContainerSpec{},
		Python:    &RuntimeSpec{},
	}

	specWithContainerAndDeno := KeptnTaskDefinitionSpec{
		Container: &ContainerSpec{},
		Deno:      &RuntimeSpec{},
	}

	specWithPythonAndDeno := KeptnTaskDefinitionSpec{
		Python: &RuntimeSpec{},
		Deno:   &RuntimeSpec{},
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
			name: "with-no-function-or-container-or-python-or-deno",
			spec: emptySpec,
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "lifecycle.keptn.sh", Kind: "KeptnTaskDefinition"},
				"with-no-function-or-container-or-python-or-deno",
				[]*field.Error{field.Invalid(
					field.NewPath("spec"),
					emptySpec,
					errors.New("Forbidden! Either Function, Container, Python, or Deno field must be defined").Error(),
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
					errors.New("Forbidden! Only one of Function, Container, Python, or Deno field can be defined").Error(),
				)},
			),
		},
		{
			name: "with-function-only",
			spec: KeptnTaskDefinitionSpec{
				Function: &RuntimeSpec{},
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
			name: "with-python-only",
			spec: KeptnTaskDefinitionSpec{
				Python: &RuntimeSpec{},
			},
			verb: "create",
		},
		{
			name: "with-deno-only",
			spec: KeptnTaskDefinitionSpec{
				Deno: &RuntimeSpec{},
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
					errors.New("Forbidden! Only one of Function, Container, Python, or Deno field can be defined").Error(),
				)},
			),
			oldSpec: &KeptnTaskDefinition{
				Spec: KeptnTaskDefinitionSpec{},
			},
			verb: "update",
		},

		{
			name: "with-both-function-and-python",
			spec: specWithFunctionAndPython,
			verb: "create",
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "lifecycle.keptn.sh", Kind: "KeptnTaskDefinition"},
				"with-both-function-and-python",
				[]*field.Error{field.Invalid(
					field.NewPath("spec"),
					specWithFunctionAndPython,
					errors.New("Forbidden! Only one of Function, Container, Python, or Deno field can be defined").Error(),
				)},
			),
		},
		{
			name: "update-with-both-function-and-python",
			spec: specWithFunctionAndPython,
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "lifecycle.keptn.sh", Kind: "KeptnTaskDefinition"},
				"update-with-both-function-and-python",
				[]*field.Error{field.Invalid(
					field.NewPath("spec"),
					specWithFunctionAndPython,
					errors.New("Forbidden! Only one of Function, Container, Python, or Deno field can be defined").Error(),
				)},
			),
			oldSpec: &KeptnTaskDefinition{
				Spec: KeptnTaskDefinitionSpec{},
			},
			verb: "update",
		},

		{
			name: "with-both-function-and-deno",
			spec: specWithFunctionAndDeno,
			verb: "create",
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "lifecycle.keptn.sh", Kind: "KeptnTaskDefinition"},
				"with-both-function-and-deno",
				[]*field.Error{field.Invalid(
					field.NewPath("spec"),
					specWithFunctionAndDeno,
					errors.New("Forbidden! Only one of Function, Container, Python, or Deno field can be defined").Error(),
				)},
			),
		},
		{
			name: "update-with-both-function-and-deno",
			spec: specWithFunctionAndDeno,
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "lifecycle.keptn.sh", Kind: "KeptnTaskDefinition"},
				"update-with-both-function-and-deno",
				[]*field.Error{field.Invalid(
					field.NewPath("spec"),
					specWithFunctionAndDeno,
					errors.New("Forbidden! Only one of Function, Container, Python, or Deno field can be defined").Error(),
				)},
			),
			oldSpec: &KeptnTaskDefinition{
				Spec: KeptnTaskDefinitionSpec{},
			},
			verb: "update",
		},

		{
			name: "with-both-container-and-python",
			spec: specWithContainerAndPython,
			verb: "create",
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "lifecycle.keptn.sh", Kind: "KeptnTaskDefinition"},
				"with-both-container-and-python",
				[]*field.Error{field.Invalid(
					field.NewPath("spec"),
					specWithContainerAndPython,
					errors.New("Forbidden! Only one of Function, Container, Python, or Deno field can be defined").Error(),
				)},
			),
		},
		{
			name: "update-with-both-container-and-python",
			spec: specWithContainerAndPython,
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "lifecycle.keptn.sh", Kind: "KeptnTaskDefinition"},
				"update-with-both-container-and-python",
				[]*field.Error{field.Invalid(
					field.NewPath("spec"),
					specWithContainerAndPython,
					errors.New("Forbidden! Only one of Function, Container, Python, or Deno field can be defined").Error(),
				)},
			),
			oldSpec: &KeptnTaskDefinition{
				Spec: KeptnTaskDefinitionSpec{},
			},
			verb: "update",
		},

		{
			name: "with-both-container-and-deno",
			spec: specWithContainerAndDeno,
			verb: "create",
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "lifecycle.keptn.sh", Kind: "KeptnTaskDefinition"},
				"with-both-container-and-deno",
				[]*field.Error{field.Invalid(
					field.NewPath("spec"),
					specWithContainerAndDeno,
					errors.New("Forbidden! Only one of Function, Container, Python, or Deno field can be defined").Error(),
				)},
			),
		},
		{
			name: "update-with-both-container-and-deno",
			spec: specWithContainerAndDeno,
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "lifecycle.keptn.sh", Kind: "KeptnTaskDefinition"},
				"update-with-both-container-and-deno",
				[]*field.Error{field.Invalid(
					field.NewPath("spec"),
					specWithContainerAndDeno,
					errors.New("Forbidden! Only one of Function, Container, Python, or Deno field can be defined").Error(),
				)},
			),
			oldSpec: &KeptnTaskDefinition{
				Spec: KeptnTaskDefinitionSpec{},
			},
			verb: "update",
		},
		{
			name: "with-both-python-and-deno",
			spec: specWithPythonAndDeno,
			verb: "create",
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "lifecycle.keptn.sh", Kind: "KeptnTaskDefinition"},
				"with-both-python-and-deno",
				[]*field.Error{field.Invalid(
					field.NewPath("spec"),
					specWithPythonAndDeno,
					errors.New("Forbidden! Only one of Function, Container, Python, or Deno field can be defined").Error(),
				)},
			),
		},
		{
			name: "update-with-both-python-and-deno",
			spec: specWithPythonAndDeno,
			want: apierrors.NewInvalid(
				schema.GroupKind{Group: "lifecycle.keptn.sh", Kind: "KeptnTaskDefinition"},
				"update-with-both-python-and-deno",
				[]*field.Error{field.Invalid(
					field.NewPath("spec"),
					specWithPythonAndDeno,
					errors.New("Forbidden! Only one of Function, Container, Python, or Deno field can be defined").Error(),
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
				_, got = ktd.ValidateCreate()
			case "update":
				_, got = ktd.ValidateUpdate(tt.oldSpec)
			case "delete":
				_, got = ktd.ValidateDelete()
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
