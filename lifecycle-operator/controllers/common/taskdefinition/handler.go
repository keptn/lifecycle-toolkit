package taskdefinition

import (
	"os"
	"reflect"

	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1/common"
)

const (
	FunctionRuntimeImageKey = "FUNCTION_RUNNER_IMAGE"
	PythonRuntimeImageKey   = "PYTHON_RUNNER_IMAGE"
	FunctionScriptMountPath = "/var/data/function.ts"
	PythonScriptMountPath   = "/var/data/function.py"
	FunctionScriptKey       = "js"
	PythonScriptKey         = "python"
)

func GetRuntimeSpec(def *apilifecycle.KeptnTaskDefinition) *apilifecycle.RuntimeSpec {
	if !IsRuntimeEmpty(def.Spec.Deno) {
		return def.Spec.Deno
	}
	if !IsRuntimeEmpty(def.Spec.Python) {
		return def.Spec.Python
	}

	return nil
}

func IsRuntimeEmpty(spec *apilifecycle.RuntimeSpec) bool {
	return spec == nil || reflect.DeepEqual(spec, &apilifecycle.RuntimeSpec{})
}

func IsContainerEmpty(spec *apilifecycle.ContainerSpec) bool {
	return spec == nil || reflect.DeepEqual(spec, &apilifecycle.ContainerSpec{})
}

func IsVolumeMountPresent(spec *apilifecycle.ContainerSpec) bool {
	return spec != nil && spec.Container != nil && spec.VolumeMounts != nil && len(spec.VolumeMounts) > 0
}

func IsInline(spec *apilifecycle.RuntimeSpec) bool {
	return spec != nil && !reflect.DeepEqual(spec.Inline, apilifecycle.Inline{})
}

func GetRuntimeImage(def *apilifecycle.KeptnTaskDefinition) string {
	image := os.Getenv(FunctionRuntimeImageKey)
	if !IsRuntimeEmpty(def.Spec.Python) && IsRuntimeEmpty(def.Spec.Deno) {
		image = os.Getenv(PythonRuntimeImageKey)
	}
	return image
}

func GetCmName(functionName string, spec *apilifecycle.RuntimeSpec) string {
	if IsInline(spec) {
		return "keptnfn-" + apicommon.TruncateString(functionName, 245)
	}
	return spec.ConfigMapReference.Name
}

func GetRuntimeMountPath(def *apilifecycle.KeptnTaskDefinition) string {
	path := FunctionScriptMountPath
	if !IsRuntimeEmpty(def.Spec.Python) && IsRuntimeEmpty(def.Spec.Deno) {
		path = PythonScriptMountPath
	}
	return path
}

// check if either the functions or container spec is set
func SpecExists(definition *apilifecycle.KeptnTaskDefinition) bool {
	if definition == nil {
		return false
	}
	runtimeSpec := GetRuntimeSpec(definition)
	if runtimeSpec != nil {
		return true
	}

	return !IsContainerEmpty(definition.Spec.Container)
}
