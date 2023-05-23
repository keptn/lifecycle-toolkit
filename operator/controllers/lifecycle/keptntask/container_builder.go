package keptntask

import (
	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"golang.org/x/net/context"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/client-go/tools/record"

	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// IContainerBuilder is the interface that describes the operations needed to help build job specs of a task
type IContainerBuilder interface {
	// AddContainers populates a job containers and volumes based on the task definition spec
	AddContainers(ctx context.Context, job *batchv1.Job) error
}

// BuilderOptions contains everything needed to build the current job
type BuilderOptions struct {
	client.Client
	recorder record.EventRecorder
	req      ctrl.Request
	Log      logr.Logger
	task     *klcv1alpha3.KeptnTask
	taskDef  *klcv1alpha3.KeptnTaskDefinition
}

func getContainerBuilder(options BuilderOptions) IContainerBuilder {
	if !reflect.DeepEqual(options.taskDef.Spec.Function, klcv1alpha3.FunctionSpec{}) {
		builder := newJSBuilder(options)
		return &builder
	}
	return nil
}
