package gateway

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Gateway) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/namespace/:namespace/keptnapp/:keptnapp/resources", s.ResourcePanelHandler)
	r.GET("/namespace/:namespace/keptnapp/:keptnapp/health", s.HealthHandler)

	return r
}

func (s *Gateway) ResourcePanelHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	keptnapp := c.Param("keptnapp")

	resp := make(map[string]string)
	resp["namespace"] = namespace
	resp["keptnapp"] = keptnapp

	c.JSON(http.StatusOK, resp)
}

func (s *Gateway) HealthHandler(c *gin.Context) {
	_ = c.Param("namespace")
	_ = c.Param("keptnapp")

	// Demo implementation
	responses := []string{"Healthy", "Unhealthy", "Warning", "Progressing"}
	c.String(http.StatusOK, responses[rand.Intn(len(responses))])
}
