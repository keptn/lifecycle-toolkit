package klcpermit

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"time"
)

// Name is the name of the plugin used in the plugin registry and configurations.
const (
	Name    = "KLCPermit"
	Wait    = "Wait"
	Success = "Success"
)

// Permit is a plugin that implements a wait for pre-deployment checks
type Permit struct {
}

var _ framework.PermitPlugin = &Permit{}

// Name returns name of the plugin.
func (pl *Permit) Name() string {
	return Name
}

func (pl *Permit) Permit(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (*framework.Status, time.Duration) {

	klog.Info("[Keptn Permit Plugin] waiting for pre-deployment checks")
	retStatus := framework.NewStatus(framework.Success)
	waitTime := 0 * time.Second
	return retStatus, waitTime
}

func (pl *Permit) decideState() string {
	return Success
}

// New initializes a new plugin and returns it.
func New(_ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
	return &Permit{}, nil
}
