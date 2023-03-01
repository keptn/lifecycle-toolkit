package v1alpha2

import (
	"fmt"

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2/common"
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	v1alpha3common "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"go.opentelemetry.io/otel/propagation"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the src v1alpha3.KeptnAppVersion to the hub version (v1alpha3.KeptnAppVersion)
//
//nolint:gocyclo
func (src *KeptnAppVersion) ConvertTo(dstRaw conversion.Hub) error {
	dst, ok := dstRaw.(*v1alpha3.KeptnAppVersion)

	if !ok {
		return fmt.Errorf("type %T %w", dstRaw, common.ErrCannotCastKeptnAppVersion)
	}

	// Copy equal stuff to new object
	// DO NOT COPY TypeMeta
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.Version = src.Spec.Version
	for _, srcWl := range src.Spec.Workloads {
		dst.Spec.Workloads = append(dst.Spec.Workloads, v1alpha3.KeptnWorkloadRef{
			Name:    srcWl.Name,
			Version: srcWl.Version,
		})
	}
	dst.Spec.PreDeploymentTasks = src.Spec.PreDeploymentTasks
	dst.Spec.PostDeploymentTasks = src.Spec.PostDeploymentTasks
	dst.Spec.PreDeploymentEvaluations = src.Spec.PreDeploymentEvaluations
	dst.Spec.PostDeploymentEvaluations = src.Spec.PostDeploymentEvaluations

	dst.Spec.AppName = src.Spec.AppName
	dst.Spec.PreviousVersion = src.Spec.PreviousVersion

	dst.Spec.TraceId = make(map[string]string, len(src.Spec.TraceId))
	for k, v := range src.Spec.TraceId {
		dst.Spec.TraceId[k] = v
	}

	dst.Status.PreDeploymentStatus = v1alpha3common.KeptnState(src.Status.PreDeploymentStatus)
	dst.Status.PostDeploymentStatus = v1alpha3common.KeptnState(src.Status.PostDeploymentStatus)
	dst.Status.PreDeploymentEvaluationStatus = v1alpha3common.KeptnState(src.Status.PreDeploymentEvaluationStatus)
	dst.Status.PostDeploymentEvaluationStatus = v1alpha3common.KeptnState(src.Status.PostDeploymentEvaluationStatus)
	dst.Status.WorkloadOverallStatus = v1alpha3common.KeptnState(src.Status.WorkloadOverallStatus)
	dst.Status.Status = v1alpha3common.KeptnState(src.Status.Status)

	for _, srcWls := range src.Status.WorkloadStatus {
		dst.Status.WorkloadStatus = append(dst.Status.WorkloadStatus, v1alpha3.WorkloadStatus{
			Workload: v1alpha3.KeptnWorkloadRef{
				Name:    srcWls.Workload.Name,
				Version: srcWls.Workload.Version,
			},
			Status: v1alpha3common.KeptnState(srcWls.Status),
		})
	}

	dst.Status.CurrentPhase = src.Status.CurrentPhase

	// Set sensible defaults for new fields
	dst.Spec.Revision = 1

	// Convert changed fields
	for _, item := range src.Status.PreDeploymentTaskStatus {
		dst.Status.PreDeploymentTaskStatus = append(dst.Status.PreDeploymentTaskStatus, v1alpha3.ItemStatus{
			DefinitionName: item.DefinitionName,
			Status:         v1alpha3common.KeptnState(item.Status),
			Name:           item.Name,
			StartTime:      item.StartTime,
			EndTime:        item.EndTime,
		})
	}

	for _, item := range src.Status.PostDeploymentTaskStatus {
		dst.Status.PostDeploymentTaskStatus = append(dst.Status.PostDeploymentTaskStatus, v1alpha3.ItemStatus{
			DefinitionName: item.DefinitionName,
			Status:         v1alpha3common.KeptnState(item.Status),
			Name:           item.Name,
			StartTime:      item.StartTime,
			EndTime:        item.EndTime,
		})
	}

	for _, item := range src.Status.PreDeploymentEvaluationTaskStatus {
		dst.Status.PreDeploymentEvaluationTaskStatus = append(dst.Status.PreDeploymentEvaluationTaskStatus, v1alpha3.ItemStatus{
			DefinitionName: item.DefinitionName,
			Status:         v1alpha3common.KeptnState(item.Status),
			Name:           item.Name,
			StartTime:      item.StartTime,
			EndTime:        item.EndTime,
		})
	}

	for _, item := range src.Status.PostDeploymentEvaluationTaskStatus {
		dst.Status.PostDeploymentEvaluationTaskStatus = append(dst.Status.PostDeploymentEvaluationTaskStatus, v1alpha3.ItemStatus{
			DefinitionName: item.DefinitionName,
			Status:         v1alpha3common.KeptnState(item.Status),
			Name:           item.Name,
			StartTime:      item.StartTime,
			EndTime:        item.EndTime,
		})
	}

	dst.Status.PhaseTraceIDs = make(v1alpha3common.PhaseTraceID, len(src.Status.PhaseTraceIDs))
	for k, v := range src.Status.PhaseTraceIDs {
		c := make(propagation.MapCarrier, len(v))
		for k1, v1 := range v {
			c[k1] = v1
		}
		dst.Status.PhaseTraceIDs[k] = c
	}

	dst.Status.StartTime = src.Status.StartTime
	dst.Status.EndTime = src.Status.EndTime

	return nil
}

// ConvertFrom converts from the hub version (v1alpha3.KeptnAppVersion) to this version (v1alpha3.KeptnAppVersion)
//
//nolint:gocyclo
func (dst *KeptnAppVersion) ConvertFrom(srcRaw conversion.Hub) error {
	src, ok := srcRaw.(*v1alpha3.KeptnAppVersion)

	if !ok {
		return fmt.Errorf("type %T %w", srcRaw, common.ErrCannotCastKeptnAppVersion)
	}

	// Copy equal stuff to new object
	// DO NOT COPY TypeMeta
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.Version = src.Spec.Version
	for _, srcWl := range src.Spec.Workloads {
		dst.Spec.Workloads = append(dst.Spec.Workloads, KeptnWorkloadRef{
			Name:    srcWl.Name,
			Version: srcWl.Version,
		})
	}
	dst.Spec.PreDeploymentTasks = src.Spec.PreDeploymentTasks
	dst.Spec.PostDeploymentTasks = src.Spec.PostDeploymentTasks
	dst.Spec.PreDeploymentEvaluations = src.Spec.PreDeploymentEvaluations
	dst.Spec.PostDeploymentEvaluations = src.Spec.PostDeploymentEvaluations

	dst.Spec.AppName = src.Spec.AppName
	dst.Spec.PreviousVersion = src.Spec.PreviousVersion

	dst.Spec.TraceId = make(map[string]string, len(src.Spec.TraceId))
	for k, v := range src.Spec.TraceId {
		dst.Spec.TraceId[k] = v
	}

	dst.Status.PreDeploymentStatus = common.KeptnState(src.Status.PreDeploymentStatus)
	dst.Status.PostDeploymentStatus = common.KeptnState(src.Status.PostDeploymentStatus)
	dst.Status.PreDeploymentEvaluationStatus = common.KeptnState(src.Status.PreDeploymentEvaluationStatus)
	dst.Status.PostDeploymentEvaluationStatus = common.KeptnState(src.Status.PostDeploymentEvaluationStatus)
	dst.Status.WorkloadOverallStatus = common.KeptnState(src.Status.WorkloadOverallStatus)
	dst.Status.Status = common.KeptnState(src.Status.Status)

	for _, srcWls := range src.Status.WorkloadStatus {
		dst.Status.WorkloadStatus = append(dst.Status.WorkloadStatus, WorkloadStatus{
			Workload: KeptnWorkloadRef{
				Name:    srcWls.Workload.Name,
				Version: srcWls.Workload.Version,
			},
			Status: common.KeptnState(srcWls.Status),
		})
	}

	dst.Status.CurrentPhase = src.Status.CurrentPhase

	// Convert changed fields
	for _, item := range src.Status.PreDeploymentTaskStatus {
		dst.Status.PreDeploymentTaskStatus = append(dst.Status.PreDeploymentTaskStatus, ItemStatus{
			DefinitionName: item.DefinitionName,
			Status:         common.KeptnState(item.Status),
			Name:           item.Name,
			StartTime:      item.StartTime,
			EndTime:        item.EndTime,
		})
	}

	for _, item := range src.Status.PostDeploymentTaskStatus {
		dst.Status.PostDeploymentTaskStatus = append(dst.Status.PostDeploymentTaskStatus, ItemStatus{
			DefinitionName: item.DefinitionName,
			Status:         common.KeptnState(item.Status),
			Name:           item.Name,
			StartTime:      item.StartTime,
			EndTime:        item.EndTime,
		})
	}

	for _, item := range src.Status.PreDeploymentEvaluationTaskStatus {
		dst.Status.PreDeploymentEvaluationTaskStatus = append(dst.Status.PreDeploymentEvaluationTaskStatus, ItemStatus{
			DefinitionName: item.DefinitionName,
			Status:         common.KeptnState(item.Status),
			Name:           item.Name,
			StartTime:      item.StartTime,
			EndTime:        item.EndTime,
		})
	}

	for _, item := range src.Status.PostDeploymentEvaluationTaskStatus {
		dst.Status.PostDeploymentEvaluationTaskStatus = append(dst.Status.PostDeploymentEvaluationTaskStatus, ItemStatus{
			DefinitionName: item.DefinitionName,
			Status:         common.KeptnState(item.Status),
			Name:           item.Name,
			StartTime:      item.StartTime,
			EndTime:        item.EndTime,
		})
	}

	dst.Status.PhaseTraceIDs = make(common.PhaseTraceID, len(src.Status.PhaseTraceIDs))
	for k, v := range src.Status.PhaseTraceIDs {
		c := make(propagation.MapCarrier, len(v))
		for k1, v1 := range v {
			c[k1] = v1
		}
		dst.Status.PhaseTraceIDs[k] = c
	}

	dst.Status.StartTime = src.Status.StartTime
	dst.Status.EndTime = src.Status.EndTime

	return nil
}
