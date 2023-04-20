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
	"github.com/keptn/lifecycle-toolkit/operator/cmd/certificates/fake"
	fakeclient "github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const CACERT = `-----BEGIN CERTIFICATE-----
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

const CAKEY = `-----BEGIN PRIVATE KEY-----
MHcCAQEEII5SAqBxINKatksyu2mTvLZZhfEOpNinYJDwlQjkfreboAoGCCqGSM49
AwEHoUQDQgAE/EA/glMl/ArP8/fZ1e7J9WLuSKdA95tJjAX+BEBRw3R0ICLoafFs
jY5eVxTSC4PMde/dVGHcRfZ+I2zNx8poJg==
-----END PRIVATE KEY-----`

const uniqueIDPEM = `-----BEGIN CERTIFICATE-----
MIIFsDCCBJigAwIBAgIIrOyC1ydafZMwDQYJKoZIhvcNAQEFBQAwgY4xgYswgYgG
A1UEAx6BgABNAGkAYwByAG8AcwBvAGYAdAAgAEYAbwByAGUAZgByAG8AbgB0ACAA
VABNAEcAIABIAFQAVABQAFMAIABJAG4AcwBwAGUAYwB0AGkAbwBuACAAQwBlAHIA
dABpAGYAaQBjAGEAdABpAG8AbgAgAEEAdQB0AGgAbwByAGkAdAB5MB4XDTE0MDEx
ODAwNDEwMFoXDTE1MTExNTA5Mzc1NlowgZYxCzAJBgNVBAYTAklEMRAwDgYDVQQI
EwdqYWthcnRhMRIwEAYDVQQHEwlJbmRvbmVzaWExHDAaBgNVBAoTE3N0aG9ub3Jl
aG90ZWxyZXNvcnQxHDAaBgNVBAsTE3N0aG9ub3JlaG90ZWxyZXNvcnQxJTAjBgNV
BAMTHG1haWwuc3Rob25vcmVob3RlbHJlc29ydC5jb20wggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQCvuu0qpI+Ko2X84Twkf84cRD/rgp6vpgc5Ebejx/D4
PEVON5edZkazrMGocK/oQqIlRxx/lefponN/chlGcllcVVPWTuFjs8k+Aat6T1qp
4iXxZekAqX+U4XZMIGJD3PckPL6G2RQSlF7/LhGCsRNRdKpMWSTbou2Ma39g52Kf
gsl3SK/GwLiWpxpcSkNQD1hugguEIsQYLxbeNwpcheXZtxbBGguPzQ7rH8c5vuKU
BkMOzaiNKLzHbBdFSrua8KWwCJg76Vdq/q36O9GlW6YgG3i+A4pCJjXWerI1lWwX
Ktk5V+SvUHGey1bkDuZKJ6myMk2pGrrPWCT7jP7WskChAgMBAAGBCQBCr1dgEleo
cKOCAfswggH3MIHDBgNVHREEgbswgbiCHG1haWwuc3Rob25vcmVob3RlbHJlc29y
dC5jb22CIGFzaGNoc3ZyLnN0aG9ub3JlaG90ZWxyZXNvcnQuY29tgiRBdXRvRGlz
Y292ZXIuc3Rob25vcmVob3RlbHJlc29ydC5jb22CHEF1dG9EaXNjb3Zlci5ob3Rl
bHJlc29ydC5jb22CCEFTSENIU1ZSghdzdGhvbm9yZWhvdGVscmVzb3J0LmNvbYIP
aG90ZWxyZXNvcnQuY29tMCEGCSsGAQQBgjcUAgQUHhIAVwBlAGIAUwBlAHIAdgBl
AHIwHQYDVR0OBBYEFMAC3UR4FwAdGekbhMgnd6lMejtbMAsGA1UdDwQEAwIFoDAT
BgNVHSUEDDAKBggrBgEFBQcDATAJBgNVHRMEAjAAMIG/BgNVHQEEgbcwgbSAFGfF
6xihk+gJJ5TfwvtWe1UFnHLQoYGRMIGOMYGLMIGIBgNVBAMegYAATQBpAGMAcgBv
AHMAbwBmAHQAIABGAG8AcgBlAGYAcgBvAG4AdAAgAFQATQBHACAASABUAFQAUABT
ACAASQBuAHMAcABlAGMAdABpAG8AbgAgAEMAZQByAHQAaQBmAGkAYwBhAHQAaQBv
AG4AIABBAHUAdABoAG8AcgBpAHQAeYIIcKhXEmBXr0IwDQYJKoZIhvcNAQEFBQAD
ggEBABlSxyCMr3+ANr+WmPSjyN5YCJBgnS0IFCwJAzIYP87bcTye/U8eQ2+E6PqG
Q7Huj7nfHEw9qnGo+HNyPp1ad3KORzXDb54c6xEoi+DeuPzYHPbn4c3hlH49I0aQ
eWW2w4RslSWpLvO6Y7Lboyz2/Thk/s2kd4RHxkkWpH2ltPqJuYYg3X6oM5+gIFHJ
WGnh+ojZ5clKvS5yXh3Wkj78M6sb32KfcBk0Hx6NkCYPt60ODYmWtvqwtw6r73u5
TnTYWRNvo2svX69TriL+CkHY9O1Hkwf2It5zHl3gNiKTJVaak8AuEz/CKWZneovt
yYLwhUhg3PX5Co1VKYE+9TxloiE=
-----END CERTIFICATE-----`

var ERR_BAD_CERT = errors.New("bad cert")

var emptySecret = v1.Secret{
	TypeMeta: metav1.TypeMeta{},
	ObjectMeta: metav1.ObjectMeta{
		Namespace: "default",
		Name:      "my-cert",
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
					return nil, nil //fake a failure in the decoding
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
			now:         time.Now(), //setting up now makes sure that the threshold is passed
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
			apiReader:             fakeclient.NewClient(),
			certificateDirectory:  t.TempDir(),
			namespace:             "default",
			certificateSecretName: "my-cert",
			wantErr:               errors.New("secrets \"my-cert\" not found"),
		},
		{
			name:                  "outdated certificate found, nothing in dir",
			apiReader:             fakeclient.NewClient(&emptySecret),
			certificateDirectory:  t.TempDir(),
			namespace:             "default",
			certificateSecretName: "my-cert",
			wantErr:               errors.New("certificate is outdated"),
		},

		{
			name:                  "outdated certificate found, not existing in dir",
			apiReader:             fakeclient.NewClient(&emptySecret),
			certificateDirectory:  oldDir,
			namespace:             "default",
			certificateSecretName: "my-cert",
			wantErr:               errors.New("certificate is outdated"),
		},
		{
			name:                  "good certificate - not stored",
			apiReader:             fakeclient.NewClient(&goodSecret),
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

func TestNewCertificateWatcher(t *testing.T) {
	logger := testr.New(t)
	client := fakeclient.NewClient()
	want := &CertificateWatcher{
		apiReader:             client,
		fs:                    afero.NewOsFs(),
		namespace:             "default",
		certificateSecretName: "my-secret",
		certificateDirectory:  "test",
		certificateTreshold:   CertThreshold,
		ICertificateHandler:   defaultCertificateHandler{},
		Log:                   testr.New(t),
	}
	got := NewCertificateWatcher(client, "test", "default", "my-secret", logger)
	require.EqualValues(t, got, want)

}
