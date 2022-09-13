package klcpermit

import (
	"context"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
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
	Name    = "KLCPermit"
	Wait    = "Wait"
	Success = "Success"
)

// Permit is a plugin that implements a wait for pre-deployment checks
type Permit struct {
	recorder record.EventRecorder
}

var _ framework.PermitPlugin = &Permit{}

// Name returns name of the plugin.
func (pl *Permit) Name() string {
	return Name
}

func (pl *Permit) Permit(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (*framework.Status, time.Duration) {

	klog.InfoS("[Keptn Permit Plugin] waiting for pre-deployment checks on", klog.KObj(p.GetObjectMeta()))

	pl.recorder.Event(p, v1.EventTypeNormal, "SomeReason", "Waiting Pre-Deployment")

	retStatus := framework.NewStatus(framework.Success)
	waitTime := 0 * time.Second
	return retStatus, waitTime
}

// New initializes a new plugin and returns it.
func New(_ runtime.Object, _ framework.Handle) (framework.Plugin, error) {
	return &Permit{
		recorder: setupRecorder(),
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
