package common

import (
	"context"
	"fmt"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	apicommon "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	controllercommon "github.com/keptn/lifecycle-toolkit/operator/controllers/common"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/unit"
	"go.opentelemetry.io/otel/sdk/metric"
	otelsdk "go.opentelemetry.io/otel/sdk/trace"
	sdktest "go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func InitKeptnMeters() apicommon.KeptnMeters {
	provider := metric.NewMeterProvider()
	meter := provider.Meter("keptn/task")
	deploymentCount, _ := meter.SyncInt64().Counter("keptn.deployment.count", instrument.WithDescription("a simple counter for Keptn Deployments"))
	deploymentDuration, _ := meter.SyncFloat64().Histogram("keptn.deployment.duration", instrument.WithDescription("a histogram of duration for Keptn Deployments"), instrument.WithUnit(unit.Unit("s")))
	taskCount, _ := meter.SyncInt64().Counter("keptn.task.count", instrument.WithDescription("a simple counter for Keptn Tasks"))
	taskDuration, _ := meter.SyncFloat64().Histogram("keptn.task.duration", instrument.WithDescription("a histogram of duration for Keptn Tasks"), instrument.WithUnit(unit.Unit("s")))
	appCount, _ := meter.SyncInt64().Counter("keptn.app.count", instrument.WithDescription("a simple counter for Keptn Apps"))
	appDuration, _ := meter.SyncFloat64().Histogram("keptn.app.duration", instrument.WithDescription("a histogram of duration for Keptn Apps"), instrument.WithUnit(unit.Unit("s")))
	evaluationCount, _ := meter.SyncInt64().Counter("keptn.evaluation.count", instrument.WithDescription("a simple counter for Keptn Evaluations"))
	evaluationDuration, _ := meter.SyncFloat64().Histogram("keptn.evaluation.duration", instrument.WithDescription("a histogram of duration for Keptn Evaluations"), instrument.WithUnit(unit.Unit("s")))

	meters := apicommon.KeptnMeters{
		TaskCount:          taskCount,
		TaskDuration:       taskDuration,
		DeploymentCount:    deploymentCount,
		DeploymentDuration: deploymentDuration,
		AppCount:           appCount,
		AppDuration:        appDuration,
		EvaluationCount:    evaluationCount,
		EvaluationDuration: evaluationDuration,
	}
	return meters
}

type TracerFactory struct {
	Tracer trace.TracerProvider
}

func (f *TracerFactory) GetTracer(name string) controllercommon.ITracer {
	return f.Tracer.Tracer(name)
}

func IgnoreAlreadyExists(err error) error {
	if apierrors.IsAlreadyExists(err) {
		return nil
	}
	return err
}

func LogErrorIfPresent(err error) {
	if err != nil {
		GinkgoLogr.Error(err, "Something went wrong while cleaning up the test environment")
	}
}

func ResetSpanRecords(tp *otelsdk.TracerProvider, spanRecorder *sdktest.SpanRecorder) {
	GinkgoLogr.Info("Removing ", fmt.Sprint(len(spanRecorder.Ended())), " spans")
	tp.UnregisterSpanProcessor(spanRecorder)
	spanRecorder = sdktest.NewSpanRecorder()
	tp.RegisterSpanProcessor(spanRecorder)
}

func MakeKLTDefaultNamespace(k8sClient client.Client, name string) *v1.Namespace {
	ns := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}

	_ = k8sClient.Create(context.TODO(), ns)

	return ns
}

func DeleteAppInCluster(ctx context.Context, k8sClient client.Client, instance *klcv1alpha3.KeptnApp) {
	By("Cleaning Up KeptnApp CRD ")
	err := k8sClient.Delete(ctx, instance)
	LogErrorIfPresent(err)
}

func AssertResourceUpdated(ctx context.Context, k8sClient client.Client, instance *klcv1alpha3.KeptnApp) *klcv1alpha3.KeptnAppVersion {

	appVersion := GetAppVersion(ctx, k8sClient, instance)

	By("Comparing expected app version")
	Expect(appVersion.Spec.AppName).To(Equal(instance.Name))
	Expect(appVersion.Spec.Version).To(Equal(instance.Spec.Version))
	Expect(appVersion.Spec.Workloads).To(Equal(instance.Spec.Workloads))

	return appVersion
}

func GetAppVersion(ctx context.Context, k8sClient client.Client, instance *klcv1alpha3.KeptnApp) *klcv1alpha3.KeptnAppVersion {
	appvName := types.NamespacedName{
		Namespace: instance.Namespace,
		Name:      fmt.Sprintf("%s-%s-%d", instance.Name, instance.Spec.Version, instance.Generation),
	}

	appVersion := &klcv1alpha3.KeptnAppVersion{}
	By("Retrieving Created app version")
	Eventually(func() error {
		return k8sClient.Get(ctx, appvName, appVersion)
	}, "20s").Should(Succeed())

	return appVersion
}
