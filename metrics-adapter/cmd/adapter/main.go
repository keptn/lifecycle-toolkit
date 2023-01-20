package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	keptnprovider "github.com/keptn/lifecycle-toolkit/metrics-adapter/pkg/provider"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"
	basecmd "sigs.k8s.io/custom-metrics-apiserver/pkg/cmd"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider"
)

type Metrics struct {
	gauges map[string]prometheus.Gauge
}

var m Metrics

type KeptnAdapter struct {
	basecmd.AdapterBase

	// the message printed on startup
	Message string
}

func main() {
	m.gauges = make(map[string]prometheus.Gauge)

	logs.InitLogs()
	defer logs.FlushLogs()

	fmt.Println("Starting Keptn Metrics Adapter")
	// initialize the flags, with one custom flag for the message
	cmd := &KeptnAdapter{}
	cmd.Flags().StringVar(&cmd.Message, "msg", "starting adapter...", "startup message")
	// make sure you get the klog flags
	logs.AddGoFlags(flag.CommandLine)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	if err := cmd.Flags().Parse(os.Args); err != nil {
		klog.Fatalf("Could not parse flags: %v", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM)
	defer cancel()

	prov := cmd.makeProviderOrDie(ctx)

	cmd.WithCustomMetrics(prov)

	klog.Infof(cmd.Message)
	if err := cmd.Run(wait.NeverStop); err != nil {
		klog.Fatalf("Could not run custom metrics adapter: %v", err)
	}
	klog.Info("Finishing Keptn Metrics Adapter")
}

func (a *KeptnAdapter) makeProviderOrDie(ctx context.Context) provider.CustomMetricsProvider {
	client, err := a.DynamicClient()
	if err != nil {
		klog.Fatalf("unable to construct dynamic client: %v", err)
	}

	mapper, err := a.RESTMapper()
	if err != nil {
		klog.Fatalf("unable to construct discovery REST mapper: %v", err)
	}

	return keptnprovider.NewProvider(ctx, client, mapper)
}
