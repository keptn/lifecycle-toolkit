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
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/go-logr/logr"
	lifecycle "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/config"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	managedByKLT = "klt"
)

// KeptnAppCreationRequestReconciler reconciles a KeptnAppCreationRequest object
type KeptnAppCreationRequestReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
	clock  clock.Clock
	config config.IConfig
}

func NewReconciler(client client.Client, scheme *runtime.Scheme, log logr.Logger) *KeptnAppCreationRequestReconciler {
	return &KeptnAppCreationRequestReconciler{
		Client: client,
		Scheme: scheme,
		Log:    log,
		config: config.Instance(),
		clock:  clock.New(),
	}
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
//
//nolint:gocyclo
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
	if appFound && keptnApp.Labels[common.K8sRecommendedManagedByAnnotations] != managedByKLT {
		r.Log.Info("User defined KeptnApp found for KeptnAppCreationRequest", "KeptnAppCreationRequest", creationRequest)
		if err := r.Delete(ctx, creationRequest); err != nil {
			r.Log.Error(err, "Could not delete KeptnAppCreationRequest", "KeptnAppCreationRequest", creationRequest)
		}
		return ctrl.Result{}, nil
	}

	// check if discovery deadline has expired or if the application is a single service app
	if !r.shouldCreateApp(creationRequest) {
		r.Log.Info("Discovery deadline not expired yet", "KeptnAppCreationRequest", creationRequest)
		return ctrl.Result{RequeueAfter: r.getCreationRequestExpirationDuration(creationRequest)}, nil
	}

	// look up all the KeptnWorkloads referencing the KeptnApp
	workloads := &lifecycle.KeptnWorkloadList{}
	if err := r.Client.List(ctx, workloads, client.InNamespace(creationRequest.Namespace), client.MatchingFields{
		"spec.app": creationRequest.Spec.AppName,
	}); err != nil {
		return ctrl.Result{}, fmt.Errorf("could not retrieve KeptnWorkloads: %w", err)
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
	if err := r.Delete(ctx, creationRequest); err != nil {
		r.Log.Error(err, "Could not delete", "KeptnAppCreationRequest", creationRequest)
	}
	return ctrl.Result{}, nil
}

func (r *KeptnAppCreationRequestReconciler) getCreationRequestExpirationDuration(cr *lifecycle.KeptnAppCreationRequest) time.Duration {
	creationRequestTimeout := r.config.GetCreationRequestTimeout()
	deadline := cr.CreationTimestamp.Add(creationRequestTimeout)

	duration := deadline.Sub(r.clock.Now())

	// make sure we return a non-negative duration
	if duration >= 0 {
		return duration
	}
	return 0
}

func (r *KeptnAppCreationRequestReconciler) shouldCreateApp(creationRequest *lifecycle.KeptnAppCreationRequest) bool {
	discoveryDeadline := r.config.GetCreationRequestTimeout()
	return creationRequest.IsSingleService() || r.clock.Now().After(creationRequest.CreationTimestamp.Add(discoveryDeadline))
}

// SetupWithManager sets up the controller with the Manager.
func (r *KeptnAppCreationRequestReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&lifecycle.KeptnAppCreationRequest{}).
		Complete(r)
}

func (r *KeptnAppCreationRequestReconciler) updateKeptnApp(ctx context.Context, keptnApp *lifecycle.KeptnApp, workloads *lifecycle.KeptnWorkloadList) error {

	updated := false

	for _, workload := range workloads.Items {
		foundWorkload := false
		workloadName := strings.TrimPrefix(workload.Name, fmt.Sprintf("%s-", keptnApp.Name))
		for index, appWorkload := range keptnApp.Spec.Workloads {
			if appWorkload.Name == workloadName {
				// make sure the version matches the current version of the workload
				if keptnApp.Spec.Workloads[index].Version != workload.Spec.Version {
					keptnApp.Spec.Workloads[index].Version = workload.Spec.Version
					// we may also want to increase the version of the app if any version has been changed
					updated = true
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
			updated = true
		}
	}

	if !updated {
		return nil
	}

	keptnApp.Spec.Version = computeVersionFromWorkloads(workloads.Items)

	return r.Update(ctx, keptnApp)
}

func (r *KeptnAppCreationRequestReconciler) createKeptnApp(ctx context.Context, creationRequest *lifecycle.KeptnAppCreationRequest, workloads *lifecycle.KeptnWorkloadList) error {
	keptnApp := &lifecycle.KeptnApp{
		ObjectMeta: metav1.ObjectMeta{
			Name:      creationRequest.Spec.AppName,
			Namespace: creationRequest.Namespace,
			Labels: map[string]string{
				common.K8sRecommendedManagedByAnnotations: managedByKLT,
			},
		},
		Spec: lifecycle.KeptnAppSpec{
			Version:                   computeVersionFromWorkloads(workloads.Items),
			PreDeploymentTasks:        []string{},
			PostDeploymentTasks:       []string{},
			PreDeploymentEvaluations:  []string{},
			PostDeploymentEvaluations: []string{},
			Workloads:                 []lifecycle.KeptnWorkloadRef{},
		},
	}

	for _, workload := range workloads.Items {
		keptnApp.Spec.Workloads = append(keptnApp.Spec.Workloads, lifecycle.KeptnWorkloadRef{
			Name:    strings.TrimPrefix(workload.Name, fmt.Sprintf("%s-", creationRequest.Spec.AppName)),
			Version: workload.Spec.Version,
		})
	}

	return r.Create(ctx, keptnApp)
}

func computeVersionFromWorkloads(workloads []lifecycle.KeptnWorkload) string {
	versionString := ""

	// iterate over all workloads and add their names + version
	for _, workload := range workloads {
		versionString += workload.Name + "-" + workload.Spec.Version
	}

	// take the string containing all workloads/versions and compute a hash
	hash := sha256.New()
	hash.Write([]byte(versionString))
	hashValue := fmt.Sprintf("%x", hash.Sum(nil))

	return common.TruncateString(hashValue, 10)
}
