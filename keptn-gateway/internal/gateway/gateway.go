package gateway

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type Gateway struct {
	port   int
	logger *zap.Logger
	client *kubernetes.Clientset
	server *http.Server
}

func NewGateway() *Gateway {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 8080
	}

	// kubernetes clientset in cluster
	config, err := rest.InClusterConfig()
	if err != nil {
		// if error fallback to .kube/config
		home := homedir.HomeDir()
		kubeconfig := filepath.Join(home, ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	NewGateway := &Gateway{
		port:   port,
		logger: zap.Must(zap.NewDevelopment()),
		server: &http.Server{
			Addr:         fmt.Sprintf(":%d", port),
			Handler:      RegisterRoutes(client),
			IdleTimeout:  time.Minute,
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 30 * time.Second,
		},
	}
	defer NewGateway.logger.Sync()

	return NewGateway
}

func (gw *Gateway) Serve() error {
	go func() {
		gw.logger.Info("server starting")
		if err := gw.server.ListenAndServe(); err != http.ErrServerClosed {
			panic(err)
		}
	}()

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c // block until signal received

	gw.logger.Info("shutting down")
	gw.server.Shutdown(nil)

	return nil
}
