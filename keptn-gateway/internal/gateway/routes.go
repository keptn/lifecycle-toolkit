package gateway

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"
)

type RouteConfig struct {
	client *kubernetes.Clientset
	logger *zap.Logger
}

func RegisterRoutes(client *kubernetes.Clientset) http.Handler {
	rc := RouteConfig{
		client: client,
		logger: zap.Must(zap.NewDevelopment()),
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/namespace/:namespace/keptnapp/:keptnapp/resources", rc.ResourcePanelHandler)
	r.GET("/namespace/:namespace/keptnapp/:keptnapp/health", rc.HealthHandler)

	return r
}

func (rc *RouteConfig) ResourcePanelHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	keptnapp := c.Param("keptnapp")

	resp := make(map[string]string)
	resp["namespace"] = namespace
	resp["keptnapp"] = keptnapp

	c.JSON(http.StatusOK, resp)
}

func (rc *RouteConfig) HealthHandler(c *gin.Context) {
	_ = c.Param("namespace")
	_ = c.Param("keptnapp")

	// Demo implementation
	responses := []string{"Healthy", "Unhealthy", "Warning", "Progressing"}
	c.String(http.StatusOK, responses[rand.Intn(len(responses))])
}
