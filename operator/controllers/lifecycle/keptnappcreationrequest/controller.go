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

package keptnappcreationrequest

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/config"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strings"
	"time"

	"github.com/go-logr/logr"
	lifecycle "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// KeptnAppCreationRequestReconciler reconciles a KeptnAppCreationRequest object
type KeptnAppCreationRequestReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnappcreationrequests,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnappcreationrequests/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=lifecycle.keptn.sh,resources=keptnappcreationrequests/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the KeptnAppCreationRequest object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *KeptnAppCreationRequestReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	creationRequest := &lifecycle.KeptnAppCreationRequest{}

	if err := r.Get(ctx, req.NamespacedName, creationRequest); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, fmt.Errorf("could not retrieve KeptnAppCreationRequest: %w", err)
	}

	// check if we already have an app that has not been created by this controller

	appFound := false
	keptnApp := &lifecycle.KeptnApp{}
	name := req.NamespacedName
	name.Name = creationRequest.Spec.AppName
	if err := r.Get(ctx, name, keptnApp); err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("No KeptnApp found for KeptnAppCreationRequest", "KeptnAppCreationRequest", creationRequest)
		} else {
			return ctrl.Result{}, fmt.Errorf("could not retrieve KeptnApp %w", err)
		}
	} else {
		appFound = true
	}

	// if the found app has not been created by this controller, we are done at this point - we don't want to mess with what the user has created
	if appFound && len(keptnApp.OwnerReferences) == 0 {
		r.Log.Info("User defined KeptnApp found for KeptnAppCreationRequest", "KeptnAppCreationRequest", creationRequest)
		if err := r.Delete(ctx, creationRequest); err != nil {
			r.Log.Error(err, "Could not delete KeptnAppCreationRequest", "KeptnAppCreationRequest", creationRequest)
		}
		return ctrl.Result{}, nil
	}

	// TODO check annotation: if the KeptnAppCreationRequest annotation keptn.sh/app-type has value single-service -> it's single-service app and the timeout is set to 0

	// check if discovery deadline has expired
	discoveryDeadline := config.Instance().GetCreationRequestTimeout() // TODO this needs to be retrieved from the KeptnConfig Object
	if !time.Now().After(creationRequest.CreationTimestamp.Add(discoveryDeadline)) {
		r.Log.Info("Discovery deadline not expired yet", "KeptnAppCreationRequest", creationRequest)
		return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	}

	// look up all the KeptnWorkloads referencing the KeptnApp

	workloads := &lifecycle.KeptnWorkloadList{}
	if err := r.Client.List(ctx, workloads, client.InNamespace(creationRequest.Namespace), client.MatchingFields{
		"spec.app": creationRequest.Spec.AppName,
	}); err != nil {
		return ctrl.Result{RequeueAfter: 10 * time.Second}, fmt.Errorf("could not retrieve KeptnWorkloads: %w", err)
	}

	var err error
	if !appFound {
		err = r.createKeptnApp(ctx, creationRequest, workloads)
	} else {
		err = r.updateKeptnApp(ctx, keptnApp, workloads)
	}

	if err != nil {
		return ctrl.Result{}, fmt.Errorf("could not update: %w", err)
	}

	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnAppCreationRequestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lifecycle.KeptnAppCreationRequest{}).
		Complete(r)
}

func (r *KeptnAppCreationRequestReconciler) updateKeptnApp(ctx context.Context, keptnApp *lifecycle.KeptnApp, workloads *lifecycle.KeptnWorkloadList) error {

	updatedVersion := false
	addedWorkload := false
	for _, workload := range workloads.Items {
		foundWorkload := false
		workloadName := strings.TrimPrefix(workload.Name, fmt.Sprintf("%s-", keptnApp.Name))
		for index, appWorkload := range keptnApp.Spec.Workloads {
			if appWorkload.Name == workloadName {
				// make sure the version matches the current version of the workload
				if keptnApp.Spec.Workloads[index].Version != workload.Spec.Version {
					keptnApp.Spec.Workloads[index].Version = workload.Spec.Version
					// we may also want to increase the version of the app if any version has been changed
					updatedVersion = true
				}
				foundWorkload = true
				break
			}
		}

		if !foundWorkload {
			keptnApp.Spec.Workloads = append(keptnApp.Spec.Workloads, lifecycle.KeptnWorkloadRef{
				Name:    workloadName,
				Version: workload.Spec.Version,
			})
			addedWorkload = true
		}
	}

	if !updatedVersion && !addedWorkload {
		return nil
	}

	if updatedVersion {
		oldVersion, _ := version.NewVersion(keptnApp.Spec.Version)
		keptnApp.Spec.Version = fmt.Sprintf("%d.0.0", oldVersion.Segments()[0]+1)
	}

	return r.Update(ctx, keptnApp)
}

func (r *KeptnAppCreationRequestReconciler) createKeptnApp(ctx context.Context, creationRequest *lifecycle.KeptnAppCreationRequest, workloads *lifecycle.KeptnWorkloadList) error {
	keptnApp := &lifecycle.KeptnApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      creationRequest.Spec.AppName,
			Namespace: creationRequest.Namespace,
		},
		Spec: lifecycle.KeptnAppSpec{
			Version:                   "1.0.0",
			PreDeploymentTasks:        []string{},
			PostDeploymentTasks:       []string{},
			PreDeploymentEvaluations:  []string{},
			PostDeploymentEvaluations: []string{},
			Workloads:                 []lifecycle.KeptnWorkloadRef{},
		},
	}

	if err := controllerutil.SetControllerReference(creationRequest, keptnApp, r.Scheme); err != nil {
		r.Log.Error(err, "could not set controller reference for KeptnApp: "+keptnApp.Name)
	}

	for _, workload := range workloads.Items {
		keptnApp.Spec.Workloads = append(keptnApp.Spec.Workloads, lifecycle.KeptnWorkloadRef{
			Name:    strings.TrimPrefix(workload.Name, fmt.Sprintf("%s-", creationRequest.Spec.AppName)),
			Version: workload.Spec.Version,
		})
	}

	return r.Create(ctx, keptnApp)
}
