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
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
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

	return r.validateKeptnTaskDefinition()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *KeptnTaskDefinition) ValidateUpdate(old runtime.Object) error {
	keptntaskdefinitionlog.Info("validate update", "name", r.Name)

	return r.validateKeptnTaskDefinition()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *KeptnTaskDefinition) ValidateDelete() error {
	keptntaskdefinitionlog.Info("validate delete", "name", r.Name)

	return nil
}

func (r *KeptnTaskDefinition) validateKeptnTaskDefinition() error {
	var allErrs field.ErrorList //defined as a list to allow returning multiple validation errors
	var err *field.Error
	if err = r.validateFields(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "lifecycle.keptn.sh", Kind: "KeptnTaskDefinition"},
		r.Name,
		allErrs)
}
func (r *KeptnTaskDefinition) validateFields() *field.Error {

	if r.Spec.Function == nil && r.Spec.Container == nil {
		return field.Invalid(
			field.NewPath("spec"),
			r.Spec,
			errors.New("Forbidden! Either Function or Container field must be defined").Error(),
		)
	}

	if r.Spec.Function != nil && r.Spec.Container != nil {
		return field.Invalid(
			field.NewPath("spec"),
			r.Spec,
			errors.New("Forbidden! Both Function and Container fields cannot be defined simultaneously").Error(),
		)
	}

	return nil
}
