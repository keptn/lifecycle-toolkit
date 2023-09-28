package dynatrace

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/fake"
	"github.com/stretchr/testify/require"
)

func TestGetSecret_NoKeyDefined(t *testing.T) {

	fakeClient := fake.NewClient()

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: "svr.URL",
		},
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
			TargetServer: "svr.URL",
		},
	}
	r, e := getDTSecret(context.TODO(), p, fakeClient)
	require.NotNil(t, e)
	require.ErrorIs(t, e, ErrSecretKeyRefNotDefined)
	require.Empty(t, r)
}

func TestGetSecret_NoTokenData(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(dtpayload))
		require.Nil(t, err)
	}))
	defer svr.Close()
	fakeClient := fake.NewClient()

	p := metricsapi.KeptnMetricsProvider{
		Spec: metricsapi.KeptnMetricsProviderSpec{
			TargetServer: svr.URL,
		},
	}
	r, e := getDTSecret(context.TODO(), p, fakeClient)
	require.NotNil(t, e)
	require.ErrorIs(t, e, ErrSecretKeyRefNotDefined)
	require.Empty(t, r)
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
