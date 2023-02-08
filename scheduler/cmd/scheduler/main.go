/*
Copyright 2020 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/keptn/lifecycle-toolkit/scheduler/pkg/klcpermit"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpgrpc"
	"go.opentelemetry.io/otel/exporters/stdout"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/semconv"
	"google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/component-base/cli"
	"k8s.io/klog/v2"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"
)

var (
	gitCommit    string
	buildTime    string
	buildVersion string
)

type envConfig struct {
	OTelCollectorURL string `envconfig:"OTEL_COLLECTOR_URL" default:""`
}

type keptnSchedulerOTelErrorHandler struct{}

func (keptnSchedulerOTelErrorHandler) Handle(_ error) {
	// ignoring any OTel errors
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Fatalf("Failed to process env var: %s", err)
	}

	tp := initOTel(env)

	rand.Seed(time.Now().UnixNano())
	command := app.NewSchedulerCommand(
		app.WithPlugin(klcpermit.PluginName, klcpermit.New),
	)

	code := cli.Run(command)

	err := tp.Shutdown(context.TODO())
	if err != nil {
		fmt.Println(err)
	}
	os.Exit(code)

}

func initOTel(env envConfig) *sdktrace.TracerProvider {
	tpOptions, err := getOTelTracerProviderOptions(env)
	if err != nil {
		log.Panicf("failed to initialize OTel options")
	}
	tp := sdktrace.NewTracerProvider(tpOptions...)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	otel.SetErrorHandler(keptnSchedulerOTelErrorHandler{})
	return tp
}

func getOTelTracerProviderOptions(env envConfig) ([]sdktrace.TracerProviderOption, error) {
	tracerProviderOptions := []sdktrace.TracerProviderOption{}

	stdOutExp, err := newStdOutExporter()
	if err != nil {
		return nil, fmt.Errorf("could not create stdout OTel exporter: %w", err)
	}
	tracerProviderOptions = append(tracerProviderOptions, sdktrace.WithBatcher(stdOutExp))

	if env.OTelCollectorURL != "" {
		// try to set OTel exporter for Jaeger
		otelExporter, err := newOTelExporter(env)
		if err != nil {
			// log the error, but do not break if Jaeger exporter cannot be created
			klog.Errorf("Could not set up OTel exporter: %v", err)
		} else if otelExporter != nil {
			tracerProviderOptions = append(tracerProviderOptions, sdktrace.WithBatcher(otelExporter))
		}
	}
	tracerProviderOptions = append(tracerProviderOptions, sdktrace.WithSampler(sdktrace.AlwaysSample()))
	tracerProviderOptions = append(tracerProviderOptions, sdktrace.WithResource(newResource()))

	return tracerProviderOptions, nil
}

func newStdOutExporter() (sdktrace.SpanExporter, error) {
	return stdout.NewExporter(stdout.WithPrettyPrint())
}

func newOTelExporter(env envConfig) (sdktrace.SpanExporter, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
	defer cancel()

	_, err := net.DialTimeout("tcp", env.OTelCollectorURL, 2*time.Second)
	if err != nil {
		return nil, err
	}

	driver := otlpgrpc.NewDriver(
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint(env.OTelCollectorURL),
		otlpgrpc.WithDialOption(grpc.WithBlock()), // useful for testing
	)
	traceExporter, err := otlp.NewExporter(ctx, driver)
	if err != nil {
		return nil, err
	}
	return traceExporter, nil
}

func newResource() *resource.Resource {
	return resource.NewWithAttributes(
		semconv.TelemetrySDKLanguageGo,
		semconv.ServiceNameKey.String("keptn-lifecycle-scheduler"),
		semconv.ServiceVersionKey.String(buildVersion+"-"+gitCommit+"-"+buildTime),
	)
}
