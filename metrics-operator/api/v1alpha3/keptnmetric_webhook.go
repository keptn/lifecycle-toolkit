/*
Copyright 2023.

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
	"time"

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
var keptnmetriclog = logf.Log.WithName("keptnmetric-resource")

func (r *KeptnMetric) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-metrics-keptn-sh-v1alpha3-keptnmetric,mutating=false,failurePolicy=fail,sideEffects=None,groups=metrics.keptn.sh,resources=keptnmetrics,verbs=create;update,versions=v1alpha3,name=vkeptnmetric.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &KeptnMetric{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *KeptnMetric) ValidateCreate() error {
	keptnmetriclog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *KeptnMetric) ValidateUpdate(old runtime.Object) error {
	keptnmetriclog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *KeptnMetric) ValidateDelete() error {
	keptnmetriclog.Info("validate delete", "name", r.Name)

	return nil
}

func (s *KeptnMetric) validateKeptnMetric() error {
	var allErrs field.ErrorList //defined as a list to allow returning multiple validation errors
	var err *field.Error
	if err = s.validateRangeInterval(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "metrics.keptn.sh", Kind: "KeptnMetric"},
		s.Name, allErrs)
}

func (s *KeptnMetric) validateRangeInterval() *field.Error {
	_, err := time.ParseDuration(s.Spec.Range.Interval)
	if err != nil {
		return field.Invalid(
			field.NewPath("spec").Child("range").Child("interval"),
			s.Spec.Range.Interval,
			errors.New("Forbidden! The time interval cannot be parsed. Please check for suitable conventions").Error(),
		)
	}
	return nil
}
