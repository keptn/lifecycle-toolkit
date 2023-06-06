package common

import (
	"os"
	"reflect"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
)

const (
	FunctionRuntimeImageKey = "FUNCTION_RUNNER_IMAGE"
	PythonRuntimeImageKey   = "PYTHON_RUNNER_IMAGE"
	FunctionScriptMountPath = "/var/data/function.ts"
	PythonScriptMountPath   = "/var/data/function.py"
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

func IsContainerEmpty(spec *klcv1alpha3.ContainerSpec) bool {
	return spec == nil || reflect.DeepEqual(spec, &klcv1alpha3.ContainerSpec{})
}

func IsVolumeMountPresent(spec *klcv1alpha3.ContainerSpec) bool {
	return spec != nil && spec.Container != nil && spec.VolumeMounts != nil && len(spec.VolumeMounts) > 0
}

func IsInline(spec *klcv1alpha3.RuntimeSpec) bool {
	return spec != nil && !reflect.DeepEqual(spec.Inline, klcv1alpha3.Inline{})
}

func GetRuntimeImage(def *klcv1alpha3.KeptnTaskDefinition) string {
	image := os.Getenv(FunctionRuntimeImageKey)
	if !IsRuntimeEmpty(def.Spec.Python) && IsRuntimeEmpty(def.Spec.Function) && IsRuntimeEmpty(def.Spec.Deno) {
		image = os.Getenv(PythonRuntimeImageKey)
	}
	return image
}

func GetRuntimeMountPath(def *klcv1alpha3.KeptnTaskDefinition) string {
	path := FunctionScriptMountPath
	if !IsRuntimeEmpty(def.Spec.Python) && IsRuntimeEmpty(def.Spec.Function) && IsRuntimeEmpty(def.Spec.Deno) {
		path = PythonScriptMountPath
	}
	return path
}

// check if either the funtions or container spec is set
func SpecExists(definition *klcv1alpha3.KeptnTaskDefinition) bool {
	if definition == nil {
		return false
	}
	runtimeSpec := GetRuntimeSpec(definition)
	if runtimeSpec != nil {
		return true
	}

	return !IsContainerEmpty(definition.Spec.Container)
}
