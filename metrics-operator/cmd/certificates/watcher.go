package certificates

import (
	"bytes"
	"context"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-logr/logr"
	"github.com/spf13/afero"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	certificateRenewalInterval = 6 * time.Hour
	ServerKey                  = "tls.key"
	ServerCert                 = "tls.crt"
)

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
		if err := watcher.updateCertificatesFromSecret(); err != nil {
			watcher.Log.Error(err, "failed to update certificates")
		} else {
			watcher.Log.Info("updated certificate successfully")
		}
	}
}

func (watcher *CertificateWatcher) updateCertificatesFromSecret() error {
	var secret corev1.Secret

	err := watcher.apiReader.Get(context.TODO(),
		client.ObjectKey{Name: watcher.certificateSecretName, Namespace: watcher.namespace}, &secret)
	if err != nil {
		return err
	}

	watcher.Log.Info("checking dir", "watcher.certificateDirectory ", watcher.certificateDirectory)
	if _, err = watcher.fs.Stat(watcher.certificateDirectory); os.IsNotExist(err) {
		err = watcher.fs.MkdirAll(watcher.certificateDirectory, 0755)
		if err != nil {
			return fmt.Errorf("could not create cert directory: %w", err)
		}
	}

	for _, filename := range []string{ServerCert, ServerKey} {
		if _, err = watcher.ensureCertificateFile(secret, filename); err != nil {
			return err
		}
	}
	isValid, err := watcher.ValidateCertificateExpiration(secret.Data[ServerCert], certificateRenewalInterval, time.Now())
	if err != nil {
		return err
	} else if !isValid {
		return fmt.Errorf("certificate is outdated")
	}
	return nil
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

		if err := watcher.updateCertificatesFromSecret(); err != nil {
			if k8serrors.IsNotFound(err) {
				watcher.Log.Info("waiting for certificate secret to be available.")
			} else {
				watcher.Log.Error(err, "failed to update certificates")
			}
			time.Sleep(10 * time.Second)
			continue
		}
		break
	}
	go watcher.watchForCertificatesSecret()
}

func (watcher *CertificateWatcher) ValidateCertificateExpiration(certData []byte, renewalThreshold time.Duration, now time.Time) (bool, error) {
	if block, _ := pem.Decode(certData); block == nil {
		watcher.Log.Error(errors.New("can't decode PEM file"), "failed to parse certificate")
		return false, nil
	} else if cert, err := x509.ParseCertificate(block.Bytes); err != nil {
		watcher.Log.Error(err, "failed to parse certificate")
		return false, err
	} else if now.After(cert.NotAfter.Add(-renewalThreshold)) {
		watcher.Log.Info("certificate is outdated, waiting for new ones", "Valid until", cert.NotAfter.UTC())
		return false, nil
	}
	return true, nil
}
