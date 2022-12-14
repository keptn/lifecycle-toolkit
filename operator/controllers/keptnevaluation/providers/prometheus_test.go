package providers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

const promWarnPayload = "{\"status\":\"success\",\"warnings\":[\"awarning\"],\"data\":{\"resultType\":\"vector\",\"result\":[{\"metric\":{\"__name__\":\"kube_pod_info\",\"container\":\"kube-rbac-proxy-main\",\"created_by_kind\":\"DaemonSet\",\"created_by_name\":\"kindnet\",\"host_ip\":\"172.18.0.2\",\"host_network\":\"true\",\"instance\":\"10.244.0.24:8443\",\"job\":\"kube-state-metrics\",\"namespace\":\"kube-system\",\"node\":\"kind-control-plane\",\"pod\":\"kindnet-llt85\",\"pod_ip\":\"172.18.0.2\",\"uid\":\"0bb9d9db-2658-439f-aed9-ab3e8502397d\"},\"value\":[1669714193.275,\"1\"]}]}}"
const promPayload = "{\"status\":\"success\",\"data\":{\"resultType\":\"vector\",\"result\":[{\"metric\":{\"__name__\":\"kube_pod_info\",\"container\":\"kube-rbac-proxy-main\",\"created_by_kind\":\"DaemonSet\",\"created_by_name\":\"kindnet\",\"host_ip\":\"172.18.0.2\",\"host_network\":\"true\",\"instance\":\"10.244.0.24:8443\",\"job\":\"kube-state-metrics\",\"namespace\":\"kube-system\",\"node\":\"kind-control-plane\",\"pod\":\"kindnet-llt85\",\"pod_ip\":\"172.18.0.2\",\"uid\":\"0bb9d9db-2658-439f-aed9-ab3e8502397d\"},\"value\":[1669714193.275,\"1\"]}]}}"
const promEmptyDataPayload = "{\"status\":\"success\",\"data\":{\"resultType\":\"vector\",\"result\":[]}}"
const promMatrixPayload = "{\"status\":\"success\",\"data\":{\"resultType\":\"matrix\",\"result\":[]}}"
const promMultiPointPayload = "{\"status\":\"success\",\"data\":{\"resultType\":\"vector\",\"result\":[{\"metric\":{\"__name__\":\"kube_pod_info\",\"container\":\"kube-rbac-proxy-main\",\"created_by_kind\":\"DaemonSet\",\"created_by_name\":\"kindnet\",\"host_ip\":\"172.18.0.2\",\"host_network\":\"true\",\"instance\":\"10.244.0.24:8443\",\"job\":\"kube-state-metrics\",\"namespace\":\"kube-system\",\"node\":\"kind-control-plane\",\"pod\":\"kindnet-llt85\",\"pod_ip\":\"172.18.0.2\",\"uid\":\"0bb9d9db-2658-439f-aed9-ab3e8502397d\"},\"value\":[1669714193.275,\"1\"]},{\"metric\":{\"__name__\":\"kube_pod_info\",\"container\":\"kube-rbac-proxy-main\",\"created_by_kind\":\"DaemonSet\",\"created_by_name\":\"kube-proxy\",\"host_ip\":\"172.18.0.2\",\"host_network\":\"true\",\"instance\":\"10.244.0.24:8443\",\"job\":\"kube-state-metrics\",\"namespace\":\"kube-system\",\"node\":\"kind-control-plane\",\"pod\":\"kube-proxy-dlq7m\",\"pod_ip\":\"172.18.0.2\",\"priority_class\":\"system-node-critical\",\"uid\":\"31240e57-5286-4bc6-ad69-80b68bf806d0\"},\"value\":[1669714193.275,\"1\"]},{\"metric\":{\"__name__\":\"kube_pod_info\",\"container\":\"kube-rbac-proxy-main\",\"created_by_kind\":\"DaemonSet\",\"created_by_name\":\"node-exporter\",\"host_ip\":\"172.18.0.2\",\"host_network\":\"true\",\"instance\":\"10.244.0.24:8443\",\"job\":\"kube-state-metrics\",\"namespace\":\"monitoring\",\"node\":\"kind-control-plane\",\"pod\":\"node-exporter-dv6nr\",\"pod_ip\":\"172.18.0.2\",\"priority_class\":\"system-cluster-critical\",\"uid\":\"cf7baf10-ac9a-4b7d-9510-a6502d7ed271\"},\"value\":[1669714193.275,\"1\"]}]}}"

func TestHandlingWrongData(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("garbage"))
		require.Nil(t, err)
	}))
	defer svr.Close()

	kpp := KeptnPrometheusProvider{
		httpClient: http.Client{},
		Log:        ctrl.Log.WithName("testytest"),
	}
	obj := klcv1alpha2.Objective{
		Query: "garbage",
	}
	p := klcv1alpha2.KeptnEvaluationProvider{
		Spec: klcv1alpha2.KeptnEvaluationProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "myapitoken",
				},
				Key: "mykey",
			},
			TargetServer: svr.URL,
		},
	}
	r, e := kpp.EvaluateQuery(context.TODO(), obj, p)
	require.Equal(t, "", r)
	require.NotNil(t, e)
}

func TestHandlingWarnings(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(promWarnPayload))
		require.Nil(t, err)
	}))
	defer svr.Close()

	kpp := KeptnPrometheusProvider{
		httpClient: http.Client{},
		Log:        ctrl.Log.WithName("testytest"),
	}
	obj := klcv1alpha2.Objective{
		Query: "garbage",
	}
	p := klcv1alpha2.KeptnEvaluationProvider{
		Spec: klcv1alpha2.KeptnEvaluationProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "myapitoken",
				},
				Key: "mykey",
			},
			TargetServer: svr.URL,
		},
	}
	r, e := kpp.EvaluateQuery(context.TODO(), obj, p)
	require.Equal(t, "1", r)
	require.Nil(t, e)
}

func TestHandlingMultiDatapoint(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(promMultiPointPayload))
		require.Nil(t, err)
	}))
	defer svr.Close()

	kpp := KeptnPrometheusProvider{
		httpClient: http.Client{},
		Log:        ctrl.Log.WithName("testytest"),
	}
	obj := klcv1alpha2.Objective{
		Query: "garbage",
	}
	p := klcv1alpha2.KeptnEvaluationProvider{
		Spec: klcv1alpha2.KeptnEvaluationProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "myapitoken",
				},
				Key: "mykey",
			},
			TargetServer: svr.URL,
		},
	}
	r, e := kpp.EvaluateQuery(context.TODO(), obj, p)
	require.Equal(t, "", r)
	require.NotNil(t, e)
}

func TestHandlingEmptyDatapoint(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(promEmptyDataPayload))
		require.Nil(t, err)
	}))
	defer svr.Close()

	kpp := KeptnPrometheusProvider{
		httpClient: http.Client{},
		Log:        ctrl.Log.WithName("testytest"),
	}
	obj := klcv1alpha2.Objective{
		Query: "garbage",
	}
	p := klcv1alpha2.KeptnEvaluationProvider{
		Spec: klcv1alpha2.KeptnEvaluationProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "myapitoken",
				},
				Key: "mykey",
			},
			TargetServer: svr.URL,
		},
	}
	r, e := kpp.EvaluateQuery(context.TODO(), obj, p)
	require.Equal(t, "", r)
	require.NotNil(t, e)
}

func TestHandlingUnsupportedAnswerType(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(promMatrixPayload))
		require.Nil(t, err)
	}))
	defer svr.Close()

	kpp := KeptnPrometheusProvider{
		httpClient: http.Client{},
		Log:        ctrl.Log.WithName("testytest"),
	}
	obj := klcv1alpha2.Objective{
		Query: "garbage",
	}
	p := klcv1alpha2.KeptnEvaluationProvider{
		Spec: klcv1alpha2.KeptnEvaluationProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "myapitoken",
				},
				Key: "mykey",
			},
			TargetServer: svr.URL,
		},
	}
	r, e := kpp.EvaluateQuery(context.TODO(), obj, p)
	require.Equal(t, "", r)
	require.NotNil(t, e)
}

func TestHappyPath(t *testing.T) {
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(promPayload))
		require.Nil(t, err)
	}))
	defer svr.Close()

	kpp := KeptnPrometheusProvider{
		httpClient: http.Client{},
		Log:        ctrl.Log.WithName("testytest"),
	}
	obj := klcv1alpha2.Objective{
		Query: "garbage",
	}
	p := klcv1alpha2.KeptnEvaluationProvider{
		Spec: klcv1alpha2.KeptnEvaluationProviderSpec{
			SecretKeyRef: v1.SecretKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: "myapitoken",
				},
				Key: "mykey",
			},
			TargetServer: svr.URL,
		},
	}
	r, e := kpp.EvaluateQuery(context.TODO(), obj, p)
	require.Equal(t, "1", r)
	require.Nil(t, e)
}
