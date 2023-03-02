package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"

	"github.com/benbjohnson/clock"
	"github.com/gorilla/mux"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha2"
	"github.com/open-feature/go-sdk/pkg/openfeature"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Metrics struct {
	gauges map[string]prometheus.Gauge
}

var metrics Metrics

var instance *serverManager
var smOnce sync.Once

type serverManager struct {
	server        *http.Server
	ticker        *clock.Ticker
	ofClient      *openfeature.Client
	exposeMetrics bool
	k8sClient     client.Client
}

// StartServerManager starts a server manager to expose metrics and runs until
// the context is cancelled (i.e. an env variable gets changes and pod is restarted)
func StartServerManager(ctx context.Context, client client.Client, ofClient *openfeature.Client, exposeMetrics bool, interval time.Duration) {
	smOnce.Do(func() {
		metrics.gauges = make(map[string]prometheus.Gauge)
		instance = &serverManager{
			ticker:        clock.New().Ticker(interval),
			ofClient:      ofClient,
			exposeMetrics: exposeMetrics,
			k8sClient:     client,
		}
		instance.start(ctx)
	})
}

func (m *serverManager) start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				if err := m.shutDownServer(); err != nil {
					klog.Errorf("Error during server shutdown: %v", err)
				}
				return
			case <-m.ticker.C:
				if err := m.setup(); err != nil {
					klog.Errorf("Error during server setup: %v", err)
				}
			}
		}
	}()
}

func (m *serverManager) shutDownServer() error {
	defer func() {
		m.server = nil
	}()
	if m.server != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
		defer cancel()
		return m.server.Shutdown(ctx)
	}
	return nil
}

func (m *serverManager) setup() error {
	maxRetries := 3
	var serverEnabled bool
	var err error

	klog.Infof("Checking configuration of keptn-metrics server")

	for i := 0; i < maxRetries; i++ {
		serverEnabled, err = m.ofClient.BooleanValue(context.TODO(), "keptn.gms.expose", m.exposeMetrics, openfeature.EvaluationContext{})
		if err == nil {
			break
		}

		if strings.Contains(err.Error(), string(openfeature.ProviderNotReadyCode)) {
			<-time.After(2 * time.Second)
			continue
		}
		break
	}

	klog.Infof("Keptn Metrics server enabled: %v", serverEnabled)

	if serverEnabled && m.server == nil {
		klog.Infof("serving Prometheus metrics at localhost:9999/metrics")
		klog.Infof("serving KeptnMetrics at localhost:9999/api/v1/metrics/{namespace}/{metric}")

		router := mux.NewRouter()
		router.Path("/metrics").Handler(promhttp.Handler())
		router.Path("/api/v1/metrics/{namespace}/{metric}").HandlerFunc(m.returnMetric)

		m.server = &http.Server{
			Addr:    ":9999",
			Handler: router,
		}

		m.recordMetrics()

		go func() {
			err := m.server.ListenAndServe()
			if err != nil {
				klog.Errorf("could not start keptn-metrics server: %w", err)
			}
		}()

	} else if !serverEnabled && m.server != nil {
		if err := m.shutDownServer(); err != nil {
			return fmt.Errorf("could not shut down keptn-metrics server: %w", err)
		}
	}
	return nil
}

func (m *serverManager) returnMetric(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namespace := vars["namespace"]
	metric := vars["metric"]

	metricObj := metricsapi.KeptnMetric{}
	err := m.k8sClient.Get(context.Background(), types.NamespacedName{Name: metric, Namespace: namespace}, &metricObj)
	if err != nil {
		fmt.Println("failed to list keptn-metrics: " + err.Error())
	}

	data := map[string]string{
		"namespace": namespace,
		"metric":    metric,
		"value":     metricObj.Status.Value,
	}

	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		fmt.Println("failed to encode data")
		os.Exit(1)
	}
}

func (m *serverManager) recordMetrics() {
	go func() {
		for {
			if m.server == nil {
				return
			}
			list := metricsapi.KeptnMetricList{}
			err := m.k8sClient.List(context.Background(), &list)
			if err != nil {
				fmt.Println("failed to list keptn-metrics" + err.Error())
			}
			for _, metric := range list.Items {
				normName := normalizeMetricName(metric.Name)
				if _, ok := metrics.gauges[normName]; !ok {
					metrics.gauges[normName] = prometheus.NewGauge(prometheus.GaugeOpts{
						Name: normName,
						Help: metric.Name,
					})
					prometheus.MustRegister(metrics.gauges[normName])
				}
				val, _ := strconv.ParseFloat(metric.Status.Value, 64)
				metrics.gauges[normName].Set(val)
			}
			<-time.After(10 * time.Second)
		}
	}()
}

// normalizeMetricName removes all characters from the name
// of the metric that are not digits nor letters and
// substitues them with underscore (_)
func normalizeMetricName(s string) string {
	return strings.Join(strings.FieldsFunc(s,
		func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsDigit(r)
		}), "_")
}
