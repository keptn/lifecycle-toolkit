package main

import (
	"flag"
	"fmt"
	keptnprovider "github.com/keptn/lifecycle-toolkit/metrics-adapter/pkg/provider"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/klog/v2"
	"net/http"
	"os"

	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/component-base/logs"
	basecmd "sigs.k8s.io/custom-metrics-apiserver/pkg/cmd"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider"
)

type KeptnAdapter struct {
	basecmd.AdapterBase

	// the message printed on startup
	Message string
}

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	go serveMetrics()

	fmt.Println("Starting Keptn Metrics Adapter")
	// initialize the flags, with one custom flag for the message
	cmd := &KeptnAdapter{}
	cmd.Flags().StringVar(&cmd.Message, "msg", "starting adapter...", "startup message")
	// make sure you get the klog flags
	logs.AddGoFlags(flag.CommandLine)
	cmd.Flags().AddGoFlagSet(flag.CommandLine)
	cmd.Flags().Parse(os.Args)

	prov := cmd.makeProviderOrDie()

	cmd.WithCustomMetrics(prov)
	// you could also set up external metrics support,
	// if your provider supported it:
	// cmd.WithExternalMetrics(provider)

	klog.Infof(cmd.Message)
	if err := cmd.Run(wait.NeverStop); err != nil {
		klog.Fatalf("unable to run custom metrics adapter: %v", err)
	}
	fmt.Println("Finishing Keptn Metrics Adapter")
}

func (a *KeptnAdapter) makeProviderOrDie() provider.CustomMetricsProvider {
	client, err := a.DynamicClient()
	if err != nil {
		klog.Fatalf("unable to construct dynamic client: %v", err)
	}

	mapper, err := a.RESTMapper()
	if err != nil {
		klog.Fatalf("unable to construct discovery REST mapper: %v", err)
	}

	return keptnprovider.NewProvider(client, mapper)
}

func serveMetrics() {
	klog.Infof("serving metrics at localhost:9999/metrics")

	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":9999", nil)
	if err != nil {
		fmt.Printf("error serving http: %v", err)
		return
	}
}
