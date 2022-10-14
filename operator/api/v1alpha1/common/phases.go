package common

type KeptnPhase KeptnPhaseType

type KeptnPhaseType struct {
	LongName  string
	ShortName string
}

var (
	PhaseWorkloadPreDeployment  = KeptnPhaseType{LongName: "Workload Pre-Deployment", ShortName: "WorkloadPreDeploy"}
	PhaseWorkloadPostDeployment = KeptnPhaseType{LongName: "Workload Post-Deployment", ShortName: "WorkloadPostDeploy"}
	PhaseWorkloadPreAnalysis    = KeptnPhaseType{LongName: "Workload Pre-Analysis", ShortName: "WorkloadPreAnalysis"}
	PhaseWorkloadPostAnalysis   = KeptnPhaseType{LongName: "Workload Post-Analysis", ShortName: "WorkloadPostAnalysis"}
	PhaseWorkloadDeployment     = KeptnPhaseType{LongName: "Workload Deployment", ShortName: "WorkloadDeploy"}
	PhaseAppPreDeployment       = KeptnPhaseType{LongName: "App Pre-Deployment", ShortName: "AppPreDeploy"}
	PhaseAppPostDeployment      = KeptnPhaseType{LongName: "App Post-Deployment", ShortName: "AppPostDeploy"}
	PhaseAppPreAnalysis         = KeptnPhaseType{LongName: "App Pre-Analysis", ShortName: "AppPreAnalysis"}
	PhaseAppPostAnalysis        = KeptnPhaseType{LongName: "App Post-Analysis", ShortName: "AppPostAnalysis"}
	PhaseAppDeployment          = KeptnPhaseType{LongName: "App Deployment", ShortName: "AppDeploy"}
)
