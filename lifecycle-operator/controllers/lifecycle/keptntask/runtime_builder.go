package keptntask

import (
	"encoding/json"
	"fmt"

	"dario.cat/mergo"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	controllererrors "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/errors"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
)

// RuntimeBuilder implements container builder interface for Deno/Python
type RuntimeBuilder struct {
	options BuilderOptions
}

func NewRuntimeBuilder(options BuilderOptions) *RuntimeBuilder {

	return &RuntimeBuilder{
		options: options,
	}
}

// RuntimeExecutionParams stores parameters related to Deno/Python container creation
type RuntimeExecutionParams struct {
	ConfigMap        string
	Parameters       map[string]string
	SecureParameters string
	CmdParameters    string
	URL              string
	Context          klcv1alpha3.TaskContext
	Image            string
	MountPath        string
}

const (
	Context           = "CONTEXT"
	SecureData        = "SECURE_DATA"
	Data              = "DATA"
	CmdArgs           = "CMD_ARGS"
	Script            = "SCRIPT"
	FunctionMountName = "function-mount"
)

func (fb *RuntimeBuilder) CreateContainer(ctx context.Context) (*corev1.Container, error) {

	var envVars []corev1.EnvVar

	params, err := fb.getParams(ctx)
	if err != nil {
		return nil, err
	}
	if len(params.Parameters) > 0 {
		jsonParams, err := json.Marshal(params.Parameters)
		if err != nil {
			return nil, err
		}
		envVars = append(envVars, corev1.EnvVar{Name: Data, Value: string(jsonParams)})
	}

	jsonParams, err := json.Marshal(params.Context)
	if err != nil {
		return nil, err
	}
	envVars = append(envVars, corev1.EnvVar{Name: Context, Value: string(jsonParams)})
	envVars = append(envVars, corev1.EnvVar{Name: CmdArgs, Value: params.CmdParameters})
	if params.SecureParameters != "" {
		envVars = append(envVars, corev1.EnvVar{
			Name: SecureData,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: params.SecureParameters},
					Key:                  SecureData,
				},
			},
		})
	}

	// Mount the function code if a ConfigMap is provided
	// The ConfigMap might be provided manually or created by the TaskDefinition controller
	container := corev1.Container{
		ImagePullPolicy: corev1.PullIfNotPresent,
		Name:            "keptn-function-runner",
		Image:           params.Image,
	}

	if params.ConfigMap != "" {
		envVars = append(envVars, corev1.EnvVar{Name: Script, Value: params.MountPath})

		container.VolumeMounts = []corev1.VolumeMount{
			{
				Name:      FunctionMountName,
				ReadOnly:  true,
				MountPath: params.MountPath,
				SubPath:   "code",
			},
		}
	} else {
		envVars = append(envVars, corev1.EnvVar{Name: Script, Value: params.URL})
	}

	container.Env = envVars
	return &container, nil

}

//nolint:nilnil
func (fb *RuntimeBuilder) CreateVolume(ctx context.Context) (*corev1.Volume, error) {
	params, err := fb.getParams(ctx)
	if err != nil {
		return nil, err
	}

	if params.ConfigMap != "" {
		return &corev1.Volume{
			Name: FunctionMountName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: params.ConfigMap,
					},
				},
			},
		}, nil
	}

	return nil, nil

}

func (fb *RuntimeBuilder) getParams(ctx context.Context) (*RuntimeExecutionParams, error) {
	params, hasParent, err := fb.parseRuntimeTaskDefinition(
		fb.options.funcSpec,
		fb.options.task.Spec.TaskDefinition,
		fb.options.task.Namespace,
		fb.options.ConfigMap,
	)
	if err != nil {
		return nil, err
	}
	// set image based on child specs
	params.Image = fb.options.Image
	params.MountPath = fb.options.MountPath

	if hasParent {
		if err := fb.handleParent(ctx, &params); err != nil {
			return nil, err
		}
	}

	params.Context = setupTaskContext(fb.options.task)

	if len(fb.options.task.Spec.Parameters.Inline) > 0 {
		err = mergo.Merge(&params.Parameters, fb.options.task.Spec.Parameters.Inline)
		if err != nil {
			fb.options.eventSender.Emit(apicommon.PhaseCreateTask, "Warning", fb.options.task, apicommon.PhaseStateFailed, fmt.Sprintf("could not merge KeptnTaskDefinition: %s ", fb.options.task.Spec.TaskDefinition), "")
			return nil, err
		}
	}

	if fb.options.task.Spec.SecureParameters.Secret != "" {
		params.SecureParameters = fb.options.task.Spec.SecureParameters.Secret
	}
	return &params, nil
}

func (fb *RuntimeBuilder) parseRuntimeTaskDefinition(spec *klcv1alpha3.RuntimeSpec, name string, namespace string, configMap string) (RuntimeExecutionParams, bool, error) {
	params := RuntimeExecutionParams{}

	// Firstly check if this task definition has a parent object
	hasParent := false
	if spec.FunctionReference != (klcv1alpha3.FunctionReference{}) {
		hasParent = true
	}
	params.ConfigMap = configMap
	if params.ConfigMap != "" && spec.HttpReference.Url != "" {
		fb.options.Log.Info(fmt.Sprintf("The JobDefinition contains a ConfigMap and a HTTP Reference, ConfigMap is used / Namespace: %s, Name: %s  ", namespace, name))
	}

	// Check if there is a ConfigMap with the function for this object
	if params.ConfigMap == "" {
		// If not, check if it has an HTTP reference. If this is also not the case and the object has no parent, something is wrong
		if spec.HttpReference.Url == "" && !hasParent {
			return params, false, fmt.Errorf(controllererrors.ErrNoConfigMapMsg, namespace, name)
		}
		params.URL = spec.HttpReference.Url
	}

	// Check if there are parameters provided
	if len(spec.Parameters.Inline) > 0 {
		params.Parameters = spec.Parameters.Inline
	}

	// Check if there is a secret for secret params provided
	if spec.SecureParameters.Secret != "" {
		params.SecureParameters = spec.SecureParameters.Secret
	}

	// Check if there is a cmd params provided
	if spec.CmdParameters != "" {
		params.CmdParameters = spec.CmdParameters
	}
	return params, hasParent, nil
}

func (fb *RuntimeBuilder) handleParent(ctx context.Context, params *RuntimeExecutionParams) error {
	var parentJobParams RuntimeExecutionParams
	parentDefinition, err := controllercommon.GetTaskDefinition(fb.options.Client, fb.options.Log, ctx, fb.options.funcSpec.FunctionReference.Name, fb.options.req.Namespace)
	if err != nil {
		fb.options.eventSender.Emit(apicommon.PhaseCreateTask, "Warning", fb.options.task, apicommon.PhaseStateNotFound, fmt.Sprintf("could not find KeptnTaskDefinition: %s ", fb.options.task.Spec.TaskDefinition), "")
		return err
	}
	parSpec := controllercommon.GetRuntimeSpec(parentDefinition)
	// if the parent has also another parent, the data from the grandparent are alredy copied to the parent and therefore parent can copy it's data to the child
	parentJobParams, _, err = fb.parseRuntimeTaskDefinition(parSpec, parentDefinition.Name, parentDefinition.Namespace, parentDefinition.Status.Function.ConfigMap)
	if err != nil {
		return err
	}
	// merge parameter to make sure we use child task data for env var and secrets
	err = mergo.Merge(params, parentJobParams)
	if err != nil {
		fb.options.eventSender.Emit(apicommon.PhaseCreateTask, "Warning", fb.options.task, apicommon.PhaseStateFailed, fmt.Sprintf("could not merge KeptnTaskDefinition: %s ", fb.options.task.Spec.TaskDefinition), "")
		return err
	}

	// make sure we take the task from the parent
	params.URL = parSpec.HttpReference.Url
	params.ConfigMap = parentDefinition.Status.Function.ConfigMap

	// rewrite image and mount based on parent
	params.Image = controllercommon.GetRuntimeImage(parentDefinition)
	params.MountPath = controllercommon.GetRuntimeMountPath(parentDefinition)

	return nil
}
