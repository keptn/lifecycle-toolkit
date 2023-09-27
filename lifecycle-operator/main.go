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
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	argov1alpha1 "github.com/argoproj/argo-rollouts/pkg/apis/rollouts/v1alpha1"
	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/keptn/lifecycle-toolkit/klt-cert-manager/pkg/certificates"
	certCommon "github.com/keptn/lifecycle-toolkit/klt-cert-manager/pkg/common"
	"github.com/keptn/lifecycle-toolkit/klt-cert-manager/pkg/webhook"
	lifecyclev1alpha1 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha1"
	lifecyclev1alpha2 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha2"
	lifecyclev1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	optionsv1alpha1 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/options/v1alpha1"
	controllercommon "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/config"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/common/telemetry"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/keptnapp"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/keptnappcreationrequest"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/keptnappversion"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/keptnevaluation"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/keptntask"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/keptntaskdefinition"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/keptnworkload"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/keptnworkloadinstance"
	controlleroptions "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/options"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/webhooks/pod_mutator"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
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
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
	ctrlWebhook "sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(lifecyclev1alpha1.AddToScheme(scheme))
	utilruntime.Must(lifecyclev1alpha2.AddToScheme(scheme))
	utilruntime.Must(optionsv1alpha1.AddToScheme(scheme))
	utilruntime.Must(lifecyclev1alpha3.AddToScheme(scheme))
	utilruntime.Must(argov1alpha1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

type envConfig struct {
	PodNamespace string `envconfig:"POD_NAMESPACE" default:""`
	PodName      string `envconfig:"POD_NAME" default:""`

	KeptnAppControllerLogLevel                int `envconfig:"KEPTN_APP_CONTROLLER_LOG_LEVEL" default:"0"`
	KeptnAppCreationRequestControllerLogLevel int `envconfig:"KEPTN_APP_CREATION_REQUEST_CONTROLLER_LOG_LEVEL" default:"0"`
	KeptnAppVersionControllerLogLevel         int `envconfig:"KEPTN_APP_VERSION_CONTROLLER_LOG_LEVEL" default:"0"`
	KeptnEvaluationControllerLogLevel         int `envconfig:"KEPTN_EVALUATION_CONTROLLER_LOG_LEVEL" default:"0"`
	KeptnTaskControllerLogLevel               int `envconfig:"KEPTN_TASK_CONTROLLER_LOG_LEVEL" default:"0"`
	KeptnTaskDefinitionControllerLogLevel     int `envconfig:"KEPTN_TASK_DEFINITION_CONTROLLER_LOG_LEVEL" default:"0"`
	KeptnWorkloadControllerLogLevel           int `envconfig:"KEPTN_WORKLOAD_CONTROLLER_LOG_LEVEL" default:"0"`
	KeptnWorkloadInstanceControllerLogLevel   int `envconfig:"KEPTN_WORKLOAD_INSTANCE_CONTROLLER_LOG_LEVEL" default:"0"`
	KeptnOptionsControllerLogLevel            int `envconfig:"OPTIONS_CONTROLLER_LOG_LEVEL" default:"0"`

	SchedulingGatesEnabled bool `envconfig:"SCHEDULING_GATES_ENABLED" default:"false"`

	KeptnOptionsCollectorURL string `envconfig:"OTEL_COLLECTOR_URL" default:""`
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
	keptnMeters := telemetry.SetUpKeptnTaskMeters(meter)

	// Start the prometheus HTTP server and pass the exporter Collector to it
	go serveMetrics()

	// As recommended by the kubebuilder docs, webhook registration should be disabled if running locally. See https://book.kubebuilder.io/cronjob-tutorial/running.html#running-webhooks-locally for reference
	flag.BoolVar(&disableWebhook, "disable-webhook", false, "Disable the registration of webhooks.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	// parse the flags, so we ensure they can be set to something else than their default values
	flag.Parse()

	// inject pod namespace into common configs
	config.Instance().SetDefaultNamespace(env.PodNamespace)

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	disableCacheFor := []ctrlclient.Object{&corev1.Secret{}}
	opt := ctrl.Options{
		Scheme: scheme,
		Metrics: server.Options{
			BindAddress: metricsAddr,
		},
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "6b866dd9.keptn.sh",
		Client: ctrlclient.Options{
			Cache: &ctrlclient.CacheOptions{
				DisableFor: disableCacheFor,
			},
		},
	}

	var webhookBuilder webhook.Builder
	if !disableWebhook {
		webhookBuilder = webhook.NewWebhookServerBuilder().
			LoadCertOptionsFromFlag().
			SetPort(9443).
			SetNamespace(env.PodNamespace).
			SetPodName(env.PodName)
		opt.WebhookServer = webhookBuilder.GetWebhookServer()
	}

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), opt)
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	// Enabling OTel
	err = telemetry.GetOtelInstance().InitOtelCollector("")
	if err != nil {
		setupLog.Error(err, "unable to initialize OTel tracer options")
	}

	spanHandler := &telemetry.SpanHandler{}

	// create Cloud Event client
	ceClient, err := ce.NewClientHTTP()
	if err != nil {
		setupLog.Error(err, "failed to create CloudEvent client")
		os.Exit(1)
	}

	taskLogger := ctrl.Log.WithName("KeptnTask Controller").V(env.KeptnTaskControllerLogLevel)
	taskRecorder := mgr.GetEventRecorderFor("keptntask-controller")
	taskReconciler := &keptntask.KeptnTaskReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		Log:           taskLogger,
		EventSender:   controllercommon.NewEventMultiplexer(taskLogger, taskRecorder, ceClient),
		Meters:        keptnMeters,
		TracerFactory: telemetry.GetOtelInstance(),
	}
	if err = (taskReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnTask")
		os.Exit(1)
	}

	taskDefinitionLogger := ctrl.Log.WithName("KeptnTaskDefinition Controller").V(env.KeptnTaskDefinitionControllerLogLevel)
	taskDefinitionRecorder := mgr.GetEventRecorderFor("keptntaskdefinition-controller")
	taskDefinitionReconciler := &keptntaskdefinition.KeptnTaskDefinitionReconciler{
		Client:      mgr.GetClient(),
		Scheme:      mgr.GetScheme(),
		Log:         taskDefinitionLogger,
		EventSender: controllercommon.NewEventMultiplexer(taskDefinitionLogger, taskDefinitionRecorder, ceClient),
	}
	if err = (taskDefinitionReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnTaskDefinition")
		os.Exit(1)
	}

	appLogger := ctrl.Log.WithName("KeptnApp Controller").V(env.KeptnAppControllerLogLevel)
	appRecorder := mgr.GetEventRecorderFor("keptnapp-controller")
	appReconciler := &keptnapp.KeptnAppReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		Log:           appLogger,
		EventSender:   controllercommon.NewEventMultiplexer(appLogger, appRecorder, ceClient),
		TracerFactory: telemetry.GetOtelInstance(),
	}
	if err = (appReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnApp")
		os.Exit(1)
	}

	appCreationRequestLogger := ctrl.Log.WithName("KeptnAppCreationRequest Controller")
	appCreationRequestReconciler := keptnappcreationrequest.NewReconciler(
		mgr.GetClient(),
		mgr.GetScheme(),
		appCreationRequestLogger.V(env.KeptnAppCreationRequestControllerLogLevel),
	)
	if err := appCreationRequestReconciler.SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnAppCreationRequest")
		os.Exit(1)
	}

	workloadLogger := ctrl.Log.WithName("KeptnWorkload Controller").V(env.KeptnWorkloadControllerLogLevel)
	workloadRecorder := mgr.GetEventRecorderFor("keptnworkload-controller")
	workloadReconciler := &keptnworkload.KeptnWorkloadReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		Log:           workloadLogger,
		EventSender:   controllercommon.NewEventMultiplexer(workloadLogger, workloadRecorder, ceClient),
		TracerFactory: telemetry.GetOtelInstance(),
	}
	if err = (workloadReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnWorkload")
		os.Exit(1)
	}

	workloadInstanceLogger := ctrl.Log.WithName("KeptnWorkloadInstance Controller").V(env.KeptnWorkloadInstanceControllerLogLevel)
	workloadInstanceRecorder := mgr.GetEventRecorderFor("keptnworkloadinstance-controller")
	workloadInstanceReconciler := &keptnworkloadinstance.KeptnWorkloadInstanceReconciler{
		SchedulingGatesHandler: controllercommon.NewSchedulingGatesHandler(mgr.GetClient(), workloadInstanceLogger, env.SchedulingGatesEnabled),
		Client:                 mgr.GetClient(),
		Scheme:                 mgr.GetScheme(),
		Log:                    workloadInstanceLogger,
		EventSender:            controllercommon.NewEventMultiplexer(workloadInstanceLogger, workloadInstanceRecorder, ceClient),
		Meters:                 keptnMeters,
		TracerFactory:          telemetry.GetOtelInstance(),
		SpanHandler:            spanHandler,
	}
	if err = (workloadInstanceReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnWorkloadInstance")
		os.Exit(1)
	}

	appVersionLogger := ctrl.Log.WithName("KeptnAppVersion Controller").V(env.KeptnAppVersionControllerLogLevel)
	appVersionRecorder := mgr.GetEventRecorderFor("keptnappversion-controller")
	appVersionReconciler := &keptnappversion.KeptnAppVersionReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		Log:           appVersionLogger,
		EventSender:   controllercommon.NewEventMultiplexer(appVersionLogger, appVersionRecorder, ceClient),
		TracerFactory: telemetry.GetOtelInstance(),
		Meters:        keptnMeters,
		SpanHandler:   spanHandler,
	}
	if err = (appVersionReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnAppVersion")
		os.Exit(1)
	}

	evaluationLogger := ctrl.Log.WithName("KeptnEvaluation Controller").V(env.KeptnEvaluationControllerLogLevel)
	evaluationRecorder := mgr.GetEventRecorderFor("keptnevaluation-controller")
	evaluationReconciler := &keptnevaluation.KeptnEvaluationReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		Log:           evaluationLogger,
		EventSender:   controllercommon.NewEventMultiplexer(evaluationLogger, evaluationRecorder, ceClient),
		TracerFactory: telemetry.GetOtelInstance(),
		Meters:        keptnMeters,
	}
	if err = (evaluationReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnEvaluation")
		os.Exit(1)
	}

	configLogger := ctrl.Log.WithName("KeptnConfig Controller").V(env.KeptnOptionsControllerLogLevel)
	configReconciler := &controlleroptions.KeptnConfigReconciler{
		Client:              mgr.GetClient(),
		Scheme:              mgr.GetScheme(),
		Log:                 configLogger,
		DefaultCollectorURL: env.KeptnOptionsCollectorURL,
	}
	if err = (configReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnConfig")
		os.Exit(1)
	}

	if err = (&lifecyclev1alpha3.KeptnApp{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "KeptnApp")
		os.Exit(1)
	}
	if err = (&lifecyclev1alpha3.KeptnEvaluationProvider{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "KeptnEvaluationProvider")
		os.Exit(1)
	}
	if err = (&lifecyclev1alpha3.KeptnAppVersion{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "KeptnAppVersion")
		os.Exit(1)
	}
	if err = (&lifecyclev1alpha3.KeptnWorkloadInstance{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "KeptnWorkloadInstance")
		os.Exit(1)
	}
	if err = (&lifecyclev1alpha3.KeptnTaskDefinition{}).SetupWebhookWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create webhook", "webhook", "KeptnTaskDefinition")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	telemetry.SetUpKeptnMeters(meter, mgr.GetClient())

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}
	if !disableWebhook {
		webhookBuilder = webhookBuilder.SetCertificateWatcher(
			certificates.NewCertificateWatcher(
				mgr.GetAPIReader(),
				webhookBuilder.GetOptions().CertDir,
				env.PodNamespace,
				certCommon.SecretName,
				setupLog,
			))

		setupLog.Info(fmt.Sprintf("%v", webhookBuilder))
		webhookLogger := ctrl.Log.WithName("Mutating Webhook")
		webhookRecorder := mgr.GetEventRecorderFor("keptn/webhook")
		webhookBuilder.Register(mgr, map[string]*ctrlWebhook.Admission{
			"/mutate-v1-pod": {
				Handler: &pod_mutator.PodMutatingWebhook{
					SchedulingGatesEnabled: env.SchedulingGatesEnabled,
					Client:                 mgr.GetClient(),
					Tracer:                 otel.Tracer("keptn/webhook"),
					EventSender:            controllercommon.NewEventMultiplexer(webhookLogger, webhookRecorder, ceClient),
					Decoder:                admission.NewDecoder(mgr.GetScheme()),
					Log:                    webhookLogger,
				},
			},
		})
		setupLog.Info("starting webhook")
	}
	setupLog.Info("starting manager")
	setupLog.Info("Keptn lifecycle-operator is alive")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
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
