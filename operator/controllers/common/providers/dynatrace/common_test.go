package dynatrace

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
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
	require.Empty(t, r)
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
	require.Empty(t, r)
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
	require.Empty(t, r)
}
