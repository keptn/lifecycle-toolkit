package webhook

import (
	"flag"

	"github.com/keptn/lifecycle-toolkit/klt-cert-manager/pkg/certificates"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

const (
	FlagCertificateDirectory   = "certs-dir"
	FlagCertificateFileName    = "cert"
	FlagCertificateKeyFileName = "cert-key"
)

//go:generate moq -pkg fake -skip-ensure -out ../fake/manager_mock.go . IManager:MockManager
type IManager manager.Manager

type Builder struct {
	namespace          string
	podName            string
	certificateWatcher certificates.ICertificateWatcher
	options            webhook.Options
}

func NewWebhookServerBuilder() Builder {
	b := Builder{}
	return b
}

func (builder Builder) SetNamespace(namespace string) Builder {
	builder.namespace = namespace
	return builder
}

func (builder Builder) SetPort(port int) Builder {
	builder.options.Port = port
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

func (builder Builder) LoadCertOptionsFromFlag() Builder {
	flag.StringVar(&builder.options.CertDir, FlagCertificateDirectory, "/tmp/webhook/certs", "Directory to look certificates for.")
	flag.StringVar(&builder.options.CertName, FlagCertificateFileName, "tls.crt", "File name for the public certificate.")
	flag.StringVar(&builder.options.KeyName, FlagCertificateKeyFileName, "tls.key", "File name for the private key.")
	flag.Parse()
	return builder
}

func (builder Builder) GetWebhookServer() webhook.Server {
	return webhook.NewServer(builder.options)
}

func (builder Builder) GetOptions() webhook.Options {
	return builder.options
}

// Register ensures that the secret containing the certificate required for the webhooks is available, and then registers the
// given webhooks at the webhookManager's webhook server.
func (builder Builder) Register(webhookManager manager.Manager, webhooks map[string]*webhook.Admission) {

	builder.certificateWatcher.WaitForCertificates()

	for path, admissionWebhook := range webhooks {
		webhookManager.GetWebhookServer().Register(path, admissionWebhook)
	}

}
