/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/imdario/mergo"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	"math/rand"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"

	klcv1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KeptnTaskReconciler reconciles a KeptnTask object
type KeptnTaskReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Log      logr.Logger
}

//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntasks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntasks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptntasks/finalizers,verbs=update

type ExecutionParams struct {
	ConfigMap        string
	Parameters       map[string]string
	SecureParameters string
	URL              string
}

func (r *KeptnTaskReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	r.Log.Info("Reconciling KeptnTask")
	task := &klcv1alpha1.KeptnTask{}

	if err := r.Client.Get(ctx, req.NamespacedName, task); err != nil {
		if errors.IsNotFound(err) {
			// taking down all associated K8s resources is handled by K8s
			r.Log.Info("KeptnTask resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		r.Log.Error(err, "Failed to get the KeptnTask")
		return ctrl.Result{Requeue: true, RequeueAfter: 30 * time.Second}, nil
	}

	if task.Status.JobName == "" {
		definition, err := r.getTaskDefinition(ctx, task.Spec.TaskDefinition, req.Namespace)
		if err != nil {
			r.Recorder.Event(task, "Warning", "TaskDefinitionNotFound", fmt.Sprintf("Could not find KeptnTaskDefinition / Namespace: %s, Name: %s ", task.Namespace, task.Spec.TaskDefinition))
			return ctrl.Result{}, err
		}

		params, hasParent, err := r.parseTaskDefinition(definition)
		var parentJobParams ExecutionParams
		if err != nil {
			return ctrl.Result{}, err
		}
		if hasParent {
			parentDefinition, err := r.getTaskDefinition(ctx, definition.Spec.Function.FunctionReference.Name, req.Namespace)
			if err != nil {
				r.Recorder.Event(task, "Warning", "TaskDefinitionNotFound", fmt.Sprintf("Could not find KeptnTaskDefinition / Namespace: %s, Name: %s ", task.Namespace, task.Spec.TaskDefinition))
				return ctrl.Result{}, err
			}
			parentJobParams, _, err = r.parseTaskDefinition(parentDefinition)
			if err != nil {
				return ctrl.Result{}, err
			}
			err = mergo.Merge(&params, parentJobParams)
			if err != nil {
				r.Recorder.Event(task, "Warning", "TaskDefinitionMergeFailure", fmt.Sprintf("Could not merge KeptnTaskDefinition / Namespace: %s, Name: %s ", task.Namespace, task.Spec.TaskDefinition))
				return ctrl.Result{}, err
			}
		}

		if len(task.Spec.Parameters.Inline) > 0 {
			err = mergo.Merge(&params.Parameters, task.Spec.Parameters.Inline)
			if err != nil {
				r.Recorder.Event(task, "Warning", "TaskDefinitionMergeFailure", fmt.Sprintf("Could not merge KeptnTaskDefinition / Namespace: %s, Name: %s ", task.Namespace, task.Spec.TaskDefinition))
				return ctrl.Result{}, err
			}
		}
		job, _ := r.createFunctionJob(ctx, task, params)
		if err != nil {
			return ctrl.Result{}, err
		}

		task.Status.JobName = job.Name
		task.Status.State = "Pending"
		err = r.Client.Status().Update(ctx, task)
		if err != nil {
			r.Log.Error(err, "could not update configmap status reference for: "+definition.Name)
		}
		r.Log.Info("updated configmap status reference for: " + definition.Name)
	}

	if task.Status.State != "Finished" {
		job, err := r.getJob(ctx, task.Status.JobName, req.Namespace)
		if err != nil {
			r.Recorder.Event(task, "Warning", "JobNotFound", fmt.Sprintf("Could not find Job / Namespace: %s, Name: %s ", task.Namespace, task.Status.JobName))
			task.Status.JobName = ""
			err = r.Client.Status().Update(ctx, task)
			if err != nil {
				r.Log.Error(err, "could not update job reference reference for: "+task.Name)
			}
			return ctrl.Result{}, err
		}
		if job.Status.Succeeded > 0 {
			task.Status.State = "Finished"
			err = r.Client.Status().Update(ctx, task)
			if err != nil {
				r.Log.Error(err, "could not update job reference reference for: "+task.Name)
			}
		}
	}

	r.Log.Info("Finished Reconciling KeptnTask")
	return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnTaskReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&klcv1alpha1.KeptnTask{}).
		Owns(&batchv1.Job{}).
		Complete(r)
}

func (r *KeptnTaskReconciler) getTaskDefinition(ctx context.Context, definitionName string, namespace string) (*klcv1alpha1.KeptnTaskDefinition, error) {
	definition := &klcv1alpha1.KeptnTaskDefinition{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: definitionName, Namespace: namespace}, definition)
	if err != nil {
		return definition, err
	}
	return definition, nil
}

func (r *KeptnTaskReconciler) getJob(ctx context.Context, jobName string, namespace string) (*batchv1.Job, error) {
	job := &batchv1.Job{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: jobName, Namespace: namespace}, job)
	if err != nil {
		return job, err
	}
	return job, nil
}

func (r *KeptnTaskReconciler) parseTaskDefinition(definition *klcv1alpha1.KeptnTaskDefinition) (ExecutionParams, bool, error) {
	params := ExecutionParams{}

	// Firstly check if this task definition has a parent object
	hasParent := false
	if definition.Spec.Function.FunctionReference != (klcv1alpha1.FunctionReference{}) {
		hasParent = true
	}

	// Check if there is a ConfigMap with the function for this object
	if definition.Status.ConfigMap != "" {
		params.ConfigMap = definition.Status.ConfigMap
	} else {
		// If not, check if it has an HTTP reference. If this is also not the case and the object has no parent, something is wrong
		if definition.Spec.Function.HttpReference.Url == "" && !hasParent {
			return params, false, fmt.Errorf("No ConfigMap specified or HTTP source specified in TaskDefinition) / Namespace: %s, Name: %s ", definition.Namespace, definition.Name)
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

func (r *KeptnTaskReconciler) createFunctionJob(ctx context.Context, task *klcv1alpha1.KeptnTask, params ExecutionParams) (*batchv1.Job, error) {
	randomId := rand.Intn(99999-10000) + 10000
	jobId := fmt.Sprintf("klc-%s-%s-%s-%d", task.Spec.Application, task.Spec.Workload, task.Spec.WorkloadVersion, randomId)
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      jobId,
			Namespace: task.Namespace,
			Annotations: map[string]string{
				"keptn.sh/app":      task.Spec.Application,
				"keptn.sh/workload": task.Spec.Workload,
				"keptn.sh/version":  task.Spec.WorkloadVersion,
			},
		},
		Spec: batchv1.JobSpec{
			Template: corev1.PodTemplateSpec{
				Spec: corev1.PodSpec{
					RestartPolicy: "OnFailure",
				},
			},
		},
	}
	err := controllerutil.SetControllerReference(task, job, r.Scheme)
	if err != nil {
		r.Log.Error(err, "could not set controller reference:")
	}

	container := corev1.Container{
		Name:  "keptn-function-runner",
		Image: "ghcr.io/keptn-sandbox/functions-runtime:main.202209211413",
	}

	var envVars []corev1.EnvVar

	if len(params.Parameters) > 0 {
		jsonParams, err := json.Marshal(params.Parameters)
		if err != nil {
			return job, fmt.Errorf("could not marshal parameters")
		}
		envVars = append(envVars, corev1.EnvVar{Name: "DATA", Value: string(jsonParams)})
	}

	if params.SecureParameters != "" {
		envVars = append(envVars, corev1.EnvVar{
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: params.SecureParameters},
					Key:                  "code",
				},
			},
		})
	}

	if params.ConfigMap != "" {
		envVars = append(envVars, corev1.EnvVar{Name: "SCRIPT", Value: "/var/data/function.ts"})

		job.Spec.Template.Spec.Volumes = []corev1.Volume{
			{
				Name: "function-mount",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: params.ConfigMap,
						},
					},
				},
			},
		}
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
	job.Spec.Template.Spec.Containers = []corev1.Container{
		container,
	}

	err = r.Client.Create(ctx, job)
	if err != nil {
		r.Log.Error(err, "could not create job")
		r.Recorder.Event(task, "Warning", "JobNotCreated", fmt.Sprintf("Could not create Job / Namespace: %s, Name: %s ", task.Namespace, task.Name))
		return job, err
	}
	return job, nil
}
