package webhook

import (
	"flag"

	"github.com/keptn/lifecycle-toolkit/klt-cert-manager/pkg/certificates"
	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
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
	managerProvider    Provider
	namespace          string
	podName            string
	certificateWatcher certificates.ICertificateWatcher
}

func NewWebhookBuilder() Builder {
	return Builder{}
}

func (builder Builder) SetManagerProvider(provider Provider) Builder {
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

func (builder Builder) SetCertificateWatcher(watcher certificates.ICertificateWatcher) Builder {
	builder.certificateWatcher = watcher
	return builder
}

func (builder Builder) GetManagerProvider() Provider {
	if builder.managerProvider == nil {
		builder.managerProvider = NewWebhookManagerProvider(certificateDirectory, certificateKeyFileName, certificateFileName)
	}

	return builder.managerProvider
}

// Run ensures that the secret containing the certificate required for the webhooks is available, and then registers the
// given webhooks at the webhookManager's webhook server.
func (builder Builder) Run(webhookManager manager.Manager, webhooks map[string]*webhook.Admission) error {

	addFlags()
	builder.GetManagerProvider().SetupWebhookServer(webhookManager)

	builder.certificateWatcher.WaitForCertificates()

	for path, admissionWebhook := range webhooks {
		webhookManager.GetWebhookServer().Register(path, admissionWebhook)
	}

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
