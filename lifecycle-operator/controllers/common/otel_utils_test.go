package common

import (
	"net"
	"testing"

	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/interfaces"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/interfaces/fake"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
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
			wantArrayLength: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// TODO also test underline return
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
	fakeMeter := &fake.IMeterMock{
		Int64ObservableGaugeFunc: func(name string, options ...metric.Int64ObservableGaugeOption) (metric.Int64ObservableGauge, error) {
			return nil, errors.New("some error")
		},
		Float64ObservableGaugeFunc: func(name string, options ...metric.Float64ObservableGaugeOption) (metric.Float64ObservableGauge, error) {
			return nil, errors.New("some error")
		},
		RegisterCallbackFunc: func(f metric.Callback, instruments ...metric.Observable) (metric.Registration, error) {
			return nil, nil
		},
	}

	type args struct {
		meter interfaces.IMeter
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
				meter: noop.NewMeterProvider().Meter(("test")),
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
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	errorFakeMeter := &fake.IMeterMock{
		Int64ObservableGaugeFunc: func(name string, options ...metric.Int64ObservableGaugeOption) (metric.Int64ObservableGauge, error) {
			return nil, errors.New("some error")
		},
		Float64ObservableGaugeFunc: func(name string, options ...metric.Float64ObservableGaugeOption) (metric.Float64ObservableGauge, error) {
			return nil, errors.New("some error")
		},
		RegisterCallbackFunc: func(f metric.Callback, instruments ...metric.Observable) (metric.Registration, error) {
			return nil, errors.New("some error")
		},
	}

	SetUpKeptnMeters(errorFakeMeter, nil)
}

func TestSetUpKeptnTaskMeters(t *testing.T) {
	got := SetUpKeptnTaskMeters(noop.NewMeterProvider().Meter(("test")))

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
	errorFakeMeter := &fake.IMeterMock{
		Int64CounterFunc: func(name string, options ...metric.Int64CounterOption) (metric.Int64Counter, error) {
			return nil, errors.New("some error")
		},
		Int64HistogramFunc: func(name string, options ...metric.Int64HistogramOption) (metric.Int64Histogram, error) {
			return nil, errors.New("some error")
		},
		Float64CounterFunc: func(name string, options ...metric.Float64CounterOption) (metric.Float64Counter, error) {
			return nil, errors.New("some error")
		},
		Float64HistogramFunc: func(name string, options ...metric.Float64HistogramOption) (metric.Float64Histogram, error) {
			return nil, errors.New("some error")
		},
		RegisterCallbackFunc: func(f metric.Callback, instruments ...metric.Observable) (metric.Registration, error) {
			return nil, errors.New("some error")
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
