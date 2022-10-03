package klcpermit

import (
	"context"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

// Name is the name of the plugin used in the plugin registry and configurations.
const (
	Name = "KLCPermit"
)

// Permit is a plugin that implements a wait for pre-deployment checks
type Permit struct {
	handler         framework.Handle
	workloadManager *WorkloadManager
}

var _ framework.PermitPlugin = &Permit{}

// Name returns name of the plugin.
func (pl *Permit) Name() string {
	return Name
}

func (pl *Permit) Permit(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (*framework.Status, time.Duration) {

	klog.InfoS("[Keptn Permit Plugin] waiting for pre-deployment checks on", p.GetObjectMeta().GetName())

	switch pl.workloadManager.Permit(ctx, p) {

	case Wait:
		klog.Infof("[Keptn Permit Plugin] waiting for pre-deployment checks on", p.GetObjectMeta().GetName())
		return framework.NewStatus(framework.Wait), 30 * time.Second
	case Failure:
		klog.Infof("[Keptn Permit Plugin] failed pre-deployment checks on", p.GetObjectMeta().GetName())
		return framework.NewStatus(framework.Error), 0 * time.Second
	case Success:
		klog.Infof("[Keptn Permit Plugin] passed pre-deployment checks on", p.GetObjectMeta().GetName())
		return framework.NewStatus(framework.Success), 0 * time.Second
	default:
		klog.Infof("[Keptn Permit Plugin] unknown status of pre-deployment checks for", p.GetObjectMeta().GetName())
		return framework.NewStatus(framework.Wait), 30 * time.Second
	}

}

// New initializes a new plugin and returns it.
func New(_ runtime.Object, h framework.Handle) (framework.Plugin, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}

	return &Permit{
		workloadManager: NewWorkloadManager(client),
		handler:         h,
	}, nil
}

func newClient() (dynamic.Interface, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return dynClient, nil
}
