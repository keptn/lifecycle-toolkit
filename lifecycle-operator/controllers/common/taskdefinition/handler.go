package taskdefinition

import (
	"os"
	"reflect"

	klcv1beta1 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1/common"
)

const (
	FunctionRuntimeImageKey = "FUNCTION_RUNNER_IMAGE"
	PythonRuntimeImageKey   = "PYTHON_RUNNER_IMAGE"
	FunctionScriptMountPath = "/var/data/function.ts"
	PythonScriptMountPath   = "/var/data/function.py"
	FunctionScriptKey       = "js"
	PythonScriptKey         = "python"
)

func GetRuntimeSpec(def *klcv1beta1.KeptnTaskDefinition) *klcv1beta1.RuntimeSpec {
	if !IsRuntimeEmpty(def.Spec.Deno) {
		return def.Spec.Deno
	}
	if !IsRuntimeEmpty(def.Spec.Python) {
		return def.Spec.Python
	}

	return nil
}

func IsRuntimeEmpty(spec *klcv1beta1.RuntimeSpec) bool {
	return spec == nil || reflect.DeepEqual(spec, &klcv1beta1.RuntimeSpec{})
}

func IsContainerEmpty(spec *klcv1beta1.ContainerSpec) bool {
	return spec == nil || reflect.DeepEqual(spec, &klcv1beta1.ContainerSpec{})
}

func IsVolumeMountPresent(spec *klcv1beta1.ContainerSpec) bool {
	return spec != nil && spec.Container != nil && spec.VolumeMounts != nil && len(spec.VolumeMounts) > 0
}

func IsInline(spec *klcv1beta1.RuntimeSpec) bool {
	return spec != nil && !reflect.DeepEqual(spec.Inline, klcv1beta1.Inline{})
}

func GetRuntimeImage(def *klcv1beta1.KeptnTaskDefinition) string {
	image := os.Getenv(FunctionRuntimeImageKey)
	if !IsRuntimeEmpty(def.Spec.Python) && IsRuntimeEmpty(def.Spec.Deno) {
		image = os.Getenv(PythonRuntimeImageKey)
	}
	return image
}

func GetCmName(functionName string, spec *klcv1beta1.RuntimeSpec) string {
	if IsInline(spec) {
		return "keptnfn-" + apicommon.TruncateString(functionName, 245)
	}
	return spec.ConfigMapReference.Name
}

func GetRuntimeMountPath(def *klcv1beta1.KeptnTaskDefinition) string {
	path := FunctionScriptMountPath
	if !IsRuntimeEmpty(def.Spec.Python) && IsRuntimeEmpty(def.Spec.Deno) {
		path = PythonScriptMountPath
	}
	return path
}

// check if either the functions or container spec is set
func SpecExists(definition *klcv1beta1.KeptnTaskDefinition) bool {
	if definition == nil {
		return false
	}
	runtimeSpec := GetRuntimeSpec(definition)
	if runtimeSpec != nil {
		return true
	}

	return !IsContainerEmpty(definition.Spec.Container)
}
