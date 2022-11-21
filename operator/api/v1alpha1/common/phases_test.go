package common

import (
	"testing"

	"github.com/stretchr/testify/require"
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
