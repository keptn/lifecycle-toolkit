package common

type KeptnPhase KeptnPhaseType

type KeptnPhaseType struct {
	LongName  string
	ShortName string
}

var (
	PhaseWorkloadPreDeployment  = KeptnPhaseType{LongName: "Workload Pre-Deployment", ShortName: "WorkloadPreDeploy"}
	PhaseWorkloadPostDeployment = KeptnPhaseType{LongName: "Workload Post-Deployment", ShortName: "WorkloadPostDeploy"}
	PhaseWorkloadDeployment     = KeptnPhaseType{LongName: "Workload Deployment", ShortName: "WorkloadDeploy"}
	PhaseAppPreDeployment       = KeptnPhaseType{LongName: "App Pre-Deployment", ShortName: "AppPreDeploy"}
	PhaseAppPostDeployment      = KeptnPhaseType{LongName: "App Post-Deployment", ShortName: "AppPostDeploy"}
	PhaseAppDeployment          = KeptnPhaseType{LongName: "App Deployment", ShortName: "AppDeploy"}
)
