package certificates

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-logr/logr/testr"
	"github.com/keptn/lifecycle-toolkit/keptn-cert-manager/pkg/certificates/fake"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	fakeClient "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

const CACERT = `-----BEGIN CERTIFICATE-----
MIIBrzCCAVmgAwIBAgIUH/zWlPkTXVBcu2zOvUy/NV1hCKkwDQYJKoZIhvcNAQEL
BQAwRTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDAeFw0yNDA0MTgwOTEzMDdaFw0zNDA0
MTYwOTEzMDdaMEUxCzAJBgNVBAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEw
HwYDVQQKDBhJbnRlcm5ldCBXaWRnaXRzIFB0eSBMdGQwXDANBgkqhkiG9w0BAQEF
AANLADBIAkEAyLJjXFVA0DzUVSJy+ANqe+tXki2MsWgm+cbYkpBMLJMKhhwnv6vW
Hxwsh5MZNwAmSoprINGb7i6Ub2OhjpVq0QIDAQABoyEwHzAdBgNVHQ4EFgQUtwGr
j5axZSNJo6o1mP7L09axxIIwDQYJKoZIhvcNAQELBQADQQDIJGtVIgsg0J3e5QRf
LZ21sKKY+xzeG5yy90ao8QMWX9CqCpZncprE1MJijkG7paCFq6Bh22g6xTZYYJ1m
yG/y
-----END CERTIFICATE-----`

const CAKEY = `-----BEGIN PRIVATE KEY-----
MIIBVAIBADANBgkqhkiG9w0BAQEFAASCAT4wggE6AgEAAkEAyLJjXFVA0DzUVSJy
+ANqe+tXki2MsWgm+cbYkpBMLJMKhhwnv6vWHxwsh5MZNwAmSoprINGb7i6Ub2Oh
jpVq0QIDAQABAkANdxJ9hmbD0eD5GUeXZjtFtyN39kBjQraiuXmcU7wYnWJ9OyaB
jsKkWlv9vx1stbMSYzlSQepDRYVcKL6AgGexAiEA7EwLkpiWT41/IwIIoYVQNgMN
Q/n8ltO47ecFljF1G6UCIQDZbm1JXYF068xo0vglnKl9HK3I69cHA4hrVww0ZUha
vQIgIDy7s3NHxnCqcK89WDPk3omKDMUVNcqKx0ImW/hBXtUCIQCvrMgCCdmp9UaP
vz0dbomGe6ByARMYKKOVTpyezOJ75QIgNqihb0lQbzEceTo6S2bQakDH7dH4Eydd
hMfh5Ml1u3o=
-----END PRIVATE KEY-----`

const OUTDATED_CACERT = `-----BEGIN CERTIFICATE-----
MIICPTCCAeKgAwIBAgIRAMIV/0UqFGHgKSYOWBdx/KcwCgYIKoZIzj0EAwIwczEL
MAkGA1UEBhMCQVQxCzAJBgNVBAgTAktMMRMwEQYDVQQHEwpLbGFnZW5mdXJ0MQ4w
DAYDVQQKEwVLZXB0bjEZMBcGA1UECxMQTGlmZWN5Y2xlVG9vbGtpdDEXMBUGA1UE
AwwOKi5rZXB0bi1ucy5zdmMwHhcNMjMwNDE5MTEwNDUzWhcNMjQwNDE4MTEwNDUz
WjBzMQswCQYDVQQGEwJBVDELMAkGA1UECBMCS0wxEzARBgNVBAcTCktsYWdlbmZ1
cnQxDjAMBgNVBAoTBUtlcHRuMRkwFwYDVQQLExBMaWZlY3ljbGVUb29sa2l0MRcw
FQYDVQQDDA4qLmtlcHRuLW5zLnN2YzBZMBMGByqGSM49AgEGCCqGSM49AwEHA0IA
BPxAP4JTJfwKz/P32dXuyfVi7kinQPebSYwF/gRAUcN0dCAi6GnxbI2OXlcU0guD
zHXv3VRh3EX2fiNszcfKaCajVzBVMA4GA1UdDwEB/wQEAwICpDATBgNVHSUEDDAK
BggrBgEFBQcDATAPBgNVHRMBAf8EBTADAQH/MB0GA1UdDgQWBBQUGe/8XYV1HsZs
nWsyrOCCGr/sQDAKBggqhkjOPQQDAgNJADBGAiEAkcPaCANDXW5Uillrof0VrnPw
ow49D22Gsrh7YM+vmTQCIQDU1L5IT0Zz+bdIyFSsDnEUXZDeydNv56DoSLh+358Y
aw==
-----END CERTIFICATE-----`

const OUTDATED_CAKEY = `-----BEGIN PRIVATE KEY-----
MHcCAQEEII5SAqBxINKatksyu2mTvLZZhfEOpNinYJDwlQjkfreboAoGCCqGSM49
AwEHoUQDQgAE/EA/glMl/ArP8/fZ1e7J9WLuSKdA95tJjAX+BEBRw3R0ICLoafFs
jY5eVxTSC4PMde/dVGHcRfZ+I2zNx8poJg==
-----END PRIVATE KEY-----`

const uniqueIDPEM = `-----BEGIN CERTIFICATE-----
MIID2jCCA0MCAg39MA0GCSqGSIb3DQEBBQUAMIGbMQswCQYDVQQGEwJKUDEOMAwG
A1UECBMFVG9reW8xEDAOBgNVBAcTB0NodW8ta3UxETAPBgNVBAoTCEZyYW5rNERE
MRgwFgYDVQQLEw9XZWJDZXJ0IFN1cHBvcnQxGDAWBgNVBAMTD0ZyYW5rNEREIFdl
YiBDQTEjMCEGCSqGSIb3DQEJARYUc3VwcG9ydEBmcmFuazRkZC5jb20wHhcNMTIw
ODIyMDUyODAwWhcNMTcwODIxMDUyODAwWjBKMQswCQYDVQQGEwJKUDEOMAwGA1UE
CAwFVG9reW8xETAPBgNVBAoMCEZyYW5rNEREMRgwFgYDVQQDDA93d3cuZXhhbXBs
ZS5jb20wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQCwvWITOLeyTbS1
Q/UacqeILIK16UHLvSymIlbbiT7mpD4SMwB343xpIlXN64fC0Y1ylT6LLeX4St7A
cJrGIV3AMmJcsDsNzgo577LqtNvnOkLH0GojisFEKQiREX6gOgq9tWSqwaENccTE
sAXuV6AQ1ST+G16s00iN92hjX9V/V66snRwTsJ/p4WRpLSdAj4272hiM19qIg9zr
h92e2rQy7E/UShW4gpOrhg2f6fcCBm+aXIga+qxaSLchcDUvPXrpIxTd/OWQ23Qh
vIEzkGbPlBA8J7Nw9KCyaxbYMBFb1i0lBjwKLjmcoihiI7PVthAOu/B71D2hKcFj
Kpfv4D1Uam/0VumKwhwuhZVNjLq1BR1FKRJ1CioLG4wCTr0LVgtvvUyhFrS+3PdU
R0T5HlAQWPMyQDHgCpbOHW0wc0hbuNeO/lS82LjieGNFxKmMBFF9lsN2zsA6Qw32
Xkb2/EFltXCtpuOwVztdk4MDrnaDXy9zMZuqFHpv5lWTbDVwDdyEQNclYlbAEbDe
vEQo/rAOZFl94Mu63rAgLiPeZN4IdS/48or5KaQaCOe0DuAb4GWNIQ42cYQ5TsEH
Wt+FIOAMSpf9hNPjDeu1uff40DOtsiyGeX9NViqKtttaHpvd7rb2zsasbcAGUl+f
NQJj4qImPSB9ThqZqPTukEcM/NtbeQIDAQABMA0GCSqGSIb3DQEBBQUAA4GBAIAi
gU3My8kYYniDuKEXSJmbVB+K1upHxWDA8R6KMZGXfbe5BRd8s40cY6JBYL52Tgqd
l8z5Ek8dC4NNpfpcZc/teT1WqiO2wnpGHjgMDuDL1mxCZNL422jHpiPWkWp3AuDI
c7tL1QjbfAUHAQYwmHkWgPP+T2wAv0pOt36GgMCM
-----END CERTIFICATE-----`

var ERR_BAD_CERT = errors.New("bad cert")

var emptySecret = v1.Secret{
	TypeMeta: metav1.TypeMeta{},
	ObjectMeta: metav1.ObjectMeta{
		Namespace: "default",
		Name:      "my-cert",
	},
}

var outdatedSecret = v1.Secret{
	TypeMeta: metav1.TypeMeta{},
	ObjectMeta: metav1.ObjectMeta{
		Namespace: "default",
		Name:      "my-cert",
	},
	Data: map[string][]byte{
		ServerCert: []byte(OUTDATED_CACERT),
		ServerKey:  []byte(OUTDATED_CAKEY),
	},
}

var goodSecret = v1.Secret{
	TypeMeta: metav1.TypeMeta{},
	ObjectMeta: metav1.ObjectMeta{
		Namespace: "default",
		Name:      "my-cert",
	},
	Data: map[string][]byte{
		ServerCert: []byte(CACERT),
		ServerKey:  []byte(CAKEY),
	},
}

func TestCertificateWatcher_ValidateCertificateExpiration(t *testing.T) {

	tests := []struct {
		name             string
		certHandler      ICertificateHandler
		certData         []byte
		renewalThreshold time.Duration
		now              time.Time
		want             bool
		wantErr          error
	}{
		{
			name: "certificate cannot be decoded",
			certHandler: &fake.ICertificateHandlerMock{
				DecodeFunc: func(data []byte) (p *pem.Block, rest []byte) {
					return nil, nil // fake a failure in the decoding
				},
				ParseFunc: nil,
			},
			want: false,
		},
		{
			name: "certificate cannot be parsed",
			certHandler: &fake.ICertificateHandlerMock{
				DecodeFunc: func(data []byte) (p *pem.Block, rest []byte) {
					return &pem.Block{Type: "test", Bytes: []byte("testdata")}, nil
				},
				ParseFunc: func(der []byte) (*x509.Certificate, error) {
					return nil, ERR_BAD_CERT
				},
			},
			want:    false,
			wantErr: ERR_BAD_CERT,
		},
		{
			name:        "good certificate - unexpired",
			certData:    []byte(uniqueIDPEM),
			certHandler: defaultCertificateHandler{},
			want:        true,
		},
		{
			name:        "good certificate - expired",
			certData:    []byte(uniqueIDPEM),
			now:         time.Now(), // setting up now makes sure that the threshold is passed
			certHandler: defaultCertificateHandler{},
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			watcher := &CertificateWatcher{
				ICertificateHandler: tt.certHandler,
				Log:                 testr.New(t),
			}
			got, err := watcher.ValidateCertificateExpiration(tt.certData, tt.renewalThreshold, tt.now)
			if tt.wantErr != nil {
				require.Error(t, err)
				t.Log("want:", tt.wantErr, "got:", err)
				require.True(t, errors.Is(tt.wantErr, err))
			}
			require.Equal(t, got, tt.want)
		})
	}
}

func TestCertificateWatcher_ensureCertificateFile(t *testing.T) {

	certdir := t.TempDir()
	f := filepath.Join(certdir, ServerCert)
	err := os.WriteFile(f, goodSecret.Data[ServerCert], 0666)
	require.Nil(t, err)
	baddir := t.TempDir()
	f = filepath.Join(baddir, ServerCert)
	err = os.WriteFile(f, goodSecret.Data[ServerKey], 0666)
	require.Nil(t, err)
	tests := []struct {
		name     string
		fs       afero.Fs
		secret   v1.Secret
		filename string
		certDir  string
		wantErr  bool
		err      string
	}{
		{
			name:     "if good cert exist in fs no error",
			secret:   goodSecret,
			certDir:  certdir,
			filename: ServerCert,
			wantErr:  false,
		},

		{
			name:     "if unexisting file name, we expect a file system error",
			secret:   emptySecret,
			filename: "$%&/())=$§%/=",
			certDir:  baddir,
			wantErr:  true,
			err:      "no such file or directory",
		},

		{
			name:     "wrong file content is replaced with updated cert",
			certDir:  baddir,
			secret:   goodSecret,
			filename: ServerCert,
			wantErr:  false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			watcher := &CertificateWatcher{
				fs:                   afero.NewOsFs(),
				certificateDirectory: tt.certDir,
			}
			err := watcher.ensureCertificateFile(tt.secret, tt.filename)
			if !tt.wantErr {
				require.Nil(t, err)
				f = filepath.Join(tt.certDir, ServerCert)
				data, err := os.ReadFile(f)
				if err != nil {
					panic(err)
				}
				if !bytes.Equal(data, tt.secret.Data[tt.filename]) {
					t.Errorf("ensureCertificateFile()data %v was not replaced with %v", data, tt.secret.Data[tt.filename])
				}
			} else {
				require.Contains(t, err.Error(), tt.err)
			}
		})
	}
}

func TestCertificateWatcher_updateCertificatesFromSecret(t *testing.T) {

	oldDir := t.TempDir()
	os.Remove(oldDir)

	tests := []struct {
		name                  string
		apiReader             client.Reader
		certificateDirectory  string
		namespace             string
		certificateSecretName string
		wantErr               error
	}{
		{
			name:                  "certificate not found",
			apiReader:             newFakeClient(),
			certificateDirectory:  t.TempDir(),
			namespace:             "default",
			certificateSecretName: "my-cert",
			wantErr:               errors.New("secrets \"my-cert\" not found"),
		},
		{
			name:                  "outdated certificate found, nothing in dir",
			apiReader:             newFakeClient(&outdatedSecret),
			certificateDirectory:  t.TempDir(),
			namespace:             "default",
			certificateSecretName: "my-cert",
			wantErr:               errors.New("certificate is outdated"),
		},

		{
			name:                  "outdated certificate found, not existing in dir",
			apiReader:             newFakeClient(&emptySecret),
			certificateDirectory:  oldDir,
			namespace:             "default",
			certificateSecretName: "my-cert",
			wantErr:               errors.New("certificate is outdated"),
		},
		{
			name:                  "good certificate - not stored",
			apiReader:             newFakeClient(&goodSecret),
			certificateDirectory:  t.TempDir(),
			namespace:             "default",
			certificateSecretName: "my-cert",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			watcher := &CertificateWatcher{
				apiReader:             tt.apiReader,
				fs:                    afero.NewOsFs(),
				certificateDirectory:  tt.certificateDirectory,
				namespace:             tt.namespace,
				certificateSecretName: tt.certificateSecretName,
				ICertificateHandler:   defaultCertificateHandler{},
				Log:                   testr.New(t),
			}
			err := watcher.updateCertificatesFromSecret()
			if tt.wantErr == nil {
				require.Nil(t, err)
			} else {
				require.NotNil(t, err)
				require.Contains(t, err.Error(), tt.wantErr.Error())
			}
		})
	}
}

func newFakeClient(objs ...client.Object) client.Reader {
	return fakeClient.NewClientBuilder().WithObjects(objs...).Build()
}

func TestNewCertificateWatcher(t *testing.T) {
	logger := testr.New(t)
	fclient := newFakeClient()
	want := &CertificateWatcher{
		apiReader:             fclient,
		fs:                    afero.NewOsFs(),
		namespace:             "default",
		certificateSecretName: "my-secret",
		certificateDirectory:  "test",
		certificateThreshold:  CertThreshold,
		ICertificateHandler:   defaultCertificateHandler{},
		Log:                   testr.New(t),
	}
	got := NewCertificateWatcher(fclient, "test", "default", "my-secret", logger)
	require.EqualValues(t, got, want)

}

func TestNewNoOpCertificateWatcher(t *testing.T) {
	require.EqualValues(t, NewNoOpCertificateWatcher(), &NoOpCertificateWatcher{})
}

func TestCertificateWatcher_watchForCertificatesSecret(t *testing.T) {
	mockReader := newFakeClient()
	logger := testr.New(t)

	watcher := &CertificateWatcher{
		apiReader:             mockReader,
		fs:                    afero.NewOsFs(),
		certificateDirectory:  "",
		namespace:             "",
		certificateSecretName: "",
		ICertificateHandler:   defaultCertificateHandler{},
		Log:                   logger,
	}

	updateSignal := make(chan struct{})
	defer close(updateSignal)

	updateCalled := false
	go func() {
		time.Sleep(10 * time.Millisecond)
		updateSignal <- struct{}{}
	}()

	go watcher.watchForCertificatesSecret()

	select {
	case <-time.After(20 * time.Millisecond):
		t.Error("Expected update but did not receive it within the specified interval")
	case <-updateSignal:
		updateCalled = true
	}

	if !updateCalled {
		t.Error("updateCertificatesFromSecret method was not called as expected")
	}
}
