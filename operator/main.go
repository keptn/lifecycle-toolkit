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
	"time"

	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/keptn-sandbox/lifecycle-controller/operator/controllers/keptnappversion"

	"github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1/common"

	"github.com/keptn-sandbox/lifecycle-controller/operator/controllers/keptnworkload"
	"github.com/keptn-sandbox/lifecycle-controller/operator/controllers/keptnworkloadinstance"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/keptn-sandbox/lifecycle-controller/operator/controllers/keptnapp"
	"github.com/keptn-sandbox/lifecycle-controller/operator/controllers/keptnevaluation"
	"github.com/keptn-sandbox/lifecycle-controller/operator/controllers/keptntask"
	"github.com/keptn-sandbox/lifecycle-controller/operator/controllers/keptntaskdefinition"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"

	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/metric"

	"sigs.k8s.io/controller-runtime/pkg/webhook"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	lifecyclev1alpha1 "github.com/keptn-sandbox/lifecycle-controller/operator/api/v1alpha1"

	"github.com/keptn-sandbox/lifecycle-controller/operator/webhooks"
	//+kubebuilder:scaffold:imports
)

var (
	scheme       = runtime.NewScheme()
	setupLog     = ctrl.Log.WithName("setup")
	gitCommit    string
	buildTime    string
	buildVersion string
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(lifecyclev1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

type envConfig struct {
	OTelCollectorURL string `envconfig:"OTEL_COLLECTOR_URL" default:""`
}

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

	// OTEL SETUP
	// The exporter embeds a default OpenTelemetry Reader and
	// implements prometheus.Collector, allowing it to be used as
	// both a Reader and Collector.

	exporter := otelprom.New()
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("keptn/task")
	deploymentCount, err := meter.SyncInt64().Counter("keptn.deployment.count", instrument.WithDescription("a simple counter for Keptn deployment"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	deploymentDuration, err := meter.SyncFloat64().Histogram("keptn.deployment.duration", instrument.WithDescription("a histogram of duration for Keptn deployment"), instrument.WithUnit(unit.Unit("s")))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	deploymentActive, err := meter.SyncInt64().UpDownCounter("keptn.deployment.active", instrument.WithDescription("a simple counter of active deployments for Keptn deployment"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	taskCount, err := meter.SyncInt64().Counter("keptn.task.count", instrument.WithDescription("a simple counter for Keptn tasks"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	taskDuration, err := meter.SyncFloat64().Histogram("keptn.task.duration", instrument.WithDescription("a histogram of duration for Keptn tasks"), instrument.WithUnit(unit.Unit("s")))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	taskActive, err := meter.SyncInt64().UpDownCounter("keptn.task.active", instrument.WithDescription("a simple counter of active tasks for Keptn tasks"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	appCount, err := meter.SyncInt64().Counter("keptn.app.count", instrument.WithDescription("a simple counter for Keptn apps"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	appDuration, err := meter.SyncFloat64().Histogram("keptn.app.duration", instrument.WithDescription("a histogram of duration for Keptn apps"), instrument.WithUnit(unit.Unit("s")))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	appActive, err := meter.SyncInt64().UpDownCounter("keptn.app.active", instrument.WithDescription("a simple counter of active apps for Keptn apps"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	analysisCount, err := meter.SyncInt64().Counter("keptn.analysis.count", instrument.WithDescription("a simple counter for Keptn analysis for Evaluations"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	analysisDuration, err := meter.SyncFloat64().Histogram("keptn.analysis.duration", instrument.WithDescription("a histogram of duration for Keptn analysis for Evaluations"), instrument.WithUnit(unit.Unit("s")))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	analysisActive, err := meter.SyncInt64().UpDownCounter("keptn.analysis.active", instrument.WithDescription("a simple counter of active apps for Keptn analysis for Evaluations"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}

	meters := common.KeptnMeters{
		TaskCount:          taskCount,
		TaskDuration:       taskDuration,
		TaskActive:         taskActive,
		DeploymentCount:    deploymentCount,
		DeploymentDuration: deploymentDuration,
		DeploymentActive:   deploymentActive,
		AppCount:           appCount,
		AppDuration:        appDuration,
		AppActive:          appActive,
		AnalysisCount:      analysisCount,
		AnalysisDuration:   analysisDuration,
		AnalysusActive:     analysisActive,
	}

	// Start the prometheus HTTP server and pass the exporter Collector to it
	go serveMetrics(exporter.Collector)

	// As recommended by the kubebuilder docs, webhook registration should be disabled if running locally. See https://book.kubebuilder.io/cronjob-tutorial/running.html#running-webhooks-locally for reference
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

	// Enabling OTel
	tpOptions, err := getOTelTracerProviderOptions(env)
	if err != nil {
		setupLog.Error(err, "unable to initialize OTel tracer options")
	}

	tp := trace.NewTracerProvider(tpOptions...)

	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			setupLog.Error(err, "unable to shutdown  OTel exporter")
			os.Exit(1)
		}
	}()
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

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
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if !disableWebhook {
		mgr.GetWebhookServer().Register("/mutate-v1-pod", &webhook.Admission{
			Handler: &webhooks.PodMutatingWebhook{
				Client:   mgr.GetClient(),
				Tracer:   otel.Tracer("keptn/webhook"),
				Recorder: mgr.GetEventRecorderFor("keptn/webhook"),
				Log:      ctrl.Log.WithName("Mutating Webhook"),
			}})
	}
	if err = (&keptntask.KeptnTaskReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Log:      ctrl.Log.WithName("KeptnTask Controller"),
		Recorder: mgr.GetEventRecorderFor("keptntask-controller"),
		Meters:   meters,
		Tracer:   otel.Tracer("keptn/operator/task"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnTask")
		os.Exit(1)
	}
	if err = (&keptntaskdefinition.KeptnTaskDefinitionReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Log:      ctrl.Log.WithName("KeptnTaskDefinition Controller"),
		Recorder: mgr.GetEventRecorderFor("keptntaskdefinition-controller"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnTaskDefinition")
		os.Exit(1)
	}
	if err = (&keptnapp.KeptnAppReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Log:      ctrl.Log.WithName("KeptnApp Controller"),
		Recorder: mgr.GetEventRecorderFor("keptnapp-controller"),
		Tracer:   otel.Tracer("keptn/operator/app"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnApp")
		os.Exit(1)
	}
	if err = (&keptnworkload.KeptnWorkloadReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Log:      ctrl.Log.WithName("KeptnWorkload Controller"),
		Recorder: mgr.GetEventRecorderFor("keptnworkload-controller"),
		Tracer:   otel.Tracer("keptn/operator/workload"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnWorkload")
		os.Exit(1)
	}
	if err = (&keptnworkloadinstance.KeptnWorkloadInstanceReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Log:      ctrl.Log.WithName("KeptnWorkloadInstance Controller"),
		Recorder: mgr.GetEventRecorderFor("keptnworkloadinstance-controller"),
		Meters:   meters,
		Tracer:   otel.Tracer("keptn/operator/workloadinstance"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnWorkloadInstance")
		os.Exit(1)
	}
	if err = (&keptnappversion.KeptnAppVersionReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Log:      ctrl.Log.WithName("KeptnAppVersion Controller"),
		Recorder: mgr.GetEventRecorderFor("keptnappversion-controller"),
		Tracer:   otel.Tracer("keptn/operator/appversion"),
		Meters:   meters,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnAppVersion")
		os.Exit(1)
	}
	if err = (&keptnevaluation.KeptnEvaluationReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Log:      ctrl.Log.WithName("KeptnAppVersion Controller"),
		Recorder: mgr.GetEventRecorderFor("keptnappversion-controller"),
		Tracer:   otel.Tracer("keptn/operator/appversion"),
		Meters:   meters,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnEvaluation")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	setupLog.Info("Keptn lifecycle operator is alive")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func getOTelTracerProviderOptions(env envConfig) ([]trace.TracerProviderOption, error) {
	tracerProviderOptions := []trace.TracerProviderOption{}

	stdOutExp, err := newStdOutExporter()
	if err != nil {
		return nil, fmt.Errorf("could not create stdout OTel exporter: %w", err)
	}
	tracerProviderOptions = append(tracerProviderOptions, trace.WithBatcher(stdOutExp))

	if env.OTelCollectorURL != "" {
		// try to set OTel exporter for Jaeger
		otelExporter, err := newOTelExporter(env)
		if err != nil {
			// log the error, but do not break if Jaeger exporter cannot be created
			setupLog.Error(err, "Could not set up OTel exporter")
		} else if otelExporter != nil {
			tracerProviderOptions = append(tracerProviderOptions, trace.WithBatcher(otelExporter))
		}
	}
	tracerProviderOptions = append(tracerProviderOptions, trace.WithResource(newResource()))

	return tracerProviderOptions, nil
}

func newStdOutExporter() (trace.SpanExporter, error) {
	return stdouttrace.New(
		// Use human readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

func newOTelExporter(env envConfig) (trace.SpanExporter, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, env.OTelCollectorURL, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector at %s: %w", env.OTelCollectorURL, err)
	}
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}
	return traceExporter, nil
}

func newResource() *resource.Resource {
	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.TelemetrySDKLanguageGo,
		semconv.ServiceNameKey.String("keptn-lifecycle-operator"),
		semconv.ServiceVersionKey.String(buildVersion+"-"+gitCommit+"-"+buildTime),
	)
	return r
}

func serveMetrics(collector prometheus.Collector) {
	registry := prometheus.NewRegistry()
	err := registry.Register(collector)
	if err != nil {
		fmt.Printf("error registering collector: %v", err)
		return
	}

	log.Printf("serving metrics at localhost:2222/metrics")
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	err = http.ListenAndServe(":2222", nil)
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}
