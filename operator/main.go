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
	"time"

	"github.com/kelseyhightower/envconfig"
	lifecyclev1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha1"
	lifecyclev1alpha2 "github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2"
	"github.com/keptn/lifecycle-toolkit/operator/api/v1alpha2/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/keptnapp"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/keptnappversion"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/keptnevaluation"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/keptntask"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/keptntaskdefinition"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/keptnworkload"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/keptnworkloadinstance"
	"github.com/keptn/lifecycle-toolkit/operator/webhooks"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
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
	utilruntime.Must(lifecyclev1alpha2.AddToScheme(scheme))
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
	// TODO uncomment after bumping OTel version to 1.11.2 or higher
	// exporter, err := otelprom.New()
	// if err != nil {
	// 	setupLog.Error(err, "unable to start OTel")
	// }
	provider := metric.NewMeterProvider(metric.WithReader(exporter))
	meter := provider.Meter("keptn/task")
	deploymentCount, err := meter.SyncInt64().Counter("keptn.deployment.count", instrument.WithDescription("a simple counter for Keptn Deployments"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	deploymentDuration, err := meter.SyncFloat64().Histogram("keptn.deployment.duration", instrument.WithDescription("a histogram of duration for Keptn Deployments"), instrument.WithUnit(unit.Unit("s")))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	deploymentActiveGauge, err := meter.AsyncInt64().Gauge("keptn.deployment.active", instrument.WithDescription("a gauge keeping track of the currently active Keptn Deployments"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	taskCount, err := meter.SyncInt64().Counter("keptn.task.count", instrument.WithDescription("a simple counter for Keptn Tasks"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	taskDuration, err := meter.SyncFloat64().Histogram("keptn.task.duration", instrument.WithDescription("a histogram of duration for Keptn Tasks"), instrument.WithUnit(unit.Unit("s")))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	taskActiveGauge, err := meter.AsyncInt64().Gauge("keptn.task.active", instrument.WithDescription("a simple counter of active Keptn Tasks"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	appCount, err := meter.SyncInt64().Counter("keptn.app.count", instrument.WithDescription("a simple counter for Keptn Apps"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	appDuration, err := meter.SyncFloat64().Histogram("keptn.app.duration", instrument.WithDescription("a histogram of duration for Keptn Apps"), instrument.WithUnit(unit.Unit("s")))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	appActiveGauge, err := meter.AsyncInt64().Gauge("keptn.app.active", instrument.WithDescription("a simple counter of active Keptn Apps"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	evaluationCount, err := meter.SyncInt64().Counter("keptn.evaluation.count", instrument.WithDescription("a simple counter for Keptn Evaluations"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	evaluationDuration, err := meter.SyncFloat64().Histogram("keptn.evaluation.duration", instrument.WithDescription("a histogram of duration for Keptn Evaluations"), instrument.WithUnit(unit.Unit("s")))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	evaluationActiveGauge, err := meter.AsyncInt64().Gauge("keptn.evaluation.active", instrument.WithDescription("a simple counter of active Keptn Evaluations"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}
	appDeploymentIntervalGauge, err := meter.AsyncFloat64().Gauge("keptn.app.deploymentinterval", instrument.WithDescription("a gauge of the interval between deployments"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}

	appDeploymentDurationGauge, err := meter.AsyncFloat64().Gauge("keptn.app.deploymentduration", instrument.WithDescription("a gauge of the duration of deployments"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}

	workloadDeploymentIntervalGauge, err := meter.AsyncFloat64().Gauge("keptn.deployment.deploymentinterval", instrument.WithDescription("a gauge of the interval between deployments"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}

	workloadDeploymentDurationGauge, err := meter.AsyncFloat64().Gauge("keptn.deployment.deploymentduration", instrument.WithDescription("a gauge of the duration of deployments"))
	if err != nil {
		setupLog.Error(err, "unable to start OTel")
	}

	meters := common.KeptnMeters{
		TaskCount:          taskCount,
		TaskDuration:       taskDuration,
		DeploymentCount:    deploymentCount,
		DeploymentDuration: deploymentDuration,
		AppCount:           appCount,
		AppDuration:        appDuration,
		EvaluationCount:    evaluationCount,
		EvaluationDuration: evaluationDuration,
	}

	// Start the prometheus HTTP server and pass the exporter Collector to it
	// TODO uncomment after bumping OTel version to 1.11.2 or higher
	//go serveMetrics()
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

	spanHandler := &controllercommon.SpanHandler{}

	if !disableWebhook {
		mgr.GetWebhookServer().Register("/mutate-v1-pod", &webhook.Admission{
			Handler: &webhooks.PodMutatingWebhook{
				Client:   mgr.GetClient(),
				Tracer:   otel.Tracer("keptn/webhook"),
				Recorder: mgr.GetEventRecorderFor("keptn/webhook"),
				Log:      ctrl.Log.WithName("Mutating Webhook"),
			}})
	}
	taskReconciler := &keptntask.KeptnTaskReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Log:      ctrl.Log.WithName("KeptnTask Controller"),
		Recorder: mgr.GetEventRecorderFor("keptntask-controller"),
		Meters:   meters,
		Tracer:   otel.Tracer("keptn/operator/task"),
	}
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
	if err = (taskDefinitionReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnTaskDefinition")
		os.Exit(1)
	}

	appReconciler := &keptnapp.KeptnAppReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Log:      ctrl.Log.WithName("KeptnApp Controller"),
		Recorder: mgr.GetEventRecorderFor("keptnapp-controller"),
		Tracer:   otel.Tracer("keptn/operator/app"),
	}
	if err = (appReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnApp")
		os.Exit(1)
	}

	workloadReconciler := &keptnworkload.KeptnWorkloadReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Log:      ctrl.Log.WithName("KeptnWorkload Controller"),
		Recorder: mgr.GetEventRecorderFor("keptnworkload-controller"),
		Tracer:   otel.Tracer("keptn/operator/workload"),
	}
	if err = (workloadReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnWorkload")
		os.Exit(1)
	}

	workloadInstanceReconciler := &keptnworkloadinstance.KeptnWorkloadInstanceReconciler{
		Client:      mgr.GetClient(),
		Scheme:      mgr.GetScheme(),
		Log:         ctrl.Log.WithName("KeptnWorkloadInstance Controller"),
		Recorder:    mgr.GetEventRecorderFor("keptnworkloadinstance-controller"),
		Meters:      meters,
		Tracer:      otel.Tracer("keptn/operator/workloadinstance"),
		SpanHandler: spanHandler,
	}
	if err = (workloadInstanceReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnWorkloadInstance")
		os.Exit(1)
	}

	appVersionReconciler := &keptnappversion.KeptnAppVersionReconciler{
		Client:      mgr.GetClient(),
		Scheme:      mgr.GetScheme(),
		Log:         ctrl.Log.WithName("KeptnAppVersion Controller"),
		Recorder:    mgr.GetEventRecorderFor("keptnappversion-controller"),
		Tracer:      otel.Tracer("keptn/operator/appversion"),
		Meters:      meters,
		SpanHandler: spanHandler,
	}
	if err = (appVersionReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnAppVersion")
		os.Exit(1)
	}

	evaluationReconciler := &keptnevaluation.KeptnEvaluationReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Log:      ctrl.Log.WithName("KeptnEvaluation Controller"),
		Recorder: mgr.GetEventRecorderFor("keptnevaluation-controller"),
		Tracer:   otel.Tracer("keptn/operator/evaluation"),
		Meters:   meters,
	}
	if err = (evaluationReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnEvaluation")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	err = meter.RegisterCallback(
		[]instrument.Asynchronous{
			deploymentActiveGauge,
			taskActiveGauge,
			appActiveGauge,
			evaluationActiveGauge,
			appDeploymentIntervalGauge,
			appDeploymentDurationGauge,
			workloadDeploymentIntervalGauge,
			workloadDeploymentDurationGauge,
		},
		func(ctx context.Context) {
			activeDeployments, err := controllercommon.GetActiveInstances(ctx, mgr.GetClient(), &lifecyclev1alpha2.KeptnWorkloadInstanceList{})
			if err != nil {
				setupLog.Error(err, "unable to gather active deployments")
			}
			for _, val := range activeDeployments {
				deploymentActiveGauge.Observe(ctx, val.Value, val.Attributes...)
			}

			activeApps, err := controllercommon.GetActiveInstances(ctx, mgr.GetClient(), &lifecyclev1alpha2.KeptnAppVersionList{})
			if err != nil {
				setupLog.Error(err, "unable to gather active apps")
			}
			for _, val := range activeApps {
				appActiveGauge.Observe(ctx, val.Value, val.Attributes...)
			}

			activeTasks, err := controllercommon.GetActiveInstances(ctx, mgr.GetClient(), &lifecyclev1alpha2.KeptnTaskList{})
			if err != nil {
				setupLog.Error(err, "unable to gather active tasks")
			}
			for _, val := range activeTasks {
				taskActiveGauge.Observe(ctx, val.Value, val.Attributes...)
			}

			activeEvaluations, err := controllercommon.GetActiveInstances(ctx, mgr.GetClient(), &lifecyclev1alpha2.KeptnEvaluationList{})
			if err != nil {
				setupLog.Error(err, "unable to gather active evaluations")
			}
			for _, val := range activeEvaluations {
				evaluationActiveGauge.Observe(ctx, val.Value, val.Attributes...)
			}

			appDeploymentInterval, err := controllercommon.GetDeploymentInterval(ctx, mgr.GetClient(), &lifecyclev1alpha2.KeptnAppVersionList{}, &lifecyclev1alpha2.KeptnAppVersion{})
			if err != nil {
				setupLog.Error(err, "unable to gather app deployment intervals")
			}
			for _, val := range appDeploymentInterval {
				appDeploymentIntervalGauge.Observe(ctx, val.Value, val.Attributes...)
			}

			appDeploymentDuration, err := controllercommon.GetDeploymentDuration(ctx, mgr.GetClient(), &lifecyclev1alpha2.KeptnAppVersionList{})
			if err != nil {
				setupLog.Error(err, "unable to gather app deployment durations")
			}
			for _, val := range appDeploymentDuration {
				appDeploymentDurationGauge.Observe(ctx, val.Value, val.Attributes...)
			}

			workloadDeploymentInterval, err := controllercommon.GetDeploymentInterval(ctx, mgr.GetClient(), &lifecyclev1alpha2.KeptnWorkloadInstanceList{}, &lifecyclev1alpha2.KeptnWorkloadInstance{})
			if err != nil {
				setupLog.Error(err, "unable to gather workload deployment intervals")
			}
			for _, val := range workloadDeploymentInterval {
				workloadDeploymentIntervalGauge.Observe(ctx, val.Value, val.Attributes...)
			}

			workloadDeploymentDuration, err := controllercommon.GetDeploymentDuration(ctx, mgr.GetClient(), &lifecyclev1alpha2.KeptnWorkloadInstanceList{})
			if err != nil {
				setupLog.Error(err, "unable to gather workload deployment durations")
			}
			for _, val := range workloadDeploymentDuration {
				workloadDeploymentDurationGauge.Observe(ctx, val.Value, val.Attributes...)
			}

		})
	if err != nil {
		fmt.Println("Failed to register callback")
		panic(err)
	}

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
