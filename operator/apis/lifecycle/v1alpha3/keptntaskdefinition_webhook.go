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

package v1alpha3

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var keptntaskdefinitionlog = logf.Log.WithName("keptntaskdefinition-resource")

func (r *KeptnTaskDefinition) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-lifecycle-keptn-sh-v1alpha3-keptntaskdefinition,mutating=false,failurePolicy=fail,sideEffects=None,groups=lifecycle.keptn.sh,resources=keptntaskdefinitions,verbs=create;update,versions=v1alpha3,name=vkeptntaskdefinition.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &KeptnTaskDefinition{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *KeptnTaskDefinition) ValidateCreate() error {
	keptntaskdefinitionlog.Info("validate create", "name", r.Name)

	return r.validateKeptnTaskDefination()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *KeptnTaskDefinition) ValidateUpdate(old runtime.Object) error {
	keptntaskdefinitionlog.Info("validate update", "name", r.Name)

	return r.validateKeptnTaskDefination()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *KeptnTaskDefinition) ValidateDelete() error {
	keptntaskdefinitionlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}

func (r *KeptnTaskDefinition) validateFields() error {
	if r.Spec.Function == nil && r.Spec.Container == nil {
		return ValidationError{Field: "spec", Message: "Either Function or Container field must be defined"}
	}

	if r.Spec.Function != nil && r.Spec.Container != nil {
		return ValidationError{Field: "spec", Message: "Both Function and Container fields cannot be defined simultaneously"}
	}

	return nil
}

// ValidationError represents a validation error with a specific field and message
type ValidationError struct {
	Field   string
	Message string
}

// Error returns the validation error message
func (e ValidationError) Error() string {
	return e.Message
}
