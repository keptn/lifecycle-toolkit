package keptnwebhookcontroller

import (
	"crypto/x509"
	"encoding/pem"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCertsValidation(t *testing.T) {
	now, _ := time.Parse(time.RFC3339, "2018-01-10T00:00:00Z")
	domain := "keptn-webhook.webhook.svc"
	firstCerts := Certs{
		Domain: domain,
		Now:    now,
	}

	require.NoError(t, firstCerts.Validate())
	require.Equal(t, len(firstCerts.Data), 5)
	requireValidCerts(t, domain, now.Add(5*time.Minute), firstCerts.Data[RootCert], firstCerts.Data[ServerCert])

	t.Run("up-to-date certs", func(t *testing.T) {
		newTime := now.Add(5 * time.Minute)

		newCerts := Certs{Domain: domain, SrcData: firstCerts.Data, Now: newTime}
		require.NoError(t, newCerts.Validate())
		requireValidCerts(t, domain, newTime, newCerts.Data[RootCert], newCerts.Data[ServerCert])

		// No changes should have been applied.
		assert.Equal(t, string(firstCerts.Data[RootCert]), string(newCerts.Data[RootCert]))
		assert.Equal(t, string(firstCerts.Data[RootCertOld]), "")
		assert.Equal(t, string(firstCerts.Data[RootKey]), string(newCerts.Data[RootKey]))
		assert.Equal(t, string(firstCerts.Data[ServerCert]), string(newCerts.Data[ServerCert]))
		assert.Equal(t, string(firstCerts.Data[ServerKey]), string(newCerts.Data[ServerKey]))
	})

	t.Run("outdated server certs", func(t *testing.T) {
		newTime := now.Add((6*24 + 22) * time.Hour) // 6d22h

		newCerts := Certs{Domain: domain, SrcData: firstCerts.Data, Now: newTime}
		require.NoError(t, newCerts.Validate())
		requireValidCerts(t, domain, newTime, newCerts.Data[RootCert], newCerts.Data[ServerCert])

		// Server certificates should have been updated.
		assert.Equal(t, string(firstCerts.Data[RootCert]), string(newCerts.Data[RootCert]))
		assert.Equal(t, string(firstCerts.Data[RootCertOld]), "")
		assert.Equal(t, string(firstCerts.Data[RootKey]), string(newCerts.Data[RootKey]))
		assert.NotEqual(t, string(firstCerts.Data[ServerCert]), string(newCerts.Data[ServerCert]))
		assert.NotEqual(t, string(firstCerts.Data[ServerKey]), string(newCerts.Data[ServerKey]))
	})

	t.Run("outdated root certs", func(t *testing.T) {
		newTime := now.Add((364*24 + 22) * time.Hour) // 364d22h

		newCerts := Certs{Domain: domain, SrcData: firstCerts.Data, Now: newTime}
		require.NoError(t, newCerts.Validate())
		requireValidCerts(t, domain, newTime, newCerts.Data[RootCert], newCerts.Data[ServerCert])

		// Server certificates should have been updated.
		assert.Equal(t, string(firstCerts.Data[RootCert]), string(newCerts.Data[RootCertOld]))
		assert.NotEqual(t, string(firstCerts.Data[RootCert]), string(newCerts.Data[RootCert]))
		assert.NotEqual(t, string(firstCerts.Data[RootKey]), string(newCerts.Data[RootKey]))
		assert.NotEqual(t, string(firstCerts.Data[ServerCert]), string(newCerts.Data[ServerCert]))
		assert.NotEqual(t, string(firstCerts.Data[ServerKey]), string(newCerts.Data[ServerKey]))
	})
}

func requireValidCerts(t *testing.T, domain string, now time.Time, caCert, tlsCert []byte) {
	caCerts := x509.NewCertPool()
	require.True(t, caCerts.AppendCertsFromPEM(caCert))

	block, _ := pem.Decode(tlsCert)
	require.NotNil(t, block)
	cert, err := x509.ParseCertificate(block.Bytes)
	require.NoError(t, err)

	_, err = cert.Verify(x509.VerifyOptions{DNSName: domain, CurrentTime: now, Roots: caCerts})
	require.NoError(t, err)
}
