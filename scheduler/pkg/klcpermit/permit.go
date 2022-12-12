package klcpermit

import (
	"context"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
	klog.Infof("[Keptn Permit Plugin] waiting for pre-deployment checks on %s", p.GetObjectMeta().GetName())

	pl.workloadManager.ObserveWorkloadForPod(ctx, p)

	return framework.NewStatus(framework.Wait), 5 * time.Minute
}

// New initializes a new plugin and returns it.
func New(_ runtime.Object, h framework.Handle) (framework.Plugin, error) {
	workloadManager, err := NewWorkloadManager(h)
	return &Permit{
		workloadManager: workloadManager,
		handler:         h,
	}, err
}
