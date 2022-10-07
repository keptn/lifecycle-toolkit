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

// PluginName is the name of the plugin used in the plugin registry and configurations.
const (
	PluginName = "KLCPermit"
)

// Permit is a plugin that waits for pre-deployment checks
type Permit struct {
	handler         framework.Handle
	workloadManager *WorkloadManager
}

var _ framework.PermitPlugin = &Permit{}

// PluginName returns name of the plugin.
func (pl *Permit) Name() string {
	return PluginName
}

func (pl *Permit) Permit(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (*framework.Status, time.Duration) {

	klog.Infof("[Keptn Permit Plugin] waiting for pre-deployment checks on %s", p.GetObjectMeta().GetName())

	// check the permit immediately, to fail early in case the pod cannot be queued
	switch pl.workloadManager.Permit(ctx, p) {

	case Wait:
		klog.Infof("[Keptn Permit Plugin] waiting for pre-deployment checks on %s", p.GetObjectMeta().GetName())
		go pl.monitorPod(ctx, p)
		return framework.NewStatus(framework.Wait), 5 * time.Minute
	case Failure:
		klog.Infof("[Keptn Permit Plugin] failed pre-deployment checks on %s", p.GetObjectMeta().GetName())
		return framework.NewStatus(framework.Error), 0 * time.Second
	case Success:
		klog.Infof("[Keptn Permit Plugin] passed pre-deployment checks on %s", p.GetObjectMeta().GetName())
		return framework.NewStatus(framework.Success), 0 * time.Second
	default:
		klog.Infof("[Keptn Permit Plugin] unknown status of pre-deployment checks for %s", p.GetObjectMeta().GetName())
		return framework.NewStatus(framework.Wait), 5 * time.Minute
	}

}

func (pl *Permit) monitorPod(ctx context.Context, p *v1.Pod) {
	waitingPodHandler := pl.handler.GetWaitingPod(p.UID)

	switch pl.workloadManager.Permit(ctx, p) {
	case Failure:
		waitingPodHandler.Reject(PluginName, "Pre Deployment Check failed")
	case Success:
		waitingPodHandler.Allow(PluginName)
	default:
		time.Sleep(10 * time.Second)
		pl.monitorPod(ctx, p)
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
