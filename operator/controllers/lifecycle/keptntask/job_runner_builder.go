package keptntask

import (
	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"golang.org/x/net/context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// JobRunnerBuilder is the interface that describes the operations needed to help build job specs of a task
type JobRunnerBuilder interface {
	// CreateContainerWithVolumes returns a job container and volumes based on the task definition spec
	CreateContainerWithVolumes(ctx context.Context) (*corev1.Container, []corev1.Volume, error)
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

func getJobRunnerBuilder(options BuilderOptions) JobRunnerBuilder {
	if options.taskDef.IsJSSpecDefined() {
		return newJSBuilder(options)
	}
	if options.taskDef.IsContainerSpecDefined() {
		return newContainerBuilder(options.taskDef)
	}
	return nil
}
