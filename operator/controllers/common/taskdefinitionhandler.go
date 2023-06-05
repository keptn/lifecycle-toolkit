package common

import (
	"os"
	"reflect"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
)

func GetRuntimeSpec(def *klcv1alpha3.KeptnTaskDefinition) *klcv1alpha3.RuntimeSpec {

	if !IsRuntimeEmpty(def.Spec.Function) {
		return def.Spec.Function
	}
	if !IsRuntimeEmpty(def.Spec.Deno) {
		return def.Spec.Deno
	}
	if !IsRuntimeEmpty(def.Spec.Python) {
		return def.Spec.Python
	}

	return nil
}

func IsRuntimeEmpty(spec *klcv1alpha3.RuntimeSpec) bool {
	return spec == nil || reflect.DeepEqual(spec, &klcv1alpha3.RuntimeSpec{})
}

func IsVolumeMountPresent(spec *klcv1alpha3.ContainerSpec) bool {
	return spec != nil && spec.Container != nil && spec.VolumeMounts != nil && len(spec.VolumeMounts) > 0
}

func IsInline(spec *klcv1alpha3.RuntimeSpec) bool {
	return spec != nil && !reflect.DeepEqual(spec.Inline, klcv1alpha3.Inline{})
}

func GetRuntimeImage(def *klcv1alpha3.KeptnTaskDefinition) string {
	image := os.Getenv("FUNCTION_RUNNER_IMAGE")
	if !IsRuntimeEmpty(def.Spec.Python) && IsRuntimeEmpty(def.Spec.Function) {
		return os.Getenv("PYTHON_RUNNER_IMAGE")
	}
	return image
}

func GetRuntimeMountPath(def *klcv1alpha3.KeptnTaskDefinition) string {
	path := "/var/data/function.ts"
	if !IsRuntimeEmpty(def.Spec.Python) && IsRuntimeEmpty(def.Spec.Function) && IsRuntimeEmpty(def.Spec.Deno) {
		path = "/var/data/function.py"
	}
	return path
}

func SpecExists(definition *klcv1alpha3.KeptnTaskDefinition) bool {
	if definition == nil {
		return false
	}
	runtimeSpec := GetRuntimeSpec(definition)
	if runtimeSpec != nil {
		return true
	}
	//hasParent := runtimeSpec != nil && reflect.DeepEqual(runtimeSpec.FunctionReference, klcv1alpha3.FunctionReference{})
	return definition.Spec.Container != nil || reflect.DeepEqual(runtimeSpec, &klcv1alpha3.ContainerSpec{})
}
