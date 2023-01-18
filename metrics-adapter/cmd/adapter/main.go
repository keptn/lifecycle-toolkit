package main

import (
	"flag"
	"fmt"
	keptnprovider "github.com/keptn/lifecycle-toolkit/metrics-adapter/pkg/provider"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"log"

	//metricsv1alpha1 "github.com/keptn/lifecycle-toolkit/operator/api/metrics/v1alpha1"
	"github.com/prometheus/client_golang/prometheus"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/component-base/logs"
	"k8s.io/klog/v2"
	"os"

	basecmd "sigs.k8s.io/custom-metrics-apiserver/pkg/cmd"
	"sigs.k8s.io/custom-metrics-apiserver/pkg/provider"
	"strings"
	"unicode"
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

	recordMetrics()

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

func watchMetrics() {
	stopCh := make(chan struct{})
	// Create a new Kubernetes dynamic client
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Error creating client config: %v", err)
	}
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Error creating client: %v", err)
	}

	// Define the resource we want to watch (e.g. pods)
	resource := schema.GroupVersionResource{Group: "metrics.keptn.sh", Version: "v1alpha1", Resource: "keptnmetrics"}

	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(client, 0, "keptn-lifecycle-toolkit-system", func(options *metav1.ListOptions) {})

	informer := factory.ForResource(resource).Informer()

	handlers := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {

		},
		UpdateFunc: func(oldObj, obj interface{}) {

		},
		DeleteFunc: func(obj interface{}) {
		},
	}
	informer.AddEventHandler(handlers)
	informer.Run(stopCh)

	// Watcher
	// Set up a watcher for the resource
	//watcher, err := client.Resource(resource).Watch(metav1.ListOptions{})
	//if err != nil {
	//	log.Fatalf("Error setting up watcher: %v", err)
	//}
	//
	//// Loop through events
	//for event := range watcher.ResultChan() {
	//	switch event.Type {
	//	case watch.Added:
	//		log.Printf("Pod added: %v", event.Object)
	//	case watch.Modified:
	//		log.Printf("Pod modified: %v", event.Object)
	//	case watch.Deleted:
	//		log.Printf("Pod deleted: %v", event.Object)
	//	}
	//}
}

func recordMetrics() {
	//go func() {
	//	scheme := runtime.NewScheme()
	//	if err := metricsv1alpha1.AddToScheme(scheme); err != nil {
	//		fmt.Println("failed to add metrics to scheme: " + err.Error())
	//	}
	//
	//	cl, err := ctrlclient.New(config.GetConfigOrDie(), ctrlclient.Options{Scheme: scheme})
	//	if err != nil {
	//		fmt.Println("failed to create client")
	//		os.Exit(1)
	//	}
	//
	//	for {
	//		list := metricsv1alpha1.MetricList{}
	//		err := cl.List(context.Background(), &list)
	//		if err != nil {
	//			fmt.Println("failed to list metrics" + err.Error())
	//		}
	//		for _, metric := range list.Items {
	//			normName := CleanUpString(metric.Name)
	//			if _, ok := m.gauges[normName]; !ok {
	//				m.gauges[normName] = prometheus.NewGauge(prometheus.GaugeOpts{
	//					Name: normName,
	//					Help: metric.Name,
	//				})
	//				prometheus.MustRegister(m.gauges[normName])
	//			}
	//			val, _ := strconv.ParseFloat(metric.Status.Value, 64)
	//			m.gauges[normName].Set(val)
	//		}
	//	}
	//}()
}

func CleanUpString(s string) string {
	return strings.Join(strings.FieldsFunc(s, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) }), "_")
}
