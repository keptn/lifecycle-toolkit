package common

import (
	"context"
	"fmt"

	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/interfaces"
	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ObserveDeploymentDuration(ctx context.Context, client client.Client, reconcileObjectList client.ObjectList, gauge asyncfloat64.Gauge) error {
	err := client.List(ctx, reconcileObjectList)
	if err != nil {
		return fmt.Errorf(controllererrors.ErrCannotRetrieveInstancesMsg, err)
	}

	piWrapper, err := interfaces.NewListItemWrapperFromClientObjectList(reconcileObjectList)
	if err != nil {
		return err
	}

	for _, ro := range piWrapper.GetItems() {
		reconcileObject, _ := interfaces.NewMetricsObjectWrapperFromClientObject(ro)
		if reconcileObject.IsEndTimeSet() {
			duration := reconcileObject.GetEndTime().Sub(reconcileObject.GetStartTime())
			gauge.Observe(ctx, duration.Seconds(), reconcileObject.GetDurationMetricsAttributes()...)
		}
	}

	return nil
}

func ObserveDeploymentInterval(ctx context.Context, client client.Client, reconcileObjectList client.ObjectList, previousObject client.Object, gauge asyncfloat64.Gauge) error {
	err := client.List(ctx, reconcileObjectList)
	if err != nil {
		return fmt.Errorf(controllererrors.ErrCannotRetrieveInstancesMsg, err)
	}

	piWrapper, err := interfaces.NewListItemWrapperFromClientObjectList(reconcileObjectList)
	if err != nil {
		return err
	}

	for _, ro := range piWrapper.GetItems() {
		reconcileObject, _ := interfaces.NewMetricsObjectWrapperFromClientObject(ro)
		if reconcileObject.GetPreviousVersion() != "" {
			err := client.Get(ctx, types.NamespacedName{Name: fmt.Sprintf("%s-%s", reconcileObject.GetParentName(), reconcileObject.GetPreviousVersion()), Namespace: reconcileObject.GetNamespace()}, previousObject)
			if err != nil {
				return nil
			}
			piWrapper2, err := interfaces.NewMetricsObjectWrapperFromClientObject(previousObject)
			if err != nil {
				return err
			}
			if reconcileObject.IsEndTimeSet() {
				previousInterval := reconcileObject.GetEndTime().Sub(piWrapper2.GetStartTime())
				gauge.Observe(ctx, previousInterval.Seconds(), reconcileObject.GetDurationMetricsAttributes()...)
			}
		}
	}

	return nil
}

func ObserveActiveInstances(ctx context.Context, client client.Client, reconcileObjectList client.ObjectList, gauge asyncint64.Gauge) error {
	err := client.List(ctx, reconcileObjectList)
	if err != nil {
		return fmt.Errorf(controllererrors.ErrCannotRetrieveInstancesMsg, err)
	}

	piWrapper, err := interfaces.NewListItemWrapperFromClientObjectList(reconcileObjectList)
	if err != nil {
		return err
	}

	for _, ro := range piWrapper.GetItems() {
		activeMetricsObject, _ := interfaces.NewActiveMetricsObjectWrapperFromClientObject(ro)
		gaugeValue := int64(0)
		if !activeMetricsObject.IsEndTimeSet() {
			gaugeValue = int64(1)
		}

		gauge.Observe(ctx, gaugeValue, activeMetricsObject.GetActiveMetricsAttributes()...)
	}

	return nil
}
