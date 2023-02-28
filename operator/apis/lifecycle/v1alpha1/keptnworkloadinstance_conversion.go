package v1alpha1

import (
	"fmt"

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha1/common"
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	v1alpha3common "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"go.opentelemetry.io/otel/propagation"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the src v1alpha1.KeptnWorkloadInstance to the hub version (v1alpha3.KeptnWorkloadInstance)
func (src *KeptnWorkloadInstance) ConvertTo(dstRaw conversion.Hub) error {
	dst, ok := dstRaw.(*v1alpha3.KeptnWorkloadInstance)

	if !ok {
		return fmt.Errorf("type %T %w", dstRaw, common.ErrCannotCastKeptnWorkloadInstance)
	}

	// Copy equal stuff to new object
	// DO NOT COPY TypeMeta
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.AppName = src.Spec.AppName
	dst.Spec.Version = src.Spec.Version
	dst.Spec.PreDeploymentTasks = src.Spec.PreDeploymentTasks
	dst.Spec.PostDeploymentTasks = src.Spec.PostDeploymentTasks
	dst.Spec.PreDeploymentEvaluations = src.Spec.PreDeploymentEvaluations
	dst.Spec.PostDeploymentEvaluations = src.Spec.PostDeploymentEvaluations
	dst.Spec.ResourceReference = v1alpha3.ResourceReference{
		UID:  src.Spec.ResourceReference.UID,
		Kind: src.Spec.ResourceReference.Kind,
		Name: src.Spec.ResourceReference.Name,
	}

	dst.Spec.WorkloadName = src.Spec.WorkloadName
	dst.Spec.PreviousVersion = src.Spec.PreviousVersion
	dst.Spec.TraceId = make(map[string]string, len(src.Spec.TraceId))
	for k, v := range src.Spec.TraceId {
		dst.Spec.TraceId[k] = v
	}

	dst.Status.PreDeploymentStatus = v1alpha3common.KeptnState(src.Status.PreDeploymentStatus)
	dst.Status.PostDeploymentStatus = v1alpha3common.KeptnState(src.Status.PostDeploymentStatus)
	dst.Status.PreDeploymentEvaluationStatus = v1alpha3common.KeptnState(src.Status.PreDeploymentEvaluationStatus)
	dst.Status.PostDeploymentEvaluationStatus = v1alpha3common.KeptnState(src.Status.PostDeploymentEvaluationStatus)
	dst.Status.DeploymentStatus = v1alpha3common.KeptnState(src.Status.DeploymentStatus)
	dst.Status.Status = v1alpha3common.KeptnState(src.Status.Status)

	dst.Status.CurrentPhase = src.Status.CurrentPhase

	// Convert changed fields
	for _, item := range src.Status.PreDeploymentTaskStatus {
		dst.Status.PreDeploymentTaskStatus = append(dst.Status.PreDeploymentTaskStatus, v1alpha3.ItemStatus{
			DefinitionName: item.TaskDefinitionName,
			Status:         v1alpha3common.KeptnState(item.Status),
			Name:           item.TaskName,
			StartTime:      item.StartTime,
			EndTime:        item.EndTime,
		})
	}

	for _, item := range src.Status.PostDeploymentTaskStatus {
		dst.Status.PostDeploymentTaskStatus = append(dst.Status.PostDeploymentTaskStatus, v1alpha3.ItemStatus{
			DefinitionName: item.TaskDefinitionName,
			Status:         v1alpha3common.KeptnState(item.Status),
			Name:           item.TaskName,
			StartTime:      item.StartTime,
			EndTime:        item.EndTime,
		})
	}

	for _, item := range src.Status.PreDeploymentEvaluationTaskStatus {
		dst.Status.PreDeploymentEvaluationTaskStatus = append(dst.Status.PreDeploymentEvaluationTaskStatus, v1alpha3.ItemStatus{
			DefinitionName: item.EvaluationDefinitionName,
			Status:         v1alpha3common.KeptnState(item.Status),
			Name:           item.EvaluationName,
			StartTime:      item.StartTime,
			EndTime:        item.EndTime,
		})
	}

	for _, item := range src.Status.PostDeploymentEvaluationTaskStatus {
		dst.Status.PostDeploymentEvaluationTaskStatus = append(dst.Status.PostDeploymentEvaluationTaskStatus, v1alpha3.ItemStatus{
			DefinitionName: item.EvaluationDefinitionName,
			Status:         v1alpha3common.KeptnState(item.Status),
			Name:           item.EvaluationName,
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

// ConvertFrom converts from the hub version (v1alpha3.KeptnWorkloadInstance) to this version (v1alpha1.KeptnWorkloadInstance)
func (dst *KeptnWorkloadInstance) ConvertFrom(srcRaw conversion.Hub) error {
	src, ok := srcRaw.(*v1alpha3.KeptnWorkloadInstance)

	if !ok {
		return fmt.Errorf("type %T %w", srcRaw, common.ErrCannotCastKeptnWorkloadInstance)
	}

	// Copy equal stuff to new object
	// DO NOT COPY TypeMeta
	dst.ObjectMeta = src.ObjectMeta

	dst.Spec.AppName = src.Spec.AppName
	dst.Spec.Version = src.Spec.Version
	dst.Spec.PreDeploymentTasks = src.Spec.PreDeploymentTasks
	dst.Spec.PostDeploymentTasks = src.Spec.PostDeploymentTasks
	dst.Spec.PreDeploymentEvaluations = src.Spec.PreDeploymentEvaluations
	dst.Spec.PostDeploymentEvaluations = src.Spec.PostDeploymentEvaluations
	dst.Spec.ResourceReference = ResourceReference{
		UID:  src.Spec.ResourceReference.UID,
		Kind: src.Spec.ResourceReference.Kind,
		Name: src.Spec.ResourceReference.Name,
	}

	dst.Spec.WorkloadName = src.Spec.WorkloadName
	dst.Spec.PreviousVersion = src.Spec.PreviousVersion
	dst.Spec.TraceId = make(map[string]string, len(src.Spec.TraceId))
	for k, v := range src.Spec.TraceId {
		dst.Spec.TraceId[k] = v
	}

	dst.Status.PreDeploymentStatus = common.KeptnState(src.Status.PreDeploymentStatus)
	dst.Status.PostDeploymentStatus = common.KeptnState(src.Status.PostDeploymentStatus)
	dst.Status.PreDeploymentEvaluationStatus = common.KeptnState(src.Status.PreDeploymentEvaluationStatus)
	dst.Status.PostDeploymentEvaluationStatus = common.KeptnState(src.Status.PostDeploymentEvaluationStatus)
	dst.Status.DeploymentStatus = common.KeptnState(src.Status.DeploymentStatus)
	dst.Status.Status = common.KeptnState(src.Status.Status)

	dst.Status.CurrentPhase = src.Status.CurrentPhase

	// Convert changed fields
	for _, item := range src.Status.PreDeploymentTaskStatus {
		dst.Status.PreDeploymentTaskStatus = append(dst.Status.PreDeploymentTaskStatus, TaskStatus{
			TaskDefinitionName: item.DefinitionName,
			Status:             common.KeptnState(item.Status),
			TaskName:           item.Name,
			StartTime:          item.StartTime,
			EndTime:            item.EndTime,
		})
	}

	for _, item := range src.Status.PostDeploymentTaskStatus {
		dst.Status.PostDeploymentTaskStatus = append(dst.Status.PostDeploymentTaskStatus, TaskStatus{
			TaskDefinitionName: item.DefinitionName,
			Status:             common.KeptnState(item.Status),
			TaskName:           item.Name,
			StartTime:          item.StartTime,
			EndTime:            item.EndTime,
		})
	}

	for _, item := range src.Status.PreDeploymentEvaluationTaskStatus {
		dst.Status.PreDeploymentEvaluationTaskStatus = append(dst.Status.PreDeploymentEvaluationTaskStatus, EvaluationStatus{
			EvaluationDefinitionName: item.DefinitionName,
			Status:                   common.KeptnState(item.Status),
			EvaluationName:           item.Name,
			StartTime:                item.StartTime,
			EndTime:                  item.EndTime,
		})
	}

	for _, item := range src.Status.PostDeploymentEvaluationTaskStatus {
		dst.Status.PostDeploymentEvaluationTaskStatus = append(dst.Status.PostDeploymentEvaluationTaskStatus, EvaluationStatus{
			EvaluationDefinitionName: item.DefinitionName,
			Status:                   common.KeptnState(item.Status),
			EvaluationName:           item.Name,
			StartTime:                item.StartTime,
			EndTime:                  item.EndTime,
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
