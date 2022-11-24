package common

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/propagation"
)

func TestKeptnPhaseType_IsEvaluation(t *testing.T) {
	tests := []struct {
		State KeptnPhaseType
		Want  bool
	}{
		{
			State: PhaseWorkloadDeployment,
			Want:  false,
		},
		{
			State: PhaseWorkloadPostEvaluation,
			Want:  true,
		},
		{
			State: PhaseWorkloadPreEvaluation,
			Want:  true,
		},
		{
			State: PhaseAppPostEvaluation,
			Want:  true,
		},
		{
			State: PhaseAppPreEvaluation,
			Want:  true,
		},
		{
			State: PhaseAppPreDeployment,
			Want:  false,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tt.State.IsEvaluation(), tt.Want)
		})
	}
}

func TestKeptnPhaseType_IsPreEvaluation(t *testing.T) {
	tests := []struct {
		State KeptnPhaseType
		Want  bool
	}{
		{
			State: PhaseWorkloadDeployment,
			Want:  false,
		},
		{
			State: PhaseWorkloadPostEvaluation,
			Want:  false,
		},
		{
			State: PhaseWorkloadPreEvaluation,
			Want:  true,
		},
		{
			State: PhaseAppPostEvaluation,
			Want:  false,
		},
		{
			State: PhaseAppPreEvaluation,
			Want:  true,
		},
		{
			State: PhaseAppPreDeployment,
			Want:  false,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tt.State.IsPreEvaluation(), tt.Want)
		})
	}
}

func TestKeptnPhaseType_IsPostEvaluation(t *testing.T) {
	tests := []struct {
		State KeptnPhaseType
		Want  bool
	}{
		{
			State: PhaseWorkloadDeployment,
			Want:  false,
		},
		{
			State: PhaseWorkloadPostEvaluation,
			Want:  true,
		},
		{
			State: PhaseWorkloadPreEvaluation,
			Want:  false,
		},
		{
			State: PhaseAppPostEvaluation,
			Want:  true,
		},
		{
			State: PhaseAppPreEvaluation,
			Want:  false,
		},
		{
			State: PhaseAppPreDeployment,
			Want:  false,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tt.State.IsPostEvaluation(), tt.Want)
		})
	}
}

func TestKeptnPhaseType_IsTask(t *testing.T) {
	tests := []struct {
		State KeptnPhaseType
		Want  bool
	}{
		{
			State: PhaseWorkloadDeployment,
			Want:  false,
		},
		{
			State: PhaseWorkloadPostDeployment,
			Want:  true,
		},
		{
			State: PhaseWorkloadPreDeployment,
			Want:  true,
		},
		{
			State: PhaseAppPostDeployment,
			Want:  true,
		},
		{
			State: PhaseAppPreDeployment,
			Want:  true,
		},
		{
			State: PhaseAppPreEvaluation,
			Want:  false,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tt.State.IsTask(), tt.Want)
		})
	}
}

func TestKeptnPhaseType_IsPreTask(t *testing.T) {
	tests := []struct {
		State KeptnPhaseType
		Want  bool
	}{
		{
			State: PhaseWorkloadDeployment,
			Want:  false,
		},
		{
			State: PhaseWorkloadPostDeployment,
			Want:  false,
		},
		{
			State: PhaseWorkloadPreDeployment,
			Want:  true,
		},
		{
			State: PhaseAppPostDeployment,
			Want:  false,
		},
		{
			State: PhaseAppPreDeployment,
			Want:  true,
		},
		{
			State: PhaseAppPreEvaluation,
			Want:  false,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tt.State.IsPreTask(), tt.Want)
		})
	}
}

func TestKeptnPhaseType_IsPostTask(t *testing.T) {
	tests := []struct {
		State KeptnPhaseType
		Want  bool
	}{
		{
			State: PhaseWorkloadDeployment,
			Want:  false,
		},
		{
			State: PhaseWorkloadPostDeployment,
			Want:  true,
		},
		{
			State: PhaseWorkloadPreDeployment,
			Want:  false,
		},
		{
			State: PhaseAppPostDeployment,
			Want:  true,
		},
		{
			State: PhaseAppPreDeployment,
			Want:  false,
		},
		{
			State: PhaseAppPreEvaluation,
			Want:  false,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			require.Equal(t, tt.State.IsPostTask(), tt.Want)
		})
	}
}

func TestPhaseTraceID(t *testing.T) {
	trace := PhaseTraceID{}

	trace.SetPhaseTraceID(PhaseAppDeployment.LongName, propagation.MapCarrier{
		"name":  "trace",
		"name2": "trace2",
	})

	require.Equal(t, PhaseTraceID{
		PhaseAppDeployment.ShortName: propagation.MapCarrier{
			"name":  "trace",
			"name2": "trace2",
		},
	}, trace)

	trace.SetPhaseTraceID(PhaseWorkloadDeployment.ShortName, propagation.MapCarrier{
		"name3": "trace3",
	})

	require.Equal(t, PhaseTraceID{
		PhaseAppDeployment.ShortName: propagation.MapCarrier{
			"name":  "trace",
			"name2": "trace2",
		},
		PhaseWorkloadDeployment.ShortName: propagation.MapCarrier{
			"name3": "trace3",
		},
	}, trace)

	require.Equal(t, propagation.MapCarrier{
		"name":  "trace",
		"name2": "trace2",
	}, trace.GetPhaseTraceID(PhaseAppDeployment.LongName))

	require.Equal(t, propagation.MapCarrier{
		"name3": "trace3",
	}, trace.GetPhaseTraceID(PhaseWorkloadDeployment.ShortName))
}

func TestGetShortPhaseName(t *testing.T) {
	require.Equal(t, "WorkloadPreDeployTasks", GetShortPhaseName("WorkloadPreDeployTasks"))
	require.Equal(t, "WorkloadPreDeployTasks", GetShortPhaseName("Workload Pre-Deployment Tasks"))
	require.Equal(t, "", GetShortPhaseName("Workload Pre-Deploycdddment Tasks"))
}
