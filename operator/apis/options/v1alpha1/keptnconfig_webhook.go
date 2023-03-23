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

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var keptnconfiglog = logf.Log.WithName("keptnconfig-resource")

func (r *KeptnConfig) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-options-keptn-sh-v1alpha1-keptnconfig,mutating=false,failurePolicy=fail,sideEffects=None,groups=options.keptn.sh,resources=keptnconfigs,verbs=create;update,versions=v1alpha1,name=vkeptnconfig.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &KeptnConfig{}

// ValidateCreate checks that there is not yet another KetpnConfig active
func (r *KeptnConfig) ValidateCreate() error {
	keptnconfiglog.Info("validate create", "name", r.Name)

	// TODO:
	// 1. Collect all KeptnConfig
	// 2. Check if # > 1 - error
	// 3. if # < 1 ok
	// 4. if # == 1 -> same name, otherwise error
	return nil
}

// ValidateUpdate immediately returns since there is nothing to validate
func (r *KeptnConfig) ValidateUpdate(old runtime.Object) error {
	keptnconfiglog.Info("validate update", "name", r.Name)
	return nil
}

// ValidateDelete immediately returns since there is nothing to validate
func (r *KeptnConfig) ValidateDelete() error {
	keptnconfiglog.Info("validate delete", "name", r.Name)
	return nil
}
