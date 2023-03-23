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
	"context"
	"errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var logger = logf.Log.WithName("keptnconfig-resource")

// client to fetch other KeptnConfig from the APIs
var _client client.Reader = nil
var _ns string

func (r *KeptnConfig) SetupWebhookWithManager(mgr ctrl.Manager, namespace string) error {
	_client = mgr.GetAPIReader()
	_ns = namespace
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/validate-options-keptn-sh-v1alpha1-keptnconfig,mutating=false,failurePolicy=fail,sideEffects=None,groups=options.keptn.sh,resources=keptnconfigs,verbs=create;update,versions=v1alpha1,name=vkeptnconfig.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &KeptnConfig{}

// ValidateCreate checks that there is not yet another KetpnConfig active
func (r *KeptnConfig) ValidateCreate() error {
	logger.Info("Validating KeptnConfig", "name", r.Name)
	configs := &KeptnConfigList{}
	if err := _client.List(context.TODO(), configs, client.InNamespace(_ns)); err != nil {
		logger.Error(err, "Impossible collecting all KeptnConfig")
		return err
	}
	if len(configs.Items) > 0 {
		return errors.New("only a single KeptnConfig can be applied")
	}
	return nil
}

// ValidateUpdate immediately returns since there is nothing to validate
func (r *KeptnConfig) ValidateUpdate(old runtime.Object) error {
	return nil
}

// ValidateDelete immediately returns since there is nothing to validate
func (r *KeptnConfig) ValidateDelete() error {
	return nil
}
