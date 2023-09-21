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
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// Analysislog is for logging in this package.
var Analysislog = logf.Log.WithName("analysis-webhook")

func (r *Analysis) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-metrics-keptn-sh-v1alpha3-analysis,mutating=false,failurePolicy=fail,sideEffects=None,groups=metrics.keptn.sh,resources=analyses,verbs=create;update,versions=v1alpha3,name=analysis.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Analysis{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (a *Analysis) ValidateCreate() (admission.Warnings, error) {
	Analysislog.Info("validate create", "name", a.Name)

	if err := a.validateTimeframe(); err != nil {
		return []string{}, err
	}

	return []string{}, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (a *Analysis) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	Analysislog.Info("validate update", "name", a.Name)

	if err := a.validateTimeframe(); err != nil {
		return []string{}, err
	}

	return []string{}, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (a *Analysis) ValidateDelete() (admission.Warnings, error) {
	Analysislog.Info("validate delete", "name", a.Name)

	return []string{}, nil
}

func (a *Analysis) validateTimeframe() error {
	// if 'Recent'  is set, this must be the only field
	if a.Spec.Timeframe.Recent.Duration != 0 {
		if !a.Spec.Timeframe.From.IsZero() || !a.Spec.Timeframe.To.IsZero() {
			return field.Invalid(
				field.NewPath("spec").Child("timeframe"),
				a.Spec.Timeframe,
				errors.New("the field 'recent' can not be used in conjunction with 'from'/'to'").Error(),
			)
		}
		return nil
	}
	// if 'Recent' is not set, both 'From' and 'To' must be set
	if a.Spec.Timeframe.From.IsZero() || a.Spec.Timeframe.To.IsZero() {
		return field.Invalid(
			field.NewPath("spec").Child("timeframe"),
			a.Spec.Timeframe,
			errors.New("either 'recent' or both 'from' and 'to'  must be set").Error(),
		)
	}
	if !a.Spec.Timeframe.To.After(a.Spec.Timeframe.From.Time) {
		return field.Invalid(
			field.NewPath("spec").Child("timeframe"),
			a.Spec.Timeframe,
			errors.New("value of 'to' must be a timestamp later than 'from'").Error(),
		)
	}

	return nil
}
