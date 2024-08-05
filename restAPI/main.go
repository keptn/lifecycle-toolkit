package main

import (
	"context"
	"encoding/json"
	"flag"
	"path/filepath"

	"fmt"
	"log"
	"net/http"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	metrics_v1 "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1"
)

var clientset *kubernetes.Clientset

func handlerMetricsCount(w http.ResponseWriter, r *http.Request) {
	log.Println("http call to /pods")

	d, err := clientset.RESTClient().
		Get().
		AbsPath("/apis/metrics.keptn.sh/v1/").
		Resource("keptnmetrics").
		Namespace("default").
		DoRaw(context.TODO())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	l := metrics_v1.KeptnMetricList{}

	err = json.Unmarshal([]byte(d), &l)

	w.WriteHeader(http.StatusOK)
	returnval := fmt.Sprintf("%#v", len(l.Items))
	w.Write([]byte(returnval))
}

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kc", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kc", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	clientConfig, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)

	if err != nil {
		panic(err.Error())
	}

	clientset, err = kubernetes.NewForConfig(clientConfig)

	if err != nil {
		panic(err)
	}

	http.Handle("/metricscount", http.HandlerFunc(handlerMetricsCount))
	log.Fatal(http.ListenAndServe(":8080", nil))
	log.Println("Listening on port 8080")
}
