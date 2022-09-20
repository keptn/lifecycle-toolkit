package klcpermit

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	clientcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/record"
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
	recorder   record.EventRecorder
	svcManager *ServiceManager
}

var _ framework.PermitPlugin = &Permit{}

// Name returns name of the plugin.
func (pl *Permit) Name() string {
	return Name
}

func (pl *Permit) Permit(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (*framework.Status, time.Duration) {

	klog.InfoS("[Keptn Permit Plugin] waiting for pre-deployment checks on", p.GetObjectMeta().GetName())

	pl.recorder.Event(p, v1.EventTypeNormal, "SomeReason", "Waiting Pre-Deployment")

	switch pl.svcManager.Permit(ctx, p) {

	case Wait:
		klog.InfoS("[Keptn Permit Plugin] waiting for pre-deployment checks on", p.GetObjectMeta().GetName())
		return framework.NewStatus(framework.Wait), 5 * time.Second
	case Failure:
		klog.InfoS("[Keptn Permit Plugin] failed pre-deployment checks on", p.GetObjectMeta().GetName())
		return framework.NewStatus(framework.Error), 0 * time.Second
	case Success:
		klog.InfoS("[Keptn Permit Plugin] passed pre-deployment checks on", p.GetObjectMeta().GetName())
		return framework.NewStatus(framework.Success), 0 * time.Second
	default:
		klog.InfoS("[Keptn Permit Plugin] unknown status of pre-deployment checks for", p.GetObjectMeta().GetName())
		return framework.NewStatus(framework.Wait), 5 * time.Second //TODO what makes sense here?
	}

}

// New initializes a new plugin and returns it.
func New(_ runtime.Object, h framework.Handle) (framework.Plugin, error) {
	client, err := newClient()
	if err != nil {
		return nil, err
	}

	return &Permit{
		recorder:   setupRecorder(),
		svcManager: NewServiceManager(client),
		handler:    h,
	}, nil
}

func setupRecorder() record.EventRecorder {

	var config *rest.Config
	config, _ = rest.InClusterConfig()

	clientset, _ := kubernetes.NewForConfig(config)
	eventSink := &clientcorev1.EventSinkImpl{
		Interface: clientset.CoreV1().Events(v1.NamespaceAll),
	}
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartRecordingToSink(eventSink)
	r := eventBroadcaster.NewRecorder(scheme.Scheme, v1.EventSource{Component: "keptn-scheduler"})
	return r
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
