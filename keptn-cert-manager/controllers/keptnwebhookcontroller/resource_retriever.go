package keptnwebhookcontroller

import (
	"context"

	"github.com/go-logr/logr"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	apiv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type IResourceRetriever interface {
	GetMutatingWebhooks(ctx context.Context) (*admissionregistrationv1.MutatingWebhookConfigurationList, error)
	GetValidatingWebhooks(ctx context.Context) (*admissionregistrationv1.ValidatingWebhookConfigurationList, error)
	GetCRDs(ctx context.Context) (*apiv1.CustomResourceDefinitionList, error)
}

func NewResourceRetriever(config CertificateReconcilerConfig) IResourceRetriever {
	if config.WatchResources != nil {
		return &ResourceNameRetriever{
			Client:         config.Client,
			WatchResources: *config.WatchResources,
			Log:            config.Log,
		}
	}
	return &LabelSelectorRetriever{
		MatchLabels: config.MatchLabels,
		Client:      config.Client,
	}
}

type LabelSelectorRetriever struct {
	MatchLabels labels.Set
	Client      client.Client
}

func (r LabelSelectorRetriever) GetMutatingWebhooks(ctx context.Context) (*admissionregistrationv1.MutatingWebhookConfigurationList, error) {
	var mutatingWebhooks admissionregistrationv1.MutatingWebhookConfigurationList

	if err := r.Client.List(ctx, &mutatingWebhooks, client.MatchingLabels(r.MatchLabels)); err != nil {
		return nil, err
	}
	return &mutatingWebhooks, nil
}

func (r LabelSelectorRetriever) GetValidatingWebhooks(ctx context.Context) (*admissionregistrationv1.ValidatingWebhookConfigurationList, error) {
	var validatingWebhooks admissionregistrationv1.ValidatingWebhookConfigurationList

	if err := r.Client.List(ctx, &validatingWebhooks, client.MatchingLabels(r.MatchLabels)); err != nil {
		return nil, err
	}
	return &validatingWebhooks, nil
}

func (r LabelSelectorRetriever) GetCRDs(ctx context.Context) (*apiv1.CustomResourceDefinitionList, error) {
	var crds apiv1.CustomResourceDefinitionList
	opt := client.MatchingLabels(r.MatchLabels)
	if err := r.Client.List(ctx, &crds, opt); err != nil {
		return nil, err
	}

	return &crds, nil
}

type ResourceNameRetriever struct {
	Client         client.Client
	WatchResources ObservedObjects
	Log            logr.Logger
}

func (r ResourceNameRetriever) GetMutatingWebhooks(ctx context.Context) (*admissionregistrationv1.MutatingWebhookConfigurationList, error) {
	result := &admissionregistrationv1.MutatingWebhookConfigurationList{
		Items: []admissionregistrationv1.MutatingWebhookConfiguration{},
	}
	for _, mwhName := range r.WatchResources.MutatingWebhooks {
		mwh := &admissionregistrationv1.MutatingWebhookConfiguration{}
		if err := r.Client.Get(ctx, client.ObjectKey{Name: mwhName}, mwh); err != nil {
			r.Log.Error(err, "Could not retrieve MutatingWebhookConfiguration", "MutatingWebhookConfiguration", mwhName)
		} else {
			result.Items = append(result.Items, *mwh)
		}
	}

	return result, nil
}

func (r ResourceNameRetriever) GetValidatingWebhooks(ctx context.Context) (*admissionregistrationv1.ValidatingWebhookConfigurationList, error) {
	result := &admissionregistrationv1.ValidatingWebhookConfigurationList{
		Items: []admissionregistrationv1.ValidatingWebhookConfiguration{},
	}
	for _, vwhName := range r.WatchResources.ValidatingWebhooks {
		vwh := &admissionregistrationv1.ValidatingWebhookConfiguration{}
		if err := r.Client.Get(ctx, client.ObjectKey{Name: vwhName}, vwh); err != nil {
			r.Log.Error(err, "Could not retrieve ValidatingWebhookConfiguration", "ValidatingWebhookConfiguration", vwhName)
		} else {
			result.Items = append(result.Items, *vwh)
		}
	}

	return result, nil
}

func (r ResourceNameRetriever) GetCRDs(ctx context.Context) (*apiv1.CustomResourceDefinitionList, error) {
	result := &apiv1.CustomResourceDefinitionList{
		Items: []apiv1.CustomResourceDefinition{},
	}
	for _, crdName := range r.WatchResources.CustomResourceDefinitions {
		crd := &apiv1.CustomResourceDefinition{}
		if err := r.Client.Get(ctx, client.ObjectKey{Name: crdName}, crd); err != nil {
			r.Log.Error(err, "Could not retrieve ValidatingWebhookConfiguration", "ValidatingWebhookConfiguration", crdName)
		} else {
			result.Items = append(result.Items, *crd)
		}
	}

	return result, nil
}
