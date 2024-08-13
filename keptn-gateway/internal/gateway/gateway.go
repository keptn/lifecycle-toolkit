package gateway

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Gateway struct {
	port   int
	logger *zap.Logger
	client *kubernetes.Clientset
}

func NewGateway() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 8080
	}

	NewGateway := &Gateway{
		port:   port,
		logger: zap.Must(zap.NewDevelopment()),
	}
	defer NewGateway.logger.Sync()

	// kubernetes clientset the clientset
	/*
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
	*/

	// alternative not in cluster config
	home := homedir.HomeDir()
	kubeconfig := filepath.Join(home, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil
	}

	NewGateway.client, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewGateway.port),
		Handler:      NewGateway.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	NewGateway.logger.Info("Listening on: " + server.Addr)

	return server
}
