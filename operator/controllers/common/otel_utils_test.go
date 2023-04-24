package common

import (
	"context"
	"net"
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/interfaces/fake"
	fakeasync "github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/interfaces/fake/async"
	fakesync "github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/interfaces/fake/sync"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	"google.golang.org/grpc"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func TestGetOTelTracerProviderOptions(t *testing.T) {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	go func() {
		err := s.Serve(listener)
		if err != nil {
			panic(err)
		}
	}()

	defer s.Stop()

	type args struct {
		oTelCollectorUrl string
	}
	tests := []struct {
		name            string
		args            args
		wantArrayLength int
		wantErr         bool
	}{
		{
			name: "Test with no URL",
			args: args{
				oTelCollectorUrl: "",
			},
			wantArrayLength: 2,
		},
		{
			name: "Test with wrong URL",
			args: args{
				oTelCollectorUrl: "error-url",
			},
			wantArrayLength: 0,
			wantErr:         true,
		},
		{
			name: "Test with URL",
			args: args{
				oTelCollectorUrl: "localhost:9000",
			},
			wantArrayLength: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//TODO also test underline return
			got, _, err := GetOTelTracerProviderOptions(tt.args.oTelCollectorUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOTelTracerProviderOptions() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantArrayLength {
				t.Errorf("GetOTelTracerProviderOptions() got length = %v, wantArrayLength %v", got, tt.wantArrayLength)

			}
		})
	}
}

func TestSetUpKeptnMeters(t *testing.T) {
	fakeAsyncIntTracerProvider := &fakeasync.ITracerProviderAsyncInt64Mock{
		GaugeFunc: func(name string, opts ...instrument.Option) (asyncint64.Gauge, error) {
			return nil, errors.New("some error")
		},
	}
	fakeAsyncFloatTracerProvider := &fakeasync.ITracerProviderAsyncFloat64Mock{
		GaugeFunc: func(name string, opts ...instrument.Option) (asyncfloat64.Gauge, error) {
			return nil, errors.New("some error")
		},
	}

	fakeMeter := &fake.IMeterMock{
		AsyncInt64Func: func() asyncint64.InstrumentProvider {
			return fakeAsyncIntTracerProvider
		},
		AsyncFloat64Func: func() asyncfloat64.InstrumentProvider {
			return fakeAsyncFloatTracerProvider
		},
		RegisterCallbackFunc: func(insts []instrument.Asynchronous, function func(context.Context)) error {
			return nil
		},
	}

	type args struct {
		meter metric.Meter
		mgr   client.Client
	}
	tests := []struct {
		name              string
		args              args
		wantRegisterCalls int
	}{
		{
			name: "Basic case",
			args: args{
				meter: metric.NewNoopMeter(),
				mgr:   nil,
			},
			wantRegisterCalls: 0,
		},
		{
			name: "Error case",
			args: args{
				meter: fakeMeter,
				mgr:   nil,
			},
			wantRegisterCalls: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetUpKeptnMeters(tt.args.meter, tt.args.mgr)
			require.Equal(t, tt.wantRegisterCalls, len(fakeMeter.RegisterCallbackCalls()))
		})
	}
}

func TestSetUpKeptnMetersError(t *testing.T) {
	fakeAsyncIntTracerProvider := &fakeasync.ITracerProviderAsyncInt64Mock{
		GaugeFunc: func(name string, opts ...instrument.Option) (asyncint64.Gauge, error) {
			return nil, errors.New("some error")
		},
	}
	fakeAsyncFloatTracerProvider := &fakeasync.ITracerProviderAsyncFloat64Mock{
		GaugeFunc: func(name string, opts ...instrument.Option) (asyncfloat64.Gauge, error) {
			return nil, errors.New("some error")
		},
	}

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	errorFakeMeter := &fake.IMeterMock{
		AsyncInt64Func: func() asyncint64.InstrumentProvider {
			return fakeAsyncIntTracerProvider
		},
		AsyncFloat64Func: func() asyncfloat64.InstrumentProvider {
			return fakeAsyncFloatTracerProvider
		},
		RegisterCallbackFunc: func(insts []instrument.Asynchronous, function func(context.Context)) error {
			return errors.New("some error")
		},
	}

	SetUpKeptnMeters(errorFakeMeter, nil)
}

func TestSetUpKeptnTaskMeters(t *testing.T) {
	noopMeter := metric.NewNoopMeter()

	got := SetUpKeptnTaskMeters(noopMeter)

	require.NotNil(t, got.TaskCount)
	require.NotNil(t, got.TaskDuration)
	require.NotNil(t, got.DeploymentCount)
	require.NotNil(t, got.DeploymentDuration)
	require.NotNil(t, got.AppCount)
	require.NotNil(t, got.AppDuration)
	require.NotNil(t, got.EvaluationCount)
	require.NotNil(t, got.EvaluationDuration)
}

func TestSetUpKeptnTaskMeters_ErrorCase(t *testing.T) {
	fakeSyncIntTracerProvider := &fakesync.ITracerProviderSyncInt64Mock{
		CounterFunc: func(name string, opts ...instrument.Option) (syncint64.Counter, error) {
			return nil, errors.New("some error")
		},
		HistogramFunc: func(name string, opts ...instrument.Option) (syncint64.Histogram, error) {
			return nil, errors.New("some error")
		},
	}
	fakeSyncFloatTracerProvider := &fakesync.ITracerProviderSyncFloat64Mock{
		CounterFunc: func(name string, opts ...instrument.Option) (syncfloat64.Counter, error) {
			return nil, errors.New("some error")
		},
		HistogramFunc: func(name string, opts ...instrument.Option) (syncfloat64.Histogram, error) {
			return nil, errors.New("some error")
		},
	}

	errorFakeMeter := &fake.IMeterMock{
		SyncInt64Func: func() syncint64.InstrumentProvider {
			return fakeSyncIntTracerProvider
		},
		SyncFloat64Func: func() syncfloat64.InstrumentProvider {
			return fakeSyncFloatTracerProvider
		},
		RegisterCallbackFunc: func(insts []instrument.Asynchronous, function func(context.Context)) error {
			return errors.New("some error")
		},
	}

	got := SetUpKeptnTaskMeters(errorFakeMeter)

	require.Nil(t, got.TaskCount)
	require.Nil(t, got.TaskDuration)
	require.Nil(t, got.DeploymentCount)
	require.Nil(t, got.DeploymentDuration)
	require.Nil(t, got.AppCount)
	require.Nil(t, got.AppDuration)
	require.Nil(t, got.EvaluationCount)
	require.Nil(t, got.EvaluationDuration)
}

func Test_otelConfig_GetTracer(t *testing.T) {
	otelConfig := GetOtelInstance()

	tracer := otelConfig.GetTracer("new-tracer")
	require.NotNil(t, tracer)

	otelConfig.cleanTracers()

	require.Empty(t, otelConfig.tracers)
}

func Test_otelConfig_InitOtelCollector_ReInitWithSameURL(t *testing.T) {

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	go func() {
		err := s.Serve(listener)
		if err != nil {
			panic(err)
		}
	}()

	defer s.Stop()

	o := GetOtelInstance()
	err = o.InitOtelCollector("localhost:9000")

	require.Nil(t, err)

	require.Equal(t, "localhost:9000", o.lastAppliedCollectorURL)

	tracer := o.GetTracer("my-tracer")
	require.NotNil(t, tracer)
	require.Len(t, o.tracers, 1)

	// init with the same URL again
	err = o.InitOtelCollector("localhost:9000")

	require.Nil(t, err)

	// in this case the init function should NOT have changed any internal state,
	// i.e. the tracers should NOT have been cleaned up
	require.Len(t, o.tracers, 1)
}

func Test_otelConfig_InitOtelCollector_ReInitWithDifferentURL(t *testing.T) {

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}

	listener2, err := net.Listen("tcp", ":9001")
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	go func() {
		err := s.Serve(listener)
		if err != nil {
			panic(err)
		}
	}()

	s2 := grpc.NewServer()
	go func() {
		err := s2.Serve(listener2)
		if err != nil {
			panic(err)
		}
	}()

	defer s.Stop()

	o := GetOtelInstance()
	err = o.InitOtelCollector("localhost:9000")

	require.Nil(t, err)

	require.Equal(t, "localhost:9000", o.lastAppliedCollectorURL)

	tracer := o.GetTracer("my-tracer")
	require.NotNil(t, tracer)
	require.Len(t, o.tracers, 1)

	// init with a different URL
	err = o.InitOtelCollector("localhost:9001")

	require.Nil(t, err)

	// in this case the init function should have changed any internal state,
	// i.e. the tracers should have been cleaned up
	require.Empty(t, o.tracers)
}
