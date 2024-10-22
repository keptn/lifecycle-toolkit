package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"k8s.io/client-go/kubernetes"

	lifecycle_v1 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
)

type RouteConfig struct {
	client *kubernetes.Clientset
	logger *zap.Logger
}

func (rc *RouteConfig) getKeptnAppVersion(namespace, keptnappname string) (keptnapp lifecycle_v1.KeptnAppVersion, err error) {
	response, err := rc.client.RESTClient().
		Get().
		AbsPath("/apis/lifecycle.keptn.sh/v1").
		Resource("keptnappversions").
		Namespace(namespace).
		DoRaw(context.TODO())

	if err != nil {
		return lifecycle_v1.KeptnAppVersion{}, err
	}

	l := lifecycle_v1.KeptnAppVersionList{}
	err = json.Unmarshal([]byte(response), &l)

	if err != nil {
		return lifecycle_v1.KeptnAppVersion{}, err
	}

	for _, keptnapp := range l.Items {
		if keptnapp.Spec.AppName == keptnappname {
			return keptnapp, nil
		}
	}
	return lifecycle_v1.KeptnAppVersion{}, errors.New("KeptnApp Not found")
}

func (rc *RouteConfig) getKeptnWorkloadVersions(namespace string, wrefs []lifecycle_v1.KeptnWorkloadRef) (workloads []lifecycle_v1.KeptnWorkloadVersion, err error) {
	response, err := rc.client.RESTClient().
		Get().
		AbsPath("/apis/lifecycle.keptn.sh/v1").
		Resource("keptnworkloadversions").
		Namespace(namespace).
		DoRaw(context.TODO())

	if err != nil {
		return nil, err
	}

	l := lifecycle_v1.KeptnWorkloadVersionList{}
	err = json.Unmarshal([]byte(response), &l)
	if err != nil {
		return nil, err
	}

	wls := make([]lifecycle_v1.KeptnWorkloadVersion, len(wrefs))
	for i, wl := range l.Items {
		for _, wref := range wrefs {
			if wl.Spec.WorkloadName == wl.Spec.AppName+"-"+wref.Name {
				wls[i] = wl
				break
			}
		}
	}

	return wls, nil
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
	keptnappname := c.Param("keptnapp")

	keptnapp, err := rc.getKeptnAppVersion(namespace, keptnappname)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	workloads, err := rc.getKeptnWorkloadVersions(namespace, keptnapp.Spec.Workloads)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	type ReturnVal struct {
		AppVersion lifecycle_v1.KeptnAppVersion        `json:"appversion"`
		Workloads  []lifecycle_v1.KeptnWorkloadVersion `json:"workloads"`
	}

	r := ReturnVal{
		AppVersion: keptnapp,
		Workloads:  workloads,
	}
	c.JSON(http.StatusOK, r)
}

func (rc *RouteConfig) HealthHandler(c *gin.Context) {
	namespace := c.Param("namespace")
	keptnappname := c.Param("keptnapp")

	keptnapp, err := rc.getKeptnAppVersion(namespace, keptnappname)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.String(http.StatusOK, string(keptnapp.Status.Status))
}
