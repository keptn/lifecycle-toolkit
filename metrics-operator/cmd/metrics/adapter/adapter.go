package adapter

import (
	"context"
	"flag"
	"fmt"

	kmprovider "github.com/keptn/lifecycle-toolkit/metrics-operator/cmd/metrics/adapter/provider"
	"k8s.io/apiserver/pkg/server/options"
	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"
	basecmd "sigs.k8s.io/custom-metrics-apiserver/pkg/cmd"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider"
)

const (
	flagPort                   = "adapter-port"
	flagCertificateDirectory   = "adapter-certs-dir"
	defaultCertificatePairName = "apiserver"
)

var (
	port    int
	certDir string
)

type MetricsAdapter struct {
	basecmd.AdapterBase
	KeptnNamespace string
}

// RunAdapter starts the Keptn Metrics adapter to provide KeptnMetrics via the Kubernetes Custom Metrics API.
// Runs until the given context is done.
func (a *MetricsAdapter) RunAdapter(ctx context.Context) {

	logs.InitLogs()
	defer logs.FlushLogs()

	addFlags()

	fmt.Println("Starting Keptn Metrics Adapter")
	// initialize the flags, with one custom flag for the message
	cmd := &MetricsAdapter{}
	// make sure you get the klog flags
	logs.AddGoFlags(flag.CommandLine)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	if err := cmd.Flags().Parse([]string{}); err != nil {
		klog.Fatalf("Could not parse flags: %v", err)
	}

	cmd.CustomMetricsAdapterServerOptions.SecureServing.BindPort = port
	cmd.CustomMetricsAdapterServerOptions.SecureServing.ServerCert = options.GeneratableKeyCert{
		PairName:      defaultCertificatePairName,
		CertDirectory: certDir,
	}

	prov := cmd.makeProviderOrDie(ctx)

	cmd.WithCustomMetrics(prov)

	if err := cmd.Run(ctx.Done()); err != nil {
		klog.Fatalf("Could not run custom metrics adapter: %v", err)
	}
	klog.Info("Finishing Keptn Metrics Adapter")
}

func (a *MetricsAdapter) makeProviderOrDie(ctx context.Context) provider.CustomMetricsProvider {
	client, err := a.DynamicClient()
	if err != nil {
		klog.Fatalf("unable to construct dynamic client: %v", err)
	}

	return kmprovider.NewProvider(ctx, client, a.KeptnNamespace)
}

func addFlags() {
	flag.IntVar(&port, flagPort, 6443, "Port of the metrics adapter endpoint")
	flag.StringVar(&certDir, flagCertificateDirectory, "/tmp/metrics-adapter/serving-certs", "Directory in which to look for certificates for the Metrics Adapter.")
	flag.Parse()
}
