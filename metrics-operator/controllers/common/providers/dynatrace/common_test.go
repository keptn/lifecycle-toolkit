package dynatrace

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const dtSecretToken = "dt0s08.XX.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
const dqlSecretString = `{"token": "dt0s08.XX.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", "authUrl":"https://my-auth-url.test"}`
const dqlSecretStringWithEmptyAuthURL = `{"token": "dt0s08.XX.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"}`

func TestGetSecret_NoKeyDefined(t *testing.T) {

	fakeClient := fake.NewClient()

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{},
	}
	r, e := getDTSecret(context.TODO(), p, fakeClient)
	require.NotNil(t, e)
	require.ErrorIs(t, e, ErrSecretKeyRefNotDefined)
	require.Empty(t, r)
}

func TestGetSecret_NoSecret(t *testing.T) {
	fakeClient := fake.NewClient()

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "my-dql-secret",
				},
				Key: "creds",
			},
		},
	}
	r, e := getDTSecret(context.TODO(), p, fakeClient)
	require.NotNil(t, e)
	require.Empty(t, r)
}

func TestGetSecret_NoTokenData(t *testing.T) {
	dtSecret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-dql-secret",
		},
		Data: map[string][]byte{},
		Type: v1.SecretTypeOpaque,
	}
	fakeClient := fake.NewClient(dtSecret)

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "my-dql-secret",
				},
				Key: "creds",
			},
		},
	}
	r, e := getDTSecret(context.TODO(), p, fakeClient)
	require.NotNil(t, e)
	require.Empty(t, r)
}

func TestGetSecret_ValidSecret(t *testing.T) {
	dtSecret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-dql-secret",
		},
		Data: map[string][]byte{
			"creds": []byte(dtSecretToken),
		},
		Type: v1.SecretTypeOpaque,
	}
	fakeClient := fake.NewClient(dtSecret)

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "my-dql-secret",
				},
				Key: "creds",
			},
		},
	}
	r, e := getDTSecret(context.TODO(), p, fakeClient)
	require.Nil(t, e)
	require.Equal(t, dtSecretToken, r)
}

func TestGetDQLSecret_NoKeyDefined(t *testing.T) {

	fakeClient := fake.NewClient()

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: "svr.URL",
		},
	}
	r, e := getDQLSecret(context.TODO(), p, fakeClient)
	require.NotNil(t, e)
	require.ErrorIs(t, e, ErrSecretKeyRefNotDefined)
	require.Empty(t, r)
}

func TestGetDQLSecret_NoSecret(t *testing.T) {
	fakeClient := fake.NewClient()

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "my-dql-secret",
				},
				Key: "creds",
			},
		},
	}
	r, e := getDQLSecret(context.TODO(), p, fakeClient)
	require.NotNil(t, e)
	require.Empty(t, r)
}

func TestGetDQLSecret_ValidSecret(t *testing.T) {
	dqlSecret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-dql-secret",
		},
		Data: map[string][]byte{
			"creds": []byte(dqlSecretString),
		},
		Type: v1.SecretTypeOpaque,
	}
	fakeClient := fake.NewClient(dqlSecret)

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "my-dql-secret",
				},
				Key: "creds",
			},
		},
	}
	secretValue, e := getDQLSecret(context.TODO(), p, fakeClient)
	require.Nil(t, e)
	require.NotEmpty(t, secretValue)
	require.Equal(t, "dt0s08.XX.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", secretValue.Token)
	require.Equal(t, "https://my-auth-url.test", secretValue.AuthUrl)
}

func TestGetDQLSecret_ValidSecretEmptyAuthURL(t *testing.T) {
	dqlSecret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-dql-secret",
		},
		Data: map[string][]byte{
			"creds": []byte(dqlSecretStringWithEmptyAuthURL),
		},
		Type: v1.SecretTypeOpaque,
	}
	fakeClient := fake.NewClient(dqlSecret)

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "my-dql-secret",
				},
				Key: "creds",
			},
		},
	}
	secretValue, e := getDQLSecret(context.TODO(), p, fakeClient)
	require.Nil(t, e)
	require.NotEmpty(t, secretValue)
	require.Equal(t, "dt0s08.XX.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX", secretValue.Token)
	require.Equal(t, "", secretValue.AuthUrl)
}

func TestGetDQLSecret_SecretIsNotAJsonObject(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(dtpayload))
		require.Nil(t, err)
	}))
	defer svr.Close()

	dqlSecret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-dql-secret",
		},
		Data: map[string][]byte{
			"creds": []byte("wrong"),
		},
		Type: v1.SecretTypeOpaque,
	}
	fakeClient := fake.NewClient(dqlSecret)

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: svr.URL,
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "my-dql-secret",
				},
				Key: "creds",
			},
		},
	}
	secretValue, e := getDQLSecret(context.TODO(), p, fakeClient)
	require.NotNil(t, e)
	require.Nil(t, secretValue)
}

func Test_urlEncodeQuery(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "encode single parameter",
			args: args{
				query: "metricSelector=my:metric-selector",
			},
			want: "metricSelector=my%3Ametric-selector",
		},
		{
			name: "encode multiple parameters",
			args: args{
				query: "metricSelector=my:metric-selector&from=now-2h",
			},
			want: "metricSelector=my%3Ametric-selector&from=now-2h",
		},
		{
			name: "omit wrongly formatted input",
			args: args{
				query: "metricSelector=my:metric-selector&from",
			},
			want: "metricSelector=my%3Ametric-selector",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := urlEncodeQuery(tt.args.query)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestDQLSecret_validate(t *testing.T) {
	type fields struct {
		Token   string
		AuthUrl string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr error
	}{
		{
			name: "good token",
			fields: fields{
				Token:   "dt0s08.XX.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				AuthUrl: "https://dev.token.internal.dynatracelabs.com/sso/oauth2/token",
			},
			wantErr: nil,
		},
		{
			name: "empty token",
			fields: fields{
				Token:   "",
				AuthUrl: "",
			},
			wantErr: ErrInvalidToken,
		},
		{
			name: "wrong format",
			fields: fields{
				Token: "dt0s08.wrong.token.with.too.many.components",
			},
			wantErr: ErrInvalidToken,
		},
		{
			name: "wrong format",
			fields: fields{
				Token: "dt0s08.wrong.length",
			},
			wantErr: ErrInvalidToken,
		},
		{
			name: "invalid auth url",
			fields: fields{
				Token:   "dt0s08.XX.XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
				AuthUrl: "wrong",
			},
			wantErr: ErrInvalidAuthURL,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := DQLSecret{
				Token:   tt.fields.Token,
				AuthUrl: tt.fields.AuthUrl,
			}
			err := s.validate()
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
