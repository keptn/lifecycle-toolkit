package webhook

import (
	"flag"

	"go.opentelemetry.io/otel"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"github.com/keptn/lifecycle-toolkit/operator/cmd/certificates"
	"github.com/keptn/lifecycle-toolkit/operator/cmd/config"
	cmdManager "github.com/keptn/lifecycle-toolkit/operator/cmd/manager"
	"github.com/keptn/lifecycle-toolkit/operator/webhooks"
	"github.com/keptn/lifecycle-toolkit/operator/webhooks/pod_mutator"
	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	FlagCertificateDirectory   = "certs-dir"
	FlagCertificateFileName    = "cert"
	FlagCertificateKeyFileName = "cert-key"
)

var (
	certificateDirectory   string
	certificateFileName    string
	certificateKeyFileName string
)

type Builder struct {
	configProvider  config.Provider
	managerProvider cmdManager.Provider
	namespace       string
	podName         string
}

func NewWebhookBuilder() Builder {
	return Builder{}
}

func (builder Builder) SetConfigProvider(provider config.Provider) Builder {
	builder.configProvider = provider
	return builder
}

func (builder Builder) SetManagerProvider(provider cmdManager.Provider) Builder {
	builder.managerProvider = provider
	return builder
}

func (builder Builder) SetNamespace(namespace string) Builder {
	builder.namespace = namespace
	return builder
}

func (builder Builder) SetPodName(podName string) Builder {
	builder.podName = podName
	return builder
}

func (builder Builder) GetManagerProvider() cmdManager.Provider {
	if builder.managerProvider == nil {
		builder.managerProvider = NewWebhookManagerProvider(certificateDirectory, certificateKeyFileName, certificateFileName)
	}

	return builder.managerProvider
}

func (builder Builder) Run(webhookManager manager.Manager) error {

	addFlags()
	builder.GetManagerProvider().SetupWebhookServer(webhookManager)

	certificates.
		NewCertificateWatcher(webhookManager, builder.namespace, webhooks.SecretCertsName, ctrl.Log.WithName("Webhook Cert Manager")).
		WaitForCertificates()

	webhookManager.GetWebhookServer().Register("/mutate-v1-pod", &webhook.Admission{
		Handler: &pod_mutator.PodMutatingWebhook{
			Client:   webhookManager.GetClient(),
			Tracer:   otel.Tracer("keptn/webhook"),
			Recorder: webhookManager.GetEventRecorderFor("keptn/webhook"),
			Log:      ctrl.Log.WithName("Mutating Webhook"),
		}})

	signalHandler := ctrl.SetupSignalHandler()
	err := webhookManager.Start(signalHandler)
	return errors.WithStack(err)
}

func addFlags() {
	flag.StringVar(&certificateDirectory, FlagCertificateDirectory, "/tmp/webhook/certs", "Directory to look certificates for.")
	flag.StringVar(&certificateFileName, FlagCertificateFileName, "tls.crt", "File name for the public certificate.")
	flag.StringVar(&certificateKeyFileName, FlagCertificateKeyFileName, "tls.key", "File name for the private key.")
	flag.Parse()
}
