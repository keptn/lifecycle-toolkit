package keptntask

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/imdario/mergo"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
)

// JSBuilder implements container builder interface for javascript deno
type JSBuilder struct {
	options BuilderOptions
}

func newJSBuilder(options BuilderOptions) *JSBuilder {
	return &JSBuilder{
		options: options,
	}
}

// FunctionExecutionParams stores parametersrelatedto js deno container creation
type FunctionExecutionParams struct {
	ConfigMap        string
	Parameters       map[string]string
	SecureParameters string
	URL              string
	Context          klcv1alpha3.TaskContext
}

func (js *JSBuilder) CreateContainerWithVolumes(ctx context.Context) (*corev1.Container, []corev1.Volume, error) {
	container := corev1.Container{
		Name:  "keptn-function-runner",
		Image: os.Getenv("FUNCTION_RUNNER_IMAGE"),
	}

	var envVars []corev1.EnvVar

	params, err := js.getParams(ctx)
	if err != nil {
		return nil, nil, err
	}
	if len(params.Parameters) > 0 {
		jsonParams, err := json.Marshal(params.Parameters)
		if err != nil {
			return nil, nil, err
		}
		envVars = append(envVars, corev1.EnvVar{Name: "DATA", Value: string(jsonParams)})
	}

	jsonParams, err := json.Marshal(params.Context)
	if err != nil {
		return nil, nil, err
	}
	envVars = append(envVars, corev1.EnvVar{Name: "CONTEXT", Value: string(jsonParams)})

	if params.SecureParameters != "" {
		envVars = append(envVars, corev1.EnvVar{
			Name: "SECURE_DATA",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: params.SecureParameters},
					Key:                  "SECURE_DATA",
				},
			},
		})
	}
	var jobVolumes []corev1.Volume
	// Mount the function code if a ConfigMap is provided
	// The ConfigMap might be provided manually or created by the TaskDefinition controller
	if params.ConfigMap != "" {
		envVars = append(envVars, corev1.EnvVar{Name: "SCRIPT", Value: "/var/data/function.ts"})

		jobVolumes = append(jobVolumes, corev1.Volume{
			Name: "function-mount",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: params.ConfigMap,
					},
				},
			},
		})

		container.VolumeMounts = []corev1.VolumeMount{
			{
				Name:      "function-mount",
				ReadOnly:  true,
				MountPath: "/var/data/function.ts",
				SubPath:   "code",
			},
		}
	} else {
		envVars = append(envVars, corev1.EnvVar{Name: "SCRIPT", Value: params.URL})
	}

	container.Env = envVars
	return &container, jobVolumes, nil

}

func (js *JSBuilder) getParams(ctx context.Context) (*FunctionExecutionParams, error) {
	params, hasParent, err := js.parseFunctionTaskDefinition(js.options.taskDef)
	if err != nil {
		return nil, err
	}
	if hasParent {
		if err := js.handleParent(ctx, &params); err != nil {
			return nil, err
		}
	}

	params.Context = setupTaskContext(js.options.task)

	if len(js.options.task.Spec.Parameters.Inline) > 0 {
		err = mergo.Merge(&params.Parameters, js.options.task.Spec.Parameters.Inline)
		if err != nil {
			controllercommon.RecordEvent(js.options.recorder, apicommon.PhaseCreateTask, "Warning", js.options.task, "TaskDefinitionMergeFailure", fmt.Sprintf("could not merge KeptnTaskDefinition: %s ", js.options.task.Spec.TaskDefinition), "")
			return nil, err
		}
	}

	if js.options.task.Spec.SecureParameters.Secret != "" {
		params.SecureParameters = js.options.task.Spec.SecureParameters.Secret
	}
	return &params, nil
}

func (js *JSBuilder) parseFunctionTaskDefinition(definition *klcv1alpha3.KeptnTaskDefinition) (FunctionExecutionParams, bool, error) {
	params := FunctionExecutionParams{}

	// Firstly check if this task definition has a parent object
	hasParent := false
	if definition.Spec.Function.FunctionReference != (klcv1alpha3.FunctionReference{}) {
		hasParent = true
	}

	if definition.Status.Function.ConfigMap != "" && definition.Spec.Function.HttpReference.Url != "" {
		js.options.Log.Info(fmt.Sprintf("The JobDefinition contains a ConfigMap and a HTTP Reference, ConfigMap is used / Namespace: %s, Name: %s  ", definition.Namespace, definition.Name))
	}

	// Check if there is a ConfigMap with the function for this object
	if definition.Status.Function.ConfigMap != "" {
		params.ConfigMap = definition.Status.Function.ConfigMap
	} else {
		// If not, check if it has an HTTP reference. If this is also not the case and the object has no parent, something is wrong
		if definition.Spec.Function.HttpReference.Url == "" && !hasParent {
			return params, false, fmt.Errorf(controllererrors.ErrNoConfigMapMsg, definition.Namespace, definition.Name)
		}
		params.URL = definition.Spec.Function.HttpReference.Url
	}

	// Check if there are parameters provided
	if len(definition.Spec.Function.Parameters.Inline) > 0 {
		params.Parameters = definition.Spec.Function.Parameters.Inline
	}

	// Check if there is a secret for secret params provided
	if definition.Spec.Function.SecureParameters.Secret != "" {
		params.SecureParameters = definition.Spec.Function.SecureParameters.Secret
	}
	return params, hasParent, nil
}

func (js *JSBuilder) handleParent(ctx context.Context, params *FunctionExecutionParams) error {
	var parentJobParams FunctionExecutionParams
	parentDefinition, err := controllercommon.GetTaskDefinition(js.options.Client, js.options.Log, ctx, js.options.taskDef.Spec.Function.FunctionReference.Name, js.options.req.Namespace)
	if err != nil {
		controllercommon.RecordEvent(js.options.recorder, apicommon.PhaseCreateTask, "Warning", js.options.task, "TaskDefinitionNotFound", fmt.Sprintf("could not find KeptnTaskDefinition: %s ", js.options.task.Spec.TaskDefinition), "")
		return err
	}
	parentJobParams, _, err = js.parseFunctionTaskDefinition(parentDefinition)
	if err != nil {
		return err
	}
	err = mergo.Merge(params, parentJobParams)
	if err != nil {
		controllercommon.RecordEvent(js.options.recorder, apicommon.PhaseCreateTask, "Warning", js.options.task, "TaskDefinitionMergeFailure", fmt.Sprintf("could not merge KeptnTaskDefinition: %s ", js.options.task.Spec.TaskDefinition), "")
		return err
	}
	return nil
}
