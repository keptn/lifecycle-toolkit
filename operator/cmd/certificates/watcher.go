package certificates

import (
	"bytes"
	"context"
	"fmt"

	"github.com/go-logr/logr"
	kubeobjects2 "keptn.sh/keptnwebhook/kubeobjects"

	"os"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

// TODO: refactor code below to be testable and also tested
const certificateRenewalInterval = 6 * time.Hour

type CertificateWatcher struct {
	apiReader             client.Reader
	fs                    afero.Fs
	certificateDirectory  string
	namespace             string
	certificateSecretName string
	Log                   logr.Logger
}

func NewCertificateWatcher(mgr manager.Manager, namespace string, secretName string, log logr.Logger) *CertificateWatcher {
	return &CertificateWatcher{
		apiReader:             mgr.GetAPIReader(),
		fs:                    afero.NewOsFs(),
		certificateDirectory:  mgr.GetWebhookServer().CertDir,
		namespace:             namespace,
		certificateSecretName: secretName,
		Log:                   log,
	}
}

func (watcher *CertificateWatcher) watchForCertificatesSecret() {
	for {
		<-time.After(certificateRenewalInterval)
		watcher.Log.Info("checking for new certificates")
		if updated, err := watcher.updateCertificatesFromSecret(); err != nil {
			watcher.Log.Info("failed to update certificates", "error", err)
		} else if updated {
			watcher.Log.Info("updated certificate successfully")
		}
	}
}

func (watcher *CertificateWatcher) updateCertificatesFromSecret() (bool, error) {
	var secret corev1.Secret

	err := watcher.apiReader.Get(context.TODO(),
		client.ObjectKey{Name: watcher.certificateSecretName, Namespace: watcher.namespace}, &secret)
	if err != nil {
		return false, err
	}

	if _, err = watcher.fs.Stat(watcher.certificateDirectory); os.IsNotExist(err) {
		err = watcher.fs.MkdirAll(watcher.certificateDirectory, 0755)
		if err != nil {
			return false, fmt.Errorf("could not create cert directory: %s", err)
		}
	}

	for _, filename := range []string{keptnwebhookcontroller.ServerCert, keptnwebhookcontroller.ServerKey} {
		if _, err = watcher.ensureCertificateFile(secret, filename); err != nil {
			return false, err
		}
	}
	isValid, err := kubeobjects2.ValidateCertificateExpiration(secret.Data[keptnwebhookcontroller.ServerCert], certificateRenewalInterval, time.Now())
	if err != nil {
		return false, err
	} else if !isValid {
		return false, fmt.Errorf("certificate is outdated")
	}
	return true, nil
}

func (watcher *CertificateWatcher) ensureCertificateFile(secret corev1.Secret, filename string) (bool, error) {
	f := filepath.Join(watcher.certificateDirectory, filename)

	data, err := afero.ReadFile(watcher.fs, f)
	if os.IsNotExist(err) || !bytes.Equal(data, secret.Data[filename]) {
		if err := afero.WriteFile(watcher.fs, f, secret.Data[filename], 0666); err != nil {
			return false, err
		}
	} else {
		return false, err
	}
	return true, nil
}

func (watcher *CertificateWatcher) WaitForCertificates() {
	for threshold := time.Now().Add(5 * time.Minute); time.Now().Before(threshold); {
		_, err := watcher.updateCertificatesFromSecret()

		if err != nil {
			if k8serrors.IsNotFound(err) {
				watcher.Log.Info("waiting for certificate secret to be available.")
			} else {
				watcher.Log.Info("failed to update certificates", "error", err)
			}
			time.Sleep(10 * time.Second)
			continue
		}
		break
	}
	go watcher.watchForCertificatesSecret()
}
