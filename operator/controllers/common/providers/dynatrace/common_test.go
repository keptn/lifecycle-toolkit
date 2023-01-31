package dynatrace

import (
	"context"
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetSecret_NoKeyDefined(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(dtpayload))
		require.Nil(t, err)
	}))
	defer svr.Close()
	fakeClient := fake.NewClient()

	p := klcv1alpha2.KeptnEvaluationProvider{
		Spec: klcv1alpha2.KeptnEvaluationProviderSpec{
			TargetServer: svr.URL,
		},
	}
	r, e := getDTSecret(context.TODO(), p, fakeClient)
	require.NotNil(t, e)
	require.True(t, strings.Contains(e.Error(), "the SecretKeyRef property with the DT API token is missing"))
}

func TestGetSecret_NoSecret(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(dtpayload))
		require.Nil(t, err)
	}))
	defer svr.Close()
	fakeClient := fake.NewClient()

	p := klcv1alpha2.KeptnEvaluationProvider{
		Spec: klcv1alpha2.KeptnEvaluationProviderSpec{
			TargetServer: svr.URL,
		},
	}
	r, e := getDTSecret(context.TODO(), p, fakeClient)
	require.NotNil(t, e)
	require.True(t, strings.Contains(e.Error(), "the SecretKeyRef property with the DT API token is missing"))
}

func TestGetSecret_NoTokenData(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(dtpayload))
		require.Nil(t, err)
	}))
	defer svr.Close()
	fakeClient := fake.NewClient()

	p := klcv1alpha2.KeptnEvaluationProvider{
		Spec: klcv1alpha2.KeptnEvaluationProviderSpec{
			TargetServer: svr.URL,
		},
	}
	r, e := getDTSecret(context.TODO(), p, fakeClient)
	require.NotNil(t, e)
	require.True(t, strings.Contains(e.Error(), "the SecretKeyRef property with the DT API token is missing"))
}

func TestValidateOAuthSecret(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		result error
	}{
		{
			name:   "good token",
			input:  "",
			result: nil,
		},
		{
			name:   "wrong prefix",
			input:  "",
			result: nil,
		},
		{
			name:   "wrong format",
			input:  "",
			result: nil,
		},
		{
			name:   "wrong secret part",
			input:  "",
			result: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := validateOAuthSecret(tt.input)
			require.Equal(t, tt.result, e)
		})

	}
}
