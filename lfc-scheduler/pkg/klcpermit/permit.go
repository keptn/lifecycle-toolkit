package klcpermit

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ref "k8s.io/client-go/tools/reference"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	"time"
)

// Name is the name of the plugin used in the plugin registry and configurations.
const (
	Name = "KLCPermit"
)

// Permit is a plugin that implements a wait for pre-deployment checks
type Permit struct {
	handler    framework.Handle
	svcManager *ServiceManager
}

var _ framework.PermitPlugin = &Permit{}

// Name returns name of the plugin.
func (pl *Permit) Name() string {
	return Name
}

func (pl *Permit) Permit(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (*framework.Status, time.Duration) {

	klog.InfoS("[Keptn Permit Plugin] waiting for pre-deployment checks on", p.GetObjectMeta().GetName())

	err := pl.sendEvent(p, "test", "started", "Waiting for pre-deployment checks")
	if err != nil {
		return framework.NewStatus(framework.Error), 0 * time.Second
	}

	switch pl.svcManager.Permit(ctx, p) {

	case Wait:
		klog.InfoS("[Keptn Permit Plugin] waiting for pre-deployment checks on", p.GetObjectMeta().GetName())
		return framework.NewStatus(framework.Wait), 30 * time.Second
	case Failure:
		klog.InfoS("[Keptn Permit Plugin] failed pre-deployment checks on", p.GetObjectMeta().GetName())
		return framework.NewStatus(framework.Error), 0 * time.Second
	case Success:
		klog.InfoS("[Keptn Permit Plugin] passed pre-deployment checks on", p.GetObjectMeta().GetName())
		return framework.NewStatus(framework.Success), 0 * time.Second
	default:
		klog.InfoS("[Keptn Permit Plugin] unknown status of pre-deployment checks for", p.GetObjectMeta().GetName())
		return framework.NewStatus(framework.Wait), 30 * time.Second //TODO what makes sense here?
	}

}

func (pl *Permit) sendEvent(p *v1.Pod, reason string, action string, note string) error {
	regarding, err := ref.GetPartialReference(scheme.Scheme, p, ".spec.containers[1]")
	if err != nil {
		return err
	}

	related, err := ref.GetPartialReference(scheme.Scheme, p, ".spec.containers[0]")
	if err != nil {
		return err
	}
	pl.handler.EventRecorder().Eventf(regarding, related, v1.EventTypeNormal, reason, action, note)
	return nil
}

// New initializes a new plugin and returns it.
func New(_ runtime.Object, h framework.Handle) (framework.Plugin, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}

	return &Permit{
		svcManager: NewServiceManager(client),
		handler:    h,
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
