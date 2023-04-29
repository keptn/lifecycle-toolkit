package webhook

import (
	"flag"

	"github.com/keptn/lifecycle-toolkit/metrics-operator/cmd/certificates"
	"github.com/keptn/lifecycle-toolkit/metrics-operator/cmd/config"
	cmdManager "github.com/keptn/lifecycle-toolkit/metrics-operator/cmd/manager"
	"github.com/pkg/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

const (
	flagCertificateDirectory   = "certs-dir"
	flagCertificateFileName    = "cert"
	flagCertificateKeyFileName = "cert-key"
	secretCertsName            = "klt-certs"
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
		NewCertificateWatcher(webhookManager.GetAPIReader(), webhookManager.GetWebhookServer().CertDir, builder.namespace, secretCertsName, ctrl.Log.WithName("Webhook Cert Manager")).
		WaitForCertificates()

	signalHandler := ctrl.SetupSignalHandler()
	err := webhookManager.Start(signalHandler)
	return errors.WithStack(err)
}

func addFlags() {
	flag.StringVar(&certificateDirectory, flagCertificateDirectory, "/tmp/webhook/certs", "Directory to look certificates for.")
	flag.StringVar(&certificateFileName, flagCertificateFileName, "tls.crt", "File name for the public certificate.")
	flag.StringVar(&certificateKeyFileName, flagCertificateKeyFileName, "tls.key", "File name for the private key.")
	flag.Parse()
}
