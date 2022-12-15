package keptnwebhookcontroller

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"

	"github.com/keptn/lifecycle-toolkit/klt-cert-manager/kubeutils"
)

var serialNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128)

const (
	renewalThreshold = 12 * time.Hour

	RootKey     = "ca.key"
	RootCert    = "ca.crt"
	RootCertOld = "ca.crt.old"
	ServerKey   = "tls.key"
	ServerCert  = "tls.crt"
)

// Certs handles creation and renewal of CA and SSL/TLS server certificates.
type Certs struct {
	Domain  string
	SrcData map[string][]byte
	Data    map[string][]byte

	Now time.Time

	rootPrivateKey *ecdsa.PrivateKey
	rootPublicCert *x509.Certificate
}

// ValidateCerts checks for certificates and keys on cs.SrcData and renews them if needed. The existing (or new)
// certificates will be stored on cs.Data.
func (cs *Certs) ValidateCerts() error {
	cs.Data = map[string][]byte{}
	if cs.SrcData != nil {
		for k, v := range cs.SrcData {
			cs.Data[k] = v
		}
	}

	now := cs.Now
	if now.IsZero() {
		now = time.Now().UTC()
	}

	renewRootCerts := cs.validateRootCerts(now)
	if renewRootCerts {
		if err := cs.generateRootCerts(cs.Domain, now); err != nil {
			return err
		}
	}

	if renewRootCerts || cs.validateServerCerts(now) {
		return cs.generateServerCerts(cs.Domain, now)
	}

	return nil
}

func (cs *Certs) validateRootCerts(now time.Time) bool {
	if cs.Data[RootKey] == nil || cs.Data[RootCert] == nil {
		////log.Info("no root certificates found, creating")
		return true
	}

	var err error

	if block, _ := pem.Decode(cs.Data[RootCert]); block == nil {
		//log.Info("failed to parse root certificates, renewing", "error", "can't decode PEM file")
		return true
	} else if cs.rootPublicCert, err = x509.ParseCertificate(block.Bytes); err != nil {
		//log.Info("failed to parse root certificates, renewing", "error", err)
		return true
	} else if now.After(cs.rootPublicCert.NotAfter.Add(-renewalThreshold)) {
		//log.Info("root certificates are about to expire, renewing", "current", now, "expiration", cs.rootPublicCert.NotAfter)
		return true
	}

	if block, _ := pem.Decode(cs.Data[RootKey]); block == nil {
		//log.Info("failed to parse root key, renewing", "error", "can't decode PEM file")
		return true
	} else if cs.rootPrivateKey, err = x509.ParseECPrivateKey(block.Bytes); err != nil {
		//log.Info("failed to parse root key, renewing", "error", err)
		return true
	}

	return false
}

func (cs *Certs) validateServerCerts(now time.Time) bool {
	if cs.Data[ServerKey] == nil || cs.Data[ServerCert] == nil {
		//log.Info("no server certificates found, creating")
		return true
	}

	isValid, err := kubeutils.ValidateCertificateExpiration(cs.Data[ServerCert], renewalThreshold, now)
	if err != nil || !isValid {
		//log.Info("server certificate failed to parse or is outdated")
		return true
	}
	return false
}

func (cs *Certs) generateRootCerts(domain string, now time.Time) error {
	// Generate CA root keys
	//log.Info("generating root certificate")
	privateKey, err := cs.generatePrivateKey(RootKey)
	if err != nil {
		return err
	}
	cs.rootPrivateKey = privateKey

	// Generate CA root certificate
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number for root certificate: %w", err)
	}

	cs.rootPublicCert = &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:            []string{"AT"},
			Province:           []string{"KL"},
			Locality:           []string{"Klagenfurt"},
			Organization:       []string{"Keptn"},
			OrganizationalUnit: []string{"LifecycleToolkit"},
			CommonName:         domain,
		},
		IsCA: true,

		NotBefore: now,
		NotAfter:  now.Add(365 * 24 * time.Hour),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	rootPublicCertDER, err := x509.CreateCertificate(
		rand.Reader,
		cs.rootPublicCert,
		cs.rootPublicCert,
		cs.rootPrivateKey.Public(),
		cs.rootPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to generate root certificate: %w", err)
	}

	cs.Data[RootCertOld] = cs.Data[RootCert]
	cs.Data[RootCert] = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: rootPublicCertDER})

	//log.Info("root certificate generated")
	return nil
}

func (cs *Certs) generateServerCerts(domain string, now time.Time) error {
	// Generate server keys
	//log.Info("generating server certificate")
	privateKey, err := cs.generatePrivateKey(ServerKey)
	if err != nil {
		return err
	}

	// Generate server certificate
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return fmt.Errorf("failed to generate serial number for server certificate: %w", err)
	}

	tpl := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Country:            []string{"AT"},
			Province:           []string{"KL"},
			Locality:           []string{"Klagenfurt"},
			Organization:       []string{"Keptn"},
			OrganizationalUnit: []string{"LifecycleToolkit"},
			CommonName:         domain,
		},

		DNSNames: []string{domain},

		NotBefore: now,
		NotAfter:  now.Add(7 * 24 * time.Hour),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	serverPublicCertDER, err := x509.CreateCertificate(rand.Reader, tpl, cs.rootPublicCert, privateKey.Public(), cs.rootPrivateKey)
	if err != nil {
		return fmt.Errorf("failed to generate server certificate: %w", err)
	}

	cs.Data[ServerCert] = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: serverPublicCertDER})
	//log.Info("server certificate generated")
	return nil
}

func (cs *Certs) generatePrivateKey(dataKey string) (*ecdsa.PrivateKey, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate server private key: %w", err)
	}

	x509Encoded, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	cs.Data[dataKey] = pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509Encoded,
	})
	return privateKey, nil
}
