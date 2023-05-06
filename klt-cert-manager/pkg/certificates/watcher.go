package certificates

import (
	"bytes"
	"context"
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
)

const (
	certificateRenewalInterval = 6 * time.Hour
	ServerKey                  = "tls.key"
	ServerCert                 = "tls.crt"
	CertThreshold              = 5 * time.Minute
)

//go:generate moq -pkg fake -skip-ensure -out ./fake/watcher_mock.go . ICertificateWatcher:MockCertificateWatcher
type ICertificateWatcher interface {
	WaitForCertificates()
}

type CertificateWatcher struct {
	apiReader             client.Reader
	fs                    afero.Fs
	certificateDirectory  string
	namespace             string
	certificateSecretName string
	certificateTreshold   time.Duration
	ICertificateHandler
	Log logr.Logger
}

func NewCertificateWatcher(reader client.Reader, certDir string, namespace string, secretName string, log logr.Logger) *CertificateWatcher {
	return &CertificateWatcher{
		apiReader:             reader,
		fs:                    afero.NewOsFs(),
		certificateDirectory:  certDir,
		namespace:             namespace,
		certificateSecretName: secretName,
		ICertificateHandler:   defaultCertificateHandler{},
		certificateTreshold:   CertThreshold,
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
		if err = watcher.ensureCertificateFile(secret, filename); err != nil {
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

func (watcher *CertificateWatcher) ensureCertificateFile(secret corev1.Secret, filename string) error {
	f := filepath.Join(watcher.certificateDirectory, filename)
	data, err := afero.ReadFile(watcher.fs, f)
	if os.IsNotExist(err) || !bytes.Equal(data, secret.Data[filename]) {
		return afero.WriteFile(watcher.fs, f, secret.Data[filename], 0666)
	}
	return err

}

func (watcher *CertificateWatcher) WaitForCertificates() {
	for threshold := time.Now().Add(watcher.certificateTreshold); time.Now().Before(threshold); {

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
	if block, _ := watcher.Decode(certData); block == nil {
		watcher.Log.Error(errors.New("can't decode PEM file"), "failed to parse certificate")
		return false, nil
	} else if cert, err := watcher.Parse(block.Bytes); err != nil {
		watcher.Log.Error(err, "failed to parse certificate")
		return false, err
	} else if now.After(cert.NotAfter.Add(-renewalThreshold)) {
		watcher.Log.Info("certificate is outdated, waiting for new ones", "Valid until", cert.NotAfter.UTC())
		return false, nil
	}
	return true, nil
}
