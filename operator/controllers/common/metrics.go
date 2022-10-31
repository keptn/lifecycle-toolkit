package common

import (
	"context"
	"fmt"

	apicommon "github.com/keptn/lifecycle-controller/operator/api/v1alpha1/common"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func GetDeploymentDuration(ctx context.Context, client client.Client, reconcileObjectList client.ObjectList) ([]apicommon.GaugeFloatValue, error) {
	err := client.List(ctx, reconcileObjectList)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve workload instances: %w", err)
	}

	piWrapper, err := NewListItemWrapperFromClientObjectList(reconcileObjectList)
	if err != nil {
		return nil, err
	}

	res := []apicommon.GaugeFloatValue{}

	for _, reconcileObject := range piWrapper.GetItems() {
		if reconcileObject.IsEndTimeSet() {
			duration := reconcileObject.GetEndTime().Sub(reconcileObject.GetStartTime())
			res = append(res, apicommon.GaugeFloatValue{
				Value:      duration.Seconds(),
				Attributes: reconcileObject.GetDurationMetricsAttributes(),
			})
		}
	}

	return res, nil
}

// func GetDeploymentInterval(ctx context.Context, client client.Client, reconcileObjectList client.ObjectList, previousObject client.Object) ([]apicommon.GaugeFloatValue, error) {
// 	err := client.List(ctx, reconcileObjectList)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not retrieve workload instances: %w", err)
// 	}

// 	piWrapper, err := NewListItemWrapperFromClientObjectList(reconcileObjectList)
// 	if err != nil {
// 		return nil, err
// 	}

// 	piWrapper2, err := NewPhaseItemWrapperFromClientObject(previousObject)
// 	if err != nil {
// 		return nil, err
// 	}

// 	res := []common.GaugeFloatValue{}
// 	for _, workloadInstance := range piWrapper.GetItems() {
// 		if workloadInstance.Spec.PreviousVersion != "" {
// 			err := r.Get(ctx, types.NamespacedName{Name: fmt.Sprintf("%s-%s", workloadInstance.Spec.WorkloadName, workloadInstance.Spec.PreviousVersion), Namespace: workloadInstance.Namespace}, previousObject)
// 			if err != nil {
// 				r.Log.Error(err, "Previous WorkloadInstance not found")
// 			} else if workloadInstance.IsEndTimeSet() {
// 				previousInterval := workloadInstance.GetEndTime().Sub(piWrapper2.GetStartTime())
// 				res = append(res, common.GaugeFloatValue{
// 					Value:      previousInterval.Seconds(),
// 					Attributes: workloadInstance.GetDurationMetricsAttributes(),
// 				})
// 			}
// 		}
// 	}
// 	return res, nil
// }
