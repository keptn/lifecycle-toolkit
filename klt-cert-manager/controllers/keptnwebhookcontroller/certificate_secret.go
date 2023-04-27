package keptnwebhookcontroller

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/go-logr/logr"
	"github.com/keptn/lifecycle-toolkit/klt-cert-manager/kubeutils"
	"github.com/keptn/lifecycle-toolkit/klt-cert-manager/pkg/common"
	"github.com/pkg/errors"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	apiv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type certificateSecret struct {
	secret          *corev1.Secret
	certificates    *Certs
	existsInCluster bool
	clt             client.Client
}

func newCertificateSecret(clt client.Client) *certificateSecret {
	return &certificateSecret{
		clt: clt,
	}
}

func (certSecret *certificateSecret) setSecretFromReader(ctx context.Context, namespace string, log logr.Logger) error {
	query := kubeutils.NewSecretQuery(nil, certSecret.clt, log)
	secret, err := query.Get(ctx, types.NamespacedName{Name: buildSecretName(), Namespace: namespace})

	if k8serrors.IsNotFound(err) {
		certSecret.secret = kubeutils.NewSecret(buildSecretName(), namespace, map[string][]byte{})
		certSecret.existsInCluster = false
	} else if err != nil {
		return errors.WithStack(err)
	} else {
		certSecret.secret = &secret
		certSecret.existsInCluster = true
	}

	return nil
}

func (certSecret *certificateSecret) isRecent() bool {
	if certSecret.secret == nil && certSecret.certificates == nil {
		return true
	} else if certSecret.secret == nil || certSecret.certificates == nil {
		return false
	} else if !reflect.DeepEqual(certSecret.certificates.Data, certSecret.secret.Data) {
		// certificates need to be updated
		return false
	}
	return true
}

func (certSecret *certificateSecret) setCertificates(namespace string) error {
	certs := Certs{
		Domain:  getDomain(namespace),
		SrcData: certSecret.secret.Data,
		Now:     time.Now(),
	}
	if err := certs.Validate(); err != nil {
		return errors.WithStack(err)
	}

	certSecret.certificates = &certs

	return nil
}

func buildSecretName() string {
	return common.SecretName
}

func getDomain(namespace string) string {
	// use a wildcard here because the secret is potentially used by multiple operators providing webhooks
	return fmt.Sprintf("*.%s.svc", namespace)
}

func (certSecret *certificateSecret) areWebhookConfigsValid(configs []*admissionregistrationv1.WebhookClientConfig) bool {
	for i := range configs {
		if configs[i] != nil && !certSecret.isBundleValid((*configs[i]).CABundle) {
			return false
		}
	}
	return true
}

func (certSecret *certificateSecret) isBundleValid(bundle []byte) bool {
	return len(bundle) != 0 && bytes.Equal(bundle, certSecret.certificates.Data[RootCert])
}

func (certSecret *certificateSecret) createOrUpdateIfNecessary(ctx context.Context) error {
	if certSecret.isRecent() && certSecret.existsInCluster {
		return nil
	}

	var err error

	certSecret.secret.Data = certSecret.certificates.Data
	if certSecret.existsInCluster {
		err = certSecret.clt.Update(ctx, certSecret.secret)
	} else {
		err = certSecret.clt.Create(ctx, certSecret.secret)
	}

	return errors.WithStack(err)
}

func (certSecret *certificateSecret) loadCombinedBundle() ([]byte, error) {
	data, hasData := certSecret.secret.Data[RootCert]
	if !hasData {
		return nil, errors.New(common.CertificatesSecretEmptyErr)
	}

	if oldData, hasOldData := certSecret.secret.Data[RootCertOld]; hasOldData {
		data = append(data, oldData...)
	}
	return data, nil
}

func (certSecret *certificateSecret) areCRDConversionsValid(crds *apiv1.CustomResourceDefinitionList) bool {
	for _, crd := range crds.Items {
		if !certSecret.isCRDConversionValid(crd.Spec.Conversion) {
			return false
		}
	}
	return true
}

func (certSecret *certificateSecret) isCRDConversionValid(conversion *apiv1.CustomResourceConversion) bool {
	if conversion.Strategy == apiv1.NoneConverter || conversion.Webhook == nil {
		return true
	}
	return certSecret.isBundleValid(conversion.Webhook.ClientConfig.CABundle)
}
