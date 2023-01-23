package adapter

import (
	"context"
	"flag"
	"fmt"
	"os"

	kmprovider "github.com/keptn/lifecycle-toolkit/operator/cmd/metrics/adapter/provider"
	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"
	basecmd "sigs.k8s.io/custom-metrics-apiserver/pkg/cmd"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider"
)

type MetricsAdapter struct {
	basecmd.AdapterBase
}

// RunAdapter starts the Keptn Metrics adapter to provide KeptnMetrics via the Kubernetes Custom Metrics API.
// Runs until the given context is done.
func (a *MetricsAdapter) RunAdapter(ctx context.Context) {

	logs.InitLogs()
	defer logs.FlushLogs()

	fmt.Println("Starting Keptn Metrics Adapter")
	// initialize the flags, with one custom flag for the message
	cmd := &MetricsAdapter{}
	// make sure you get the klog flags
	logs.AddGoFlags(flag.CommandLine)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	if err := cmd.Flags().Parse(os.Args); err != nil {
		klog.Fatalf("Could not parse flags: %v", err)
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

	return kmprovider.NewProvider(ctx, client)
}
