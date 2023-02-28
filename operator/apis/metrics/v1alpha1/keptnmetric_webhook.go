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
	"github.com/keptn/lifecycle-toolkit/operator/apis/metrics/v1alpha1/common"
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

func (s *KeptnMetric) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(s).
		Complete()
}

// +kubebuilder:webhook:path=/validate-metrics-keptn-sh-v1alpha1-keptnmetric,mutating=false,failurePolicy=fail,sideEffects=None,groups=metrics.keptn.sh,resources=keptnmetrics,verbs=create;update,versions=v1alpha1,name=vkeptnmetric.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &KeptnMetric{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (s *KeptnMetric) ValidateCreate() error {
	keptnmetriclog.Info("validate create", "name", s.Name)
	return s.validateKeptnMetric()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (s *KeptnMetric) ValidateUpdate(old runtime.Object) error {
	keptnmetriclog.Info("validate update", "name", s.Name)
	return s.validateKeptnMetric()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (s *KeptnMetric) ValidateDelete() error {
	keptnmetriclog.Info("validate delete", "name", s.Name)

	return nil
}

func (s *KeptnMetric) validateKeptnMetric() error {
	var allErrs field.ErrorList // defined as a list to allow returning multiple validation errors
	var err *field.Error
	if err = s.validateProvider(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		schema.GroupKind{Group: "metrics.keptn.sh", Kind: "KeptnMetric"},
		s.Name, allErrs)
}

func (s *KeptnMetric) validateProvider() *field.Error {
	// The field helpers from the kubernetes API machinery help us return nicely
	// structured validation errors.
	// we explicitly use spec.provider.name to make sure the error path corresponds
	if !s.IsProviderValid(s.Spec.Provider.Name) {
		return field.Invalid(
			field.NewPath("spec").Child("provider").Child("name"),
			s.Spec.Provider.Name,
			common.ErrForbiddenProvider.Error(),
		)
	}
	return nil
}
