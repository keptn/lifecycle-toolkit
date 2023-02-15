/*
Copyright 2022.

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
	lifecyclev1alpha1 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha1"
	lifecyclev1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	metricsv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/apis/metrics/v1alpha1"
	optionsv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/apis/options/v1alpha1"
	cmdConfig "github.com/keptn/lifecycle-toolkit/operator/cmd/config"
	"github.com/keptn/lifecycle-toolkit/operator/cmd/metrics/adapter"
	"github.com/keptn/lifecycle-toolkit/operator/cmd/webhook"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnapp"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnappversion"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnevaluation"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptntask"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptntaskdefinition"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnworkload"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnworkloadinstance"
	keptnmetric "github.com/keptn/lifecycle-toolkit/operator/controllers/metrics"
	controlleroptions "github.com/keptn/lifecycle-toolkit/operator/controllers/options"
	keptnserver "github.com/keptn/lifecycle-toolkit/operator/pkg/metrics"
	"github.com/open-feature/go-sdk/pkg/openfeature"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/sdk/metric"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	//+kubebuilder:scaffold:imports
)

var (
	scheme                     = runtime.NewScheme()
	setupLog                   = ctrl.Log.WithName("setup")
	metricServerTickerInterval = 10 * time.Second
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(lifecyclev1alpha1.AddToScheme(scheme))
	utilruntime.Must(lifecyclev1alpha2.AddToScheme(scheme))
	utilruntime.Must(metricsv1alpha1.AddToScheme(scheme))
	utilruntime.Must(optionsv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

type envConfig struct {
	PodNamespace       string `envconfig:"POD_NAMESPACE" default:""`
	PodName            string `envconfig:"POD_NAME" default:""`
	ExposeKeptnMetrics bool   `envconfig:"EXPOSE_KEPTN_METRICS" default:"true"`

	KeptnAppControllerLogLevel              int `envconfig:"KEPTN_APP_CONTROLLER_LOG_LEVEL" default:"0"`
	KeptnAppVersionControllerLogLevel       int `envconfig:"KEPTN_APP_VERSION_CONTROLLER_LOG_LEVEL" default:"0"`
	KeptnEvaluationControllerLogLevel       int `envconfig:"KEPTN_EVALUATION_CONTROLLER_LOG_LEVEL" default:"0"`
	KeptnTaskControllerLogLevel             int `envconfig:"KEPTN_TASK_CONTROLLER_LOG_LEVEL" default:"0"`
	KeptnTaskDefinitionControllerLogLevel   int `envconfig:"KEPTN_TASK_DEFINITION_CONTROLLER_LOG_LEVEL" default:"0"`
	KeptnWorkloadControllerLogLevel         int `envconfig:"KEPTN_WORKLOAD_CONTROLLER_LOG_LEVEL" default:"0"`
	KeptnWorkloadInstanceControllerLogLevel int `envconfig:"KEPTN_WORKLOAD_INSTANCE_CONTROLLER_LOG_LEVEL" default:"0"`
	KeptnMetricControllerLogLevel           int `envconfig:"METRICS_CONTROLLER_LOG_LEVEL" default:"0"`
	KptnOptionsControllerLogLevel           int `envconfig:"OPTIONS_CONTROLLER_LOG_LEVEL" default:"0"`
}

//nolint:funlen,gocognit,gocyclo
func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}
	var metricsAddr string
	var enableLeaderElection bool
	var disableWebhook bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")

	// OTEL SET UP
	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.

	exporter, err := otelprom.New()
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("keptn/task")
	keptnMeters := controllercommon.SetUpKeptnTaskMeters(meter)

	// Start the prometheus HTTP server and pass the exporter Collector to it
	go serveMetrics()

	// Start the custom metrics adapter
	go startCustomMetricsAdapter(env.PodNamespace)

	// As recommended by the kubebuilder docs, webhook registration should be disabled if running locally. See https://book.kubebuilder.io/cronjob-tutorial/running.html#running-webhooks-locally for reference
	flag.BoolVar(&disableWebhook, "disable-webhook", false, "Disable the registration of webhooks.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
	disableCacheFor := []ctrlclient.Object{&corev1.Secret{}}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "6b866dd9.keptn.sh",
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

	// Enabling OTel
	err = controllercommon.GetOtelInstance().InitOtelCollector("")
	if err != nil {
		setupLog.Error(err, "unable to initialize OTel tracer options")
	}

	spanHandler := &controllercommon.SpanHandler{}

	taskReconciler := &keptntask.KeptnTaskReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		Log:           ctrl.Log.WithName("KeptnTask Controller"),
		Recorder:      mgr.GetEventRecorderFor("keptntask-controller"),
		Meters:        keptnMeters,
		TracerFactory: controllercommon.GetOtelInstance(),
	}
	taskReconciler.Log.V(env.KeptnTaskControllerLogLevel)
	if err = (taskReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnTask")
		os.Exit(1)
	}

	taskDefinitionReconciler := &keptntaskdefinition.KeptnTaskDefinitionReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Log:      ctrl.Log.WithName("KeptnTaskDefinition Controller"),
		Recorder: mgr.GetEventRecorderFor("keptntaskdefinition-controller"),
	}
	taskDefinitionReconciler.Log.V(env.KeptnTaskDefinitionControllerLogLevel)
	if err = (taskDefinitionReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnTaskDefinition")
		os.Exit(1)
	}

	appReconciler := &keptnapp.KeptnAppReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		Log:           ctrl.Log.WithName("KeptnApp Controller"),
		Recorder:      mgr.GetEventRecorderFor("keptnapp-controller"),
		TracerFactory: controllercommon.GetOtelInstance(),
	}
	appReconciler.Log.V(env.KeptnAppControllerLogLevel)
	if err = (appReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnApp")
		os.Exit(1)
	}

	workloadReconciler := &keptnworkload.KeptnWorkloadReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		Log:           ctrl.Log.WithName("KeptnWorkload Controller"),
		Recorder:      mgr.GetEventRecorderFor("keptnworkload-controller"),
		TracerFactory: controllercommon.GetOtelInstance(),
	}
	workloadReconciler.Log.V(env.KeptnWorkloadControllerLogLevel)
	if err = (workloadReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnWorkload")
		os.Exit(1)
	}

	workloadInstanceReconciler := &keptnworkloadinstance.KeptnWorkloadInstanceReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		Log:           ctrl.Log.WithName("KeptnWorkloadInstance Controller"),
		Recorder:      mgr.GetEventRecorderFor("keptnworkloadinstance-controller"),
		Meters:        keptnMeters,
		TracerFactory: controllercommon.GetOtelInstance(),
		SpanHandler:   spanHandler,
	}
	workloadInstanceReconciler.Log.V(env.KeptnWorkloadInstanceControllerLogLevel)
	if err = (workloadInstanceReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnWorkloadInstance")
		os.Exit(1)
	}

	appVersionReconciler := &keptnappversion.KeptnAppVersionReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		Log:           ctrl.Log.WithName("KeptnAppVersion Controller"),
		Recorder:      mgr.GetEventRecorderFor("keptnappversion-controller"),
		TracerFactory: controllercommon.GetOtelInstance(),
		Meters:        keptnMeters,
		SpanHandler:   spanHandler,
	}
	appVersionReconciler.Log.V(env.KeptnAppVersionControllerLogLevel)
	if err = (appVersionReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnAppVersion")
		os.Exit(1)
	}

	evaluationReconciler := &keptnevaluation.KeptnEvaluationReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		Log:           ctrl.Log.WithName("KeptnEvaluation Controller"),
		Recorder:      mgr.GetEventRecorderFor("keptnevaluation-controller"),
		TracerFactory: controllercommon.GetOtelInstance(),
		Meters:        keptnMeters,
		Namespace:     env.PodNamespace,
	}
	evaluationReconciler.Log.V(env.KeptnEvaluationControllerLogLevel)
	if err = (evaluationReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnEvaluation")
		os.Exit(1)
	}

	metricsReconciler := &keptnmetric.KeptnMetricReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log.WithName("KeptnMetric Controller"),
	}
	metricsReconciler.Log.V(env.KeptnMetricControllerLogLevel)
	if err = (metricsReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnMetric")
		os.Exit(1)
	}

	configReconciler := &controlleroptions.KeptnConfigReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log.WithName("KeptnConfig Controller"),
	}
	configReconciler.Log.V(env.KptnOptionsControllerLogLevel)
	if err = (configReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnConfig")
		os.Exit(1)
	}

	if err = (&lifecyclev1alpha2.KeptnApp{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "KeptnApp")
		os.Exit(1)
	}
	if err = (&lifecyclev1alpha2.KeptnEvaluationProvider{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "KeptnEvaluationProvider")
		os.Exit(1)
	}
	if err = (&lifecyclev1alpha2.KeptnAppVersion{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "KeptnAppVersion")
		os.Exit(1)
	}
	if err = (&lifecyclev1alpha2.KeptnWorkloadInstance{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "KeptnWorkloadInstance")
		os.Exit(1)
	}
	if err = (&metricsv1alpha1.KeptnMetric{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "KeptnMetric")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	controllercommon.SetUpKeptnMeters(meter, mgr.GetClient())

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}
	if !disableWebhook {
		webhookBuilder := webhook.NewWebhookBuilder().
			SetNamespace(env.PodNamespace).
			SetPodName(env.PodName).
			SetConfigProvider(cmdConfig.NewKubeConfigProvider())

		setupLog.Info("starting webhook and manager")
		if err1 := webhookBuilder.Run(mgr); err1 != nil {
			setupLog.Error(err, "problem running manager")
			os.Exit(1)
		}

	} else {
		flag.Parse()
		setupLog.Info("starting manager")
		setupLog.Info("Keptn lifecycle operator is alive")
		if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
			setupLog.Error(err, "problem running manager")
			os.Exit(1)
		}
	}

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

func startCustomMetricsAdapter(namespace string) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer cancel()

	adapter := adapter.MetricsAdapter{KltNamespace: namespace}
	adapter.RunAdapter(ctx)
}
