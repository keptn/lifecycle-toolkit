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
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var analysisdefinitionlog = logf.Log.WithName("analysisdefinition-webhook")

func (r *AnalysisDefinition) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-metrics-keptn-sh-v1alpha3-analysisdefinition,mutating=false,failurePolicy=fail,sideEffects=None,groups=metrics.keptn.sh,resources=analysisdefinitions,verbs=create;update,versions=v1alpha3,name=vanalysisdefinition.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &AnalysisDefinition{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *AnalysisDefinition) ValidateCreate() error {
	analysisdefinitionlog.Info("validate create", "name", r.Name)

	for _, o := range r.Spec.Objectives {
		if err := o.validate(); err != nil {
			return err
		}
	}

	return r.Spec.TotalScore.validate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *AnalysisDefinition) ValidateUpdate(old runtime.Object) error {
	analysisdefinitionlog.Info("validate update", "name", r.Name)

	for _, o := range r.Spec.Objectives {
		if err := o.validate(); err != nil {
			return err
		}
	}

	return r.Spec.TotalScore.validate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *AnalysisDefinition) ValidateDelete() error {
	analysisdefinitionlog.Info("validate delete", "name", r.Name)

	return nil
}

func (s *Score) validate() error {
	if s.WarningPercentage >= s.PassPercentage {
		return fmt.Errorf("Warn percentage score cannot be higher or equal than Pass percentage score")
	}
	return nil
}

func (o *Objective) validate() error {
	if err := o.Target.validate(); err != nil {
		return err
	}
	return nil
}

func (t *Target) validate() error {
	if t.Failure == nil && t.Warning != nil {
		return fmt.Errorf("Warning criteria cannot be set without Failure criteria")
	}
	if t.Failure != nil {
		if err := t.Failure.validate(); err != nil {
			return err
		}
	}

	if t.Warning != nil {
		if err := t.Warning.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (o *Operator) validate() error {
	counter := 0
	if o.LessThan != nil {
		counter++
	}
	if o.LessThanOrEqual != nil {
		counter++
	}
	if o.GreaterThan != nil {
		counter++
	}
	if o.GreaterThanOrEqual != nil {
		counter++
	}
	if o.EqualTo != nil {
		counter++
	}
	if counter > 1 {
		return fmt.Errorf("Operator: multiple operators can not be set")
	}
	if counter == 0 {
		return fmt.Errorf("Operator: no operator set")
	}
	return nil
}
