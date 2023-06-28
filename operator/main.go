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
	"github.com/kelseyhightower/envconfig"
	"github.com/keptn/lifecycle-toolkit/klt-cert-manager/pkg/certificates"
	certCommon "github.com/keptn/lifecycle-toolkit/klt-cert-manager/pkg/common"
	"github.com/keptn/lifecycle-toolkit/klt-cert-manager/pkg/webhook"
	metricsapi "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1alpha3"
	lifecyclev1alpha1 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha1"
	lifecyclev1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	lifecyclev1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	optionsv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/apis/options/v1alpha1"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnapp"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnappcreationrequest"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnappversion"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnevaluation"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptntask"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptntaskdefinition"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnworkload"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/keptnworkloadinstance"
	controlleroptions "github.com/keptn/lifecycle-toolkit/operator/controllers/options"
	"github.com/keptn/lifecycle-toolkit/operator/webhooks/pod_mutator"
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
	utilruntime.Must(metricsapi.AddToScheme(scheme))
	utilruntime.Must(argov1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
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
	keptnMeters := controllercommon.SetUpKeptnTaskMeters(meter)

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

	// Enabling OTel
	err = controllercommon.GetOtelInstance().InitOtelCollector("")
	if err != nil {
		setupLog.Error(err, "unable to initialize OTel tracer options")
	}

	spanHandler := &controllercommon.SpanHandler{}

	taskLogger := ctrl.Log.WithName("KeptnTask Controller")
	taskReconciler := &keptntask.KeptnTaskReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		Log:           taskLogger.V(env.KeptnTaskControllerLogLevel),
		Recorder:      mgr.GetEventRecorderFor("keptntask-controller"),
		Meters:        keptnMeters,
		TracerFactory: controllercommon.GetOtelInstance(),
	}
	if err = (taskReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnTask")
		os.Exit(1)
	}

	taskDefinitionLogger := ctrl.Log.WithName("KeptnTaskDefinition Controller")
	taskDefinitionReconciler := &keptntaskdefinition.KeptnTaskDefinitionReconciler{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Log:      taskDefinitionLogger.V(env.KeptnTaskDefinitionControllerLogLevel),
		Recorder: mgr.GetEventRecorderFor("keptntaskdefinition-controller"),
	}
	if err = (taskDefinitionReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnTaskDefinition")
		os.Exit(1)
	}

	appLogger := ctrl.Log.WithName("KeptnApp Controller")
	appReconciler := &keptnapp.KeptnAppReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		Log:           appLogger.V(env.KeptnAppControllerLogLevel),
		Recorder:      mgr.GetEventRecorderFor("keptnapp-controller"),
		TracerFactory: controllercommon.GetOtelInstance(),
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

	workloadLogger := ctrl.Log.WithName("KeptnWorkload Controller")
	workloadReconciler := &keptnworkload.KeptnWorkloadReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		Log:           workloadLogger.V(env.KeptnWorkloadControllerLogLevel),
		Recorder:      mgr.GetEventRecorderFor("keptnworkload-controller"),
		TracerFactory: controllercommon.GetOtelInstance(),
	}
	if err = (workloadReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnWorkload")
		os.Exit(1)
	}

	workloadInstanceLogger := ctrl.Log.WithName("KeptnWorkloadInstance Controller")
	workloadInstanceReconciler := &keptnworkloadinstance.KeptnWorkloadInstanceReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		Log:           workloadInstanceLogger.V(env.KeptnWorkloadInstanceControllerLogLevel),
		Recorder:      mgr.GetEventRecorderFor("keptnworkloadinstance-controller"),
		Meters:        keptnMeters,
		TracerFactory: controllercommon.GetOtelInstance(),
		SpanHandler:   spanHandler,
	}
	if err = (workloadInstanceReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnWorkloadInstance")
		os.Exit(1)
	}

	appVersionLogger := ctrl.Log.WithName("KeptnAppVersion Controller")
	appVersionReconciler := &keptnappversion.KeptnAppVersionReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		Log:           appVersionLogger.V(env.KeptnAppVersionControllerLogLevel),
		Recorder:      mgr.GetEventRecorderFor("keptnappversion-controller"),
		TracerFactory: controllercommon.GetOtelInstance(),
		Meters:        keptnMeters,
		SpanHandler:   spanHandler,
	}
	if err = (appVersionReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnAppVersion")
		os.Exit(1)
	}

	evaluationLogger := ctrl.Log.WithName("KeptnEvaluation Controller")
	evaluationReconciler := &keptnevaluation.KeptnEvaluationReconciler{
		Client:        mgr.GetClient(),
		Scheme:        mgr.GetScheme(),
		Log:           evaluationLogger.V(env.KeptnEvaluationControllerLogLevel),
		Recorder:      mgr.GetEventRecorderFor("keptnevaluation-controller"),
		TracerFactory: controllercommon.GetOtelInstance(),
		Meters:        keptnMeters,
		Namespace:     env.PodNamespace,
	}
	if err = (evaluationReconciler).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "KeptnEvaluation")
		os.Exit(1)
	}

	configLogger := ctrl.Log.WithName("KeptnConfig Controller")
	configReconciler := &controlleroptions.KeptnConfigReconciler{
		Client:              mgr.GetClient(),
		Scheme:              mgr.GetScheme(),
		Log:                 configLogger.V(env.KeptnOptionsControllerLogLevel),
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
			SetManagerProvider(
				webhook.NewWebhookManagerProvider(
					mgr.GetWebhookServer().(*ctrlWebhook.DefaultServer).Options.CertDir, "tls.key", "tls.crt"),
			).
			SetCertificateWatcher(
				certificates.NewCertificateWatcher(
					mgr.GetAPIReader(),
					mgr.GetWebhookServer().(*ctrlWebhook.DefaultServer).Options.CertDir,
					env.PodNamespace,
					certCommon.SecretName,
					setupLog,
				),
			)

		setupLog.Info("starting webhook and manager")

		decoder, err := admission.NewDecoder(mgr.GetScheme())
		if err != nil {
			setupLog.Error(err, "unable to initialize decoder")
			os.Exit(1)
		}

		if err := webhookBuilder.Run(mgr, map[string]*ctrlWebhook.Admission{
			"/mutate-v1-pod": {
				Handler: &pod_mutator.PodMutatingWebhook{
					Client:   mgr.GetClient(),
					Tracer:   otel.Tracer("keptn/webhook"),
					Recorder: mgr.GetEventRecorderFor("keptn/webhook"),
					Decoder:  decoder,
					Log:      ctrl.Log.WithName("Mutating Webhook"),
				},
			},
		}); err != nil {
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
