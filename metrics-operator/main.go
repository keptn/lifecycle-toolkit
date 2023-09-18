/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/keptn/lifecycle-toolkit/klt-cert-manager/pkg/certificates"
	certCommon "github.com/keptn/lifecycle-toolkit/klt-cert-manager/pkg/common"
	"github.com/keptn/lifecycle-toolkit/klt-cert-manager/pkg/webhook"
	metricsv1alpha1 "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha1"
	metricsv1alpha2 "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha2"
	metricsv1alpha3 "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/cmd/metrics/adapter"
	analysiscontroller "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/analysis"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/common/analysis"
	metricscontroller "github.com/keptn/lifecycle-toolkit/metrics-operator/controllers/metrics"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/converter"
	keptnserver "github.com/keptn/lifecycle-toolkit/metrics-operator/pkg/metrics"
	"github.com/open-feature/go-sdk/pkg/openfeature"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"

	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var (
	scheme                     = runtime.NewScheme()
	setupLog                   = ctrl.Log.WithName("setup")
	metricServerTickerInterval = 10 * time.Second
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(metricsv1alpha1.AddToScheme(scheme))
	utilruntime.Must(metricsv1alpha2.AddToScheme(scheme))
	utilruntime.Must(metricsv1alpha3.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

type envConfig struct {
	PodNamespace                  string `envconfig:"POD_NAMESPACE" default:""`
	PodName                       string `envconfig:"POD_NAME" default:""`
	KeptnMetricControllerLogLevel int    `envconfig:"METRICS_CONTROLLER_LOG_LEVEL" default:"0"`
	AnalysisControllerLogLevel    int    `envconfig:"ANALYSIS_CONTROLLER_LOG_LEVEL" default:"0"`
	ExposeKeptnMetrics            bool   `envconfig:"EXPOSE_KEPTN_METRICS" default:"true"`
	EnableKeptnAnalysis           bool   `envconfig:"ENABLE_ANALYSIS" default:"false"`
}

//nolint:gocyclo,funlen
func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}
	var metricsAddr string
	var SLIFilePath string
	var provider string
	var namespace string
	var SLOFilePath string
	var analysisDefinition string
	var enableLeaderElection bool
	var disableWebhook bool
	var probeAddr string
	flag.StringVar(&SLIFilePath, "convert-sli", "", "The path the the SLI file to be converted")
	flag.StringVar(&provider, "keptn-provider-name", "", "The name of KeptnMetricsProvider referenced in KeptnValueTemplates")
	flag.StringVar(&namespace, "keptn-provider-namespace", "", "The namespace of the referenced KeptnMetricsProvider")
	flag.StringVar(&SLOFilePath, "convert-slo", "", "The path the the SLO file to be converted")
	flag.StringVar(&analysisDefinition, "analysis-definition-name", "", "The name of AnalysisDefinition to be created")
	flag.StringVar(&namespace, "analysis-value-template-namespace", "", "The namespace of the referenced AnalysisValueTemplate")
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&disableWebhook, "disable-webhook", false, "Disable the registration of webhooks.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	if SLIFilePath != "" {
		// convert
		content, err := convertSLI(SLIFilePath, provider, namespace)
		if err != nil {
			log.Fatalf(err.Error())
			return
		}
		// write out converted result
		fmt.Print(content)
		return
	}

	if SLOFilePath != "" {
		// convert
		content, err := convertSLO(SLOFilePath, analysisDefinition, namespace)
		if err != nil {
			log.Fatalf(err.Error())
			return
		}
		// write out converted result
		fmt.Print(content)
		return
	}

	exporter, err := otelprom.New()
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	metricProvider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := metricProvider.Meter("keptn/metric")

	// Initialize your metric
	keptnMetricActive, err := meter.Int64Counter("keptn_metric_active")
	if err != nil {
		setupLog.Error(err, "unable to create metric keptn_metric_active")
		os.Exit(1)
	}

	go serveMetrics()

	// Set the metric value as soon as the operator starts
	keptnMetricActive.Add(context.Background(), 1)

	// Start the custom metrics adapter
	go startCustomMetricsAdapter(env.PodNamespace)

	disableCacheFor := []ctrlclient.Object{&corev1.Secret{}}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "3f8532ca.keptn.sh",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
		ClientDisableCacheFor: disableCacheFor, // due to https://github.com/kubernetes-sigs/controller-runtime/issues/550
		// We disable secret informer cache so that the operator won't need clusterrole list access to secrets
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	keptnserver.StartServerManager(ctx, mgr.GetClient(), openfeature.NewClient("klt"), env.ExposeKeptnMetrics, metricServerTickerInterval)

	metricsLogger := ctrl.Log.WithName("KeptnMetric Controller")
	if err = (&metricscontroller.KeptnMetricReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    metricsLogger.V(env.KeptnMetricControllerLogLevel),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnMetric")
		os.Exit(1)
	}

	if env.EnableKeptnAnalysis {

		analysisLogger := ctrl.Log.WithName("KeptnAnalysis Controller")
		targetEval := analysis.NewTargetEvaluator(&analysis.OperatorEvaluator{})
		objEval := analysis.NewObjectiveEvaluator(&targetEval)
		analysisEval := analysis.NewAnalysisEvaluator(&objEval)

		if err = (&analysiscontroller.AnalysisReconciler{
			Client:                mgr.GetClient(),
			Scheme:                mgr.GetScheme(),
			Log:                   analysisLogger.V(env.AnalysisControllerLogLevel),
			MaxWorkers:            2,
			Namespace:             env.PodNamespace,
			NewWorkersPoolFactory: analysiscontroller.NewWorkersPool,
			IAnalysisEvaluator:    &analysisEval,
		}).SetupWithManager(mgr); err != nil {
			setupLog.Error(err, "unable to create controller", "controller", "KeptnMetric")
			os.Exit(1)
		}
	}
	// +kubebuilder:scaffold:builder

	setupValidationWebhooks(mgr)
	setupProbes(mgr)

	if !disableWebhook {
		webhookBuilder := webhook.NewWebhookBuilder().
			SetNamespace(env.PodNamespace).
			SetPodName(env.PodName).
			SetManagerProvider(
				webhook.NewWebhookManagerProvider(
					mgr.GetWebhookServer().CertDir, "tls.key", "tls.crt"),
			).
			SetCertificateWatcher(
				certificates.NewCertificateWatcher(
					mgr.GetAPIReader(),
					mgr.GetWebhookServer().CertDir,
					env.PodNamespace,
					certCommon.SecretName,
					setupLog,
				),
			)

		setupLog.Info("starting webhook and manager")
		if err := webhookBuilder.Run(mgr, nil); err != nil {
			setupLog.Error(err, "problem running manager")
			os.Exit(1)
		}

	} else {
		flag.Parse()
		setupLog.Info("starting manager")
		setupLog.Info("Keptn metrics-operator is alive")
		if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
			setupLog.Error(err, "problem running manager")
			os.Exit(1)
		}
	}
}

func setupValidationWebhooks(mgr manager.Manager) {
	if err := (&metricsv1alpha3.KeptnMetric{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "KeptnMetric")
		os.Exit(1)
	}
	if err := (&metricsv1alpha3.AnalysisDefinition{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "AnalysisDefinition")
		os.Exit(1)
	}
}

func setupProbes(mgr manager.Manager) {

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}
}

func startCustomMetricsAdapter(namespace string) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer cancel()

	metricsAdapter := adapter.MetricsAdapter{KltNamespace: namespace}
	metricsAdapter.RunAdapter(ctx)
}

func serveMetrics() {
	log.Printf("serving metrics at localhost:2222/metrics")

	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":2222", nil)
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}

func convertSLI(SLIFilePath, provider, namespace string) (string, error) {
	//read file content
	fileContent, err := os.ReadFile(SLIFilePath)
	if err != nil {
		return "", fmt.Errorf("error reading file content: %s", err.Error())
	}

	// convert
	c := converter.NewSLIConverter()
	content, err := c.Convert(fileContent, provider, namespace)
	if err != nil {
		return "", err
	}

	return content, nil
}

func convertSLO(SLOFilePath, analysisDefinition, namespace string) (string, error) {
	//read file content
	fileContent, err := os.ReadFile(SLOFilePath)
	if err != nil {
		return "", fmt.Errorf("error reading file content: %s", err.Error())
	}

	// convert
	c := converter.NewSLOConverter()
	content, err := c.Convert(fileContent, analysisDefinition, namespace)
	if err != nil {
		return "", err
	}

	return content, nil
}
