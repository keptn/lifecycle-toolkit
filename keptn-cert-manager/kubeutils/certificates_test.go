package kubeutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"testing"
	"time"
)

func TestValidateCertificateExpiration(t *testing.T) {
	certTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now(),
		NotAfter:     time.Now().Add(24 * time.Hour),
	}

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("Failed to generate private key: %v", err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, certTemplate, certTemplate, &privateKey.PublicKey, privateKey)
	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})

	testCases := []struct {
		name             string
		certData         []byte
		renewalThreshold time.Duration
		now              time.Time
		expectedValid    bool
		expectedErrorNil bool
	}{
		{
			name:             "Valid certificate",
			certData:         certPEM,
			renewalThreshold: 2 * time.Hour,
			now:              time.Now().Add(21 * time.Hour), // Certificate is still valids
			expectedValid:    true,
			expectedErrorNil: true,
		},
		{
			name:             "Expired certificate",
			certData:         certPEM,
			renewalThreshold: 2 * time.Hour,
			now:              time.Now().Add(25 * time.Hour), // Certificate has expired
			expectedValid:    false,
			expectedErrorNil: true,
		},
		{
			name:             "Invalid PEM data",
			certData:         []byte("invalid PEM data"),
			renewalThreshold: 2 * time.Hour,
			now:              time.Now(),
			expectedValid:    false,
			expectedErrorNil: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			valid, err := ValidateCertificateExpiration(tc.certData, tc.renewalThreshold, tc.now)
			if valid != tc.expectedValid {
				t.Errorf("Expected valid=%v, got %v", tc.expectedValid, valid)
			}
			if (err == nil) != tc.expectedErrorNil {
				t.Errorf("Expected error nil=%v, got error=%v", tc.expectedErrorNil, err)
			}
		})
	}
}
