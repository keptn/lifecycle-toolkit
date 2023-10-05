package telemetry

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	lifecyclev1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3/common"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/interfaces"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	logger       = ctrl.Log.WithName("otel-utils")
	gitCommit    string
	buildTime    string
	buildVersion string
	otelInitOnce sync.Once
)

type otelConfig struct {
	TracerProvider *trace.TracerProvider
	OtelExporter   *trace.SpanExporter

	lastAppliedCollectorURL string

	mtx     sync.RWMutex
	tracers map[string]ITracer
}

// do not export this type to make it accessible only via the GetInstance method (i.e Singleton)
var otelInstance *otelConfig

func GetOtelInstance() *otelConfig {
	// initialize once
	otelInitOnce.Do(func() {
		otelInstance = &otelConfig{
			tracers: map[string]ITracer{},
		}
	})

	return otelInstance
}

func (o *otelConfig) InitOtelCollector(otelCollectorUrl string) error {
	if o.lastAppliedCollectorURL == otelCollectorUrl {
		return nil
	}
	tpOptions, otelExporter, err := GetOTelTracerProviderOptions(otelCollectorUrl)
	if err != nil {
		return err
	}

	o.TracerProvider = trace.NewTracerProvider(tpOptions...)
	otel.SetTracerProvider(o.TracerProvider)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	o.OtelExporter = &otelExporter
	o.cleanTracers()
	o.lastAppliedCollectorURL = otelCollectorUrl
	logger.Info("Successfully initialized OTel collector")
	return nil
}

func (o *otelConfig) ShutDown() {
	if err := o.TracerProvider.Shutdown(context.Background()); err != nil {
		os.Exit(1)
	}
}

func (o *otelConfig) GetTracer(name string) ITracer {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	if o.tracers[name] == nil {
		o.tracers[name] = otel.Tracer(name)
	}
	return o.tracers[name]
}

func (o *otelConfig) cleanTracers() {
	o.mtx.Lock()
	defer o.mtx.Unlock()
	o.tracers = map[string]ITracer{}
}

func GetOTelTracerProviderOptions(oTelCollectorUrl string) ([]trace.TracerProviderOption, trace.SpanExporter, error) {
	var tracerProviderOptions []trace.TracerProviderOption
	var otelExporter trace.SpanExporter

	if oTelCollectorUrl != "" {
		// try to set OTel exporter for Jaeger
		otelExporter, err := newOTelExporter(oTelCollectorUrl)
		if err != nil {
			// log the error, but do not break if Jaeger exporter cannot be created
			logger.Error(err, "Could not set up OTel exporter")
			return nil, nil, err
		} else if otelExporter != nil {
			tracerProviderOptions = append(tracerProviderOptions, trace.WithBatcher(otelExporter))
		}
	} else {
		// if no collector is set, we use std::out to print trace info
		stdOutExp, err := newStdOutExporter()
		if err != nil {
			return nil, nil, fmt.Errorf("could not create stdout OTel exporter: %w", err)
		}
		tracerProviderOptions = append(tracerProviderOptions, trace.WithBatcher(stdOutExp))
	}
	tracerProviderOptions = append(tracerProviderOptions, trace.WithResource(newResource()))

	return tracerProviderOptions, otelExporter, nil
}

func newStdOutExporter() (trace.SpanExporter, error) {
	return stdouttrace.New(
		// Use human readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

func newOTelExporter(oTelCollectorUrl string) (trace.SpanExporter, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	conn, err := grpc.DialContext(ctx, oTelCollectorUrl, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector at %s: %w", oTelCollectorUrl, err)
	}
	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, fmt.Errorf("failed to create trace exporter: %w", err)
	}
	return traceExporter, nil
}

func newResource() *resource.Resource {
	r := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.TelemetrySDKLanguageGo,
		semconv.ServiceNameKey.String("lifecycle-operator"),
		semconv.ServiceVersionKey.String(buildVersion+"-"+gitCommit+"-"+buildTime),
	)
	return r
}

func SetUpKeptnMeters(meter interfaces.IMeter, mgr client.Client) {
	deploymentActiveGauge, err := meter.Int64ObservableGauge("keptn.deployment.active", metric.WithDescription("a gauge keeping track of the currently active Keptn Deployments"))
	if err != nil {
		logger.Error(err, "unable to initialize active deployments OTel gauge")
	}
	taskActiveGauge, err := meter.Int64ObservableGauge("keptn.task.active", metric.WithDescription("a simple counter of active Keptn Tasks"))
	if err != nil {
		logger.Error(err, "unable to initialize active tasks OTel gauge")
	}
	appActiveGauge, err := meter.Int64ObservableGauge("keptn.app.active", metric.WithDescription("a simple counter of active Keptn Apps"))
	if err != nil {
		logger.Error(err, "unable to initialize active apps OTel gauge")
	}
	evaluationActiveGauge, err := meter.Int64ObservableGauge("keptn.evaluation.active", metric.WithDescription("a simple counter of active Keptn Evaluations"))
	if err != nil {
		logger.Error(err, "unable to initialize active evaluations OTel gauge")
	}
	appDeploymentIntervalGauge, err := meter.Float64ObservableGauge("keptn.app.deploymentinterval", metric.WithDescription("a gauge of the interval between app deployments"))
	if err != nil {
		logger.Error(err, "unable to initialize app deployment interval OTel gauge")
	}

	appDeploymentDurationGauge, err := meter.Float64ObservableGauge("keptn.app.deploymentduration", metric.WithDescription("a gauge of the duration of app deployments"))
	if err != nil {
		logger.Error(err, "unable to initialize app deployment duration OTel gauge")
	}

	workloadDeploymentIntervalGauge, err := meter.Float64ObservableGauge("keptn.deployment.deploymentinterval", metric.WithDescription("a gauge of the interval between workload deployments"))
	if err != nil {
		logger.Error(err, "unable to initialize workload deployment interval OTel gauge")
	}

	workloadDeploymentDurationGauge, err := meter.Float64ObservableGauge("keptn.deployment.deploymentduration", metric.WithDescription("a gauge of the duration of workload deployments"))
	if err != nil {
		logger.Error(err, "unable to initialize workload deployment duration OTel gauge")
	}

	_, err = meter.RegisterCallback(
		func(ctx context.Context, o metric.Observer) error {
			observeActiveInstances(ctx, mgr, deploymentActiveGauge, appActiveGauge, taskActiveGauge, evaluationActiveGauge, o)
			observeDeploymentInterval(ctx, mgr, appDeploymentIntervalGauge, workloadDeploymentIntervalGauge, o)
			observeDuration(ctx, mgr, appDeploymentDurationGauge, workloadDeploymentDurationGauge, o)
			return nil
		},
		deploymentActiveGauge,
		taskActiveGauge,
		appActiveGauge,
		evaluationActiveGauge,
		appDeploymentIntervalGauge,
		appDeploymentDurationGauge,
		workloadDeploymentIntervalGauge,
		workloadDeploymentDurationGauge,
	)
	if err != nil {
		fmt.Println("Failed to register callback")
		panic(err)
	}
}

func observeDuration(ctx context.Context, mgr client.Client, appDeploymentDurationGauge metric.Float64ObservableGauge, workloadDeploymentDurationGauge metric.Float64ObservableGauge, observer metric.Observer) {

	err := ObserveDeploymentDuration(ctx, mgr, &lifecyclev1alpha3.KeptnAppVersionList{}, appDeploymentDurationGauge, observer)
	if err != nil {
		logger.Error(err, "unable to gather app deployment durations")
	}

	err = ObserveDeploymentDuration(ctx, mgr, &lifecyclev1alpha3.KeptnWorkloadVersionList{}, workloadDeploymentDurationGauge, observer)
	if err != nil {
		logger.Error(err, "unable to gather workload deployment durations")
	}

}

func observeDeploymentInterval(ctx context.Context, mgr client.Client, appDeploymentIntervalGauge metric.Float64ObservableGauge, workloadDeploymentIntervalGauge metric.Float64ObservableGauge, observer metric.Observer) {
	err := ObserveDeploymentInterval(ctx, mgr, &lifecyclev1alpha3.KeptnAppVersionList{}, appDeploymentIntervalGauge, observer)
	if err != nil {
		logger.Error(err, "unable to gather app deployment intervals")
	}

	err = ObserveDeploymentInterval(ctx, mgr, &lifecyclev1alpha3.KeptnWorkloadVersionList{}, workloadDeploymentIntervalGauge, observer)
	if err != nil {
		logger.Error(err, "unable to gather workload deployment intervals")
	}
}

func observeActiveInstances(ctx context.Context, mgr client.Client, deploymentActiveGauge metric.Int64ObservableGauge, appActiveGauge metric.Int64ObservableGauge, taskActiveGauge metric.Int64ObservableGauge, evaluationActiveGauge metric.Int64ObservableGauge, observer metric.Observer) {

	err := ObserveActiveInstances(ctx, mgr, &lifecyclev1alpha3.KeptnWorkloadVersionList{}, deploymentActiveGauge, observer)
	if err != nil {
		logger.Error(err, "unable to gather active deployments")
	}
	err = ObserveActiveInstances(ctx, mgr, &lifecyclev1alpha3.KeptnAppVersionList{}, appActiveGauge, observer)
	if err != nil {
		logger.Error(err, "unable to gather active apps")
	}
	err = ObserveActiveInstances(ctx, mgr, &lifecyclev1alpha3.KeptnTaskList{}, taskActiveGauge, observer)
	if err != nil {
		logger.Error(err, "unable to gather active tasks")
	}
	err = ObserveActiveInstances(ctx, mgr, &lifecyclev1alpha3.KeptnEvaluationList{}, evaluationActiveGauge, observer)
	if err != nil {
		logger.Error(err, "unable to gather active evaluations")
	}
}

func SetUpKeptnTaskMeters(meter interfaces.IMeter) common.KeptnMeters {
	deploymentCount, err := meter.Int64Counter("keptn.deployment.count", metric.WithDescription("a simple counter for Keptn Deployments"))
	if err != nil {
		logger.Error(err, "unable to initialize deployment count OTel counter")
	}
	deploymentDuration, err := meter.Float64Histogram("keptn.deployment.duration", metric.WithDescription("a histogram of duration for Keptn Deployments"), metric.WithUnit("s"))
	if err != nil {
		logger.Error(err, "unable to initialize deployment duration OTel histogram")
	}
	taskCount, err := meter.Int64Counter("keptn.task.count", metric.WithDescription("a simple counter for Keptn Tasks"))
	if err != nil {
		logger.Error(err, "unable to initialize task OTel counter")
	}
	taskDuration, err := meter.Float64Histogram("keptn.task.duration", metric.WithDescription("a histogram of duration for Keptn Tasks"), metric.WithUnit("s"))
	if err != nil {
		logger.Error(err, "unable to initialize task duration OTel histogram")
	}
	appCount, err := meter.Int64Counter("keptn.app.count", metric.WithDescription("a simple counter for Keptn Apps"))
	if err != nil {
		logger.Error(err, "unable to initialize app OTel counter")
	}
	appDuration, err := meter.Float64Histogram("keptn.app.duration", metric.WithDescription("a histogram of duration for Keptn Apps"), metric.WithUnit("s"))
	if err != nil {
		logger.Error(err, "unable to initialize app duration OTel histogram")
	}
	evaluationCount, err := meter.Int64Counter("keptn.evaluation.count", metric.WithDescription("a simple counter for Keptn Evaluations"))
	if err != nil {
		logger.Error(err, "unable to initialize evaluation OTel counter")
	}
	evaluationDuration, err := meter.Float64Histogram("keptn.evaluation.duration", metric.WithDescription("a histogram of duration for Keptn Evaluations"), metric.WithUnit("s"))
	if err != nil {
		logger.Error(err, "unable to initialize evaluation duration OTel histogram")
	}

	meters := common.KeptnMeters{
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
