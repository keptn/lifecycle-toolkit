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

func NewJSBuilder(options BuilderOptions) *JSBuilder {
	return &JSBuilder{
		options: options,
	}
}

// FunctionExecutionParams stores parameters related to js deno container creation
type FunctionExecutionParams struct {
	ConfigMap        string
	Parameters       map[string]string
	SecureParameters string
	CmdParameters    string
	URL              string
	Context          klcv1alpha3.TaskContext
}

func (fb *FunctionBuilder) CreateContainerWithVolumes(ctx context.Context) (*corev1.Container, []corev1.Volume, error) {

	container := corev1.Container{
		ImagePullPolicy: corev1.PullIfNotPresent,
		Name:  "keptn-function-runner",
		Image: fb.options.taskDef.GetImage(),
	}

	var envVars []corev1.EnvVar

	params, err := fb.getParams(ctx)
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
	envVars = append(envVars, corev1.EnvVar{Name: "CMD_ARGS", Value: params.CmdParameters})
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
		envVars = append(envVars, corev1.EnvVar{Name: "SCRIPT", Value: fb.options.taskDef.GetMountPath()})

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
				MountPath: fb.options.taskDef.GetMountPath(),
				SubPath:   "code",
			},
		}
	} else {
		envVars = append(envVars, corev1.EnvVar{Name: "SCRIPT", Value: params.URL})
	}

	container.Env = envVars
	return &container, jobVolumes, nil

}

func (fb *FunctionBuilder) getParams(ctx context.Context) (*FunctionExecutionParams, error) {
	params, hasParent, err := fb.parseFunctionTaskDefinition(fb.options.taskDef)
	if err != nil {
		return nil, err
	}
	if hasParent {
		if err := fb.handleParent(ctx, &params); err != nil {
			return nil, err
		}
	}

	params.Context = setupTaskContext(fb.options.task)

	if len(fb.options.task.Spec.Parameters.Inline) > 0 {
		err = mergo.Merge(&params.Parameters, fb.options.task.Spec.Parameters.Inline)
		if err != nil {
			controllercommon.RecordEvent(fb.options.recorder, apicommon.PhaseCreateTask, "Warning", fb.options.task, "TaskDefinitionMergeFailure", fmt.Sprintf("could not merge KeptnTaskDefinition: %s ", fb.options.task.Spec.TaskDefinition), "")
			return nil, err
		}
	}

	if fb.options.task.Spec.SecureParameters.Secret != "" {
		params.SecureParameters = fb.options.task.Spec.SecureParameters.Secret
	}
	return &params, nil
}

func (fb *FunctionBuilder) parseFunctionTaskDefinition(definition *klcv1alpha3.KeptnTaskDefinition) (FunctionExecutionParams, bool, error) {
	params := FunctionExecutionParams{}

	// Firstly check if this task definition has a parent object
	hasParent := false
	if definition.Spec.Function.FunctionReference != (klcv1alpha3.FunctionReference{}) {
		hasParent = true
	}

	if definition.Status.Function.ConfigMap != "" && definition.Spec.Function.HttpReference.Url != "" {
		fb.options.Log.Info(fmt.Sprintf("The JobDefinition contains a ConfigMap and a HTTP Reference, ConfigMap is used / Namespace: %s, Name: %s  ", definition.Namespace, definition.Name))
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

	// Check if there is a cmd params provided
	if definition.Spec.Function.CmdParameters != "" {
		params.CmdParameters = definition.Spec.Function.CmdParameters
	}
	return params, hasParent, nil
}

func (fb *FunctionBuilder) handleParent(ctx context.Context, params *FunctionExecutionParams) error {
	var parentJobParams FunctionExecutionParams
	parentDefinition, err := controllercommon.GetTaskDefinition(fb.options.Client, fb.options.Log, ctx, fb.options.taskDef.Spec.Function.FunctionReference.Name, fb.options.req.Namespace)
	if err != nil {
		controllercommon.RecordEvent(fb.options.recorder, apicommon.PhaseCreateTask, "Warning", fb.options.task, "TaskDefinitionNotFound", fmt.Sprintf("could not find KeptnTaskDefinition: %s ", fb.options.task.Spec.TaskDefinition), "")
		return err
	}
	parentJobParams, _, err = fb.parseFunctionTaskDefinition(parentDefinition)
	if err != nil {
		return err
	}
	// merge parameter to make sure we use child task data for env var and secrets
	err = mergo.Merge(params, parentJobParams)
	if err != nil {
		controllercommon.RecordEvent(fb.options.recorder, apicommon.PhaseCreateTask, "Warning", fb.options.task, "TaskDefinitionMergeFailure", fmt.Sprintf("could not merge KeptnTaskDefinition: %s ", fb.options.task.Spec.TaskDefinition), "")
		return err
	}

	// make sure we take the task from the parent
	params.URL = parentDefinition.Spec.Function.HttpReference.Url
	params.ConfigMap = parentDefinition.Spec.Function.ConfigMapReference.Name

	// the task definition needs to inherit the runtime of the parent
	fb.options.taskDef.Spec.Function.FunctionRuntime = parentDefinition.Spec.Function.FunctionRuntime

	return nil
}
