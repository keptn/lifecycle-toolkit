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
var analysisdefinitionlog = logf.Log.WithName("analysisdefinition-resource")

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

	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *AnalysisDefinition) ValidateUpdate(old runtime.Object) error {
	analysisdefinitionlog.Info("validate update", "name", r.Name)

	for _, o := range r.Spec.Objectives {
		if err := o.validate(); err != nil {
			return err
		}
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *AnalysisDefinition) ValidateDelete() error {
	analysisdefinitionlog.Info("validate delete", "name", r.Name)

	return nil
}

func (o *Objective) validate() error {
	if err := o.SLOTargets.Pass.validate(); err != nil {
		return err
	}
	if o.SLOTargets.Warning != nil {
		if err := o.SLOTargets.Warning.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (c *CriteriaSet) validate() error {
	if len(c.AllOf) > 0 && len(c.AnyOf) > 0 {
		return fmt.Errorf("CriteriaSet: AllOf and Anyof are set simultaneusly")
	}

	if len(c.AllOf) == 0 && len(c.AnyOf) == 0 {
		return fmt.Errorf("CriteriaSet: AllOf nor Anyof set")
	}

	for _, a := range c.AllOf {
		if err := a.validate(); err != nil {
			return err
		}
	}

	for _, a := range c.AnyOf {
		if err := a.validate(); err != nil {
			return err
		}
	}

	return nil
}

func (c *Criteria) validate() error {
	if len(c.AllOf) > 0 && len(c.AnyOf) > 0 {
		return fmt.Errorf("Criteria: AllOf and Anyof are set simultaneusly")
	}

	if len(c.AllOf) == 0 && len(c.AnyOf) == 0 {
		return fmt.Errorf("Criteria: AllOf nor Anyof set")
	}

	for _, t := range c.AllOf {
		if err := t.validate(); err != nil {
			return err
		}
	}

	for _, t := range c.AnyOf {
		if err := t.validate(); err != nil {
			return err
		}
	}
	return nil
}

func (t *Target) validate() error {
	counter := 0
	if t.LessThan != nil {
		counter++
	}
	if t.LessThanOrEqual != nil {
		counter++
	}
	if t.GreaterThan != nil {
		counter++
	}
	if t.GreaterThanOrEqual != nil {
		counter++
	}
	if t.EqualTo != nil {
		counter++
	}
	if counter > 1 {
		return fmt.Errorf("Target: multiple targets set anot allowed per Analysis")
	}
	if counter == 0 {
		return fmt.Errorf("Target: not set")
	}
	return nil
}
