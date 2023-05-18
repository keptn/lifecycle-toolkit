package common

import (
	"context"
	"fmt"
	"strings"

	controllererrors "github.com/keptn/lifecycle-toolkit/operator/controllers/errors"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/lifecycle/interfaces"
	"go.opentelemetry.io/otel/metric"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ObserveDeploymentDuration(ctx context.Context, client client.Client, reconcileObjectList client.ObjectList, gauge metric.Float64ObservableGauge, o metric.Observer) error {
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
			o.ObserveFloat64(gauge, duration.Seconds(), metric.WithAttributes(reconcileObject.GetDurationMetricsAttributes()...))
		}
	}

	return nil
}

func ObserveDeploymentInterval(ctx context.Context, client client.Client, reconcileObjectList client.ObjectList, gauge metric.Float64ObservableGauge, o metric.Observer) error {
	err := client.List(ctx, reconcileObjectList)
	if err != nil {
		return fmt.Errorf(controllererrors.ErrCannotRetrieveInstancesMsg, err)
	}

	piWrapper, err := interfaces.NewListItemWrapperFromClientObjectList(reconcileObjectList)
	if err != nil {
		return err
	}

	items := piWrapper.GetItems()
	for index := 0; index < len(items); index++ {
		ro := items[index]
		reconcileObject, _ := interfaces.NewMetricsObjectWrapperFromClientObject(ro)
		if reconcileObject.GetPreviousVersion() != "" {
			if !reconcileObject.IsEndTimeSet() {
				continue
			}
			predecessor := getPredecessor(reconcileObject, items)
			if predecessor == nil {
				continue
			}

			previousInterval := reconcileObject.GetEndTime().Sub(predecessor.GetEndTime())
			o.ObserveFloat64(gauge, previousInterval.Seconds(), metric.WithAttributes(reconcileObject.GetDurationMetricsAttributes()...))
		}
	}

	return nil
}

func getPredecessor(successor *interfaces.MetricsObjectWrapper, items []client.Object) interfaces.MetricsObject {
	var predecessor interfaces.MetricsObject
	for i := 0; i < len(items); i++ {
		// to calculate the interval, we take the earliest revision of the previous version as a reference
		if strings.HasPrefix(items[i].GetName(), fmt.Sprintf("%s-%s", successor.GetParentName(), successor.GetPreviousVersion())) {
			predecessorCandidate, err := interfaces.NewMetricsObjectWrapperFromClientObject(items[i])
			if err != nil {
				// continue with the other items
				continue
			}
			if predecessor == nil || predecessorCandidate.GetStartTime().Before(predecessor.GetStartTime()) {
				predecessor = predecessorCandidate
			}
		}
	}
	return predecessor
}

func ObserveActiveInstances(ctx context.Context, client client.Client, reconcileObjectList client.ObjectList, gauge metric.Int64ObservableGauge, o metric.Observer) error {
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

		o.ObserveInt64(gauge, gaugeValue, metric.WithAttributes(activeMetricsObject.GetActiveMetricsAttributes()...))
	}

	return nil
}
