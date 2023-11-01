package klcpermit

import (
	"context"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

// PluginName is the name of the plugin used in the plugin registry and configurations.
const (
	PluginName = "KLCPermit"
)

// Permit is a plugin that waits for pre-deployment checks to be successfully finished
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

	klog.Infof("[Keptn Permit Plugin] waiting for pre-deployment checks on %s in namespace %s", p.GetObjectMeta().GetName(), p.GetObjectMeta().GetNamespace())

	// check the permit immediately, to fail early in case the pod cannot be queued
	switch pl.workloadManager.Permit(ctx, p) {

	case Success:
		klog.Infof("[Keptn Permit Plugin] passed pre-deployment checks on %s in namespace %s", p.GetObjectMeta().GetName(), p.GetObjectMeta().GetNamespace())
		return framework.NewStatus(framework.Success), 0 * time.Second
	default:
		klog.Infof("[Keptn Permit Plugin] waiting for pre-deployment checks on %s in namespace %s to succeed", p.GetObjectMeta().GetName(), p.GetObjectMeta().GetNamespace())
		go func() {
			// create a new context since we are in a new goroutine
			ctx2, cancel := context.WithCancel(context.Background())
			defer cancel()
			pl.monitorPod(ctx2, p)
		}()
		return framework.NewStatus(framework.Wait), 5 * time.Minute
	}

}

func (pl *Permit) monitorPod(ctx context.Context, p *v1.Pod) {
	waitingPodHandler := pl.handler.GetWaitingPod(p.UID)

	for {
		switch pl.workloadManager.Permit(ctx, p) {
		case Success:
			waitingPodHandler.Allow(PluginName)
			return
		default:
			time.Sleep(10 * time.Second)
		}
	}

}

// New initializes a new plugin and returns it.
func New(_ runtime.Object, h framework.Handle) (framework.Plugin, error) {
	client, err := newClient(h)
	if err != nil {
		return nil, err
	}

	return &Permit{
		workloadManager: NewWorkloadManager(client),
		handler:         h,
	}, nil
}

func newClient(handle framework.Handle) (dynamic.Interface, error) {
	config := handle.KubeConfig()

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return dynClient, nil
}
