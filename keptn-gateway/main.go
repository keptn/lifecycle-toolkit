package main

import (
	"context"
	"encoding/json"

	"fmt"
	"log"
	"net/http"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"go.uber.org/zap"

	metrics_v1 "github.com/keptn/lifecycle-toolkit/metrics-operator/api/v1"
)

type Gateway struct {
	clientset *kubernetes.Clientset
	logger    *zap.Logger
}

var gateway Gateway

func handlerMetricsCount(w http.ResponseWriter, r *http.Request) {
	log.Println("http call to /pods")

	d, err := gateway.clientset.RESTClient().
		Get().
		AbsPath("/apis/metrics.keptn.sh/v1/").
		Resource("keptnmetrics").
		Namespace("default").
		DoRaw(context.TODO())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
		if err != nil {
			gateway.logger.Error(err.Error())
		}
	}

	l := metrics_v1.KeptnMetricList{}

	err = json.Unmarshal([]byte(d), &l)
	if err != nil {
		gateway.logger.Error(err.Error())
	}

	w.WriteHeader(http.StatusOK)
	returnval := fmt.Sprintf("%#v", len(l.Items))

	_, err = w.Write([]byte(returnval))
	if err != nil {
		gateway.logger.Error(err.Error())
	}
}

func main() {
	gateway.logger = zap.Must(zap.NewDevelopment())
	defer gateway.logger.Sync()

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	gateway.clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	http.Handle("/metricscount", http.HandlerFunc(handlerMetricsCount))
	log.Fatal(http.ListenAndServe(":8080", nil))
	log.Println("Listening on port 8080")
}
