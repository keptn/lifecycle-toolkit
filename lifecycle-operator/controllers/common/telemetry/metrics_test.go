package telemetry

import (
	"context"
	"testing"
	"time"

	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	controllererrors "github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/errors"
	"github.com/keptn/lifecycle-toolkit/lifecycle-operator/controllers/lifecycle/interfaces"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestMetrics_ObserveDeploymentDuration(t *testing.T) {
	gauge := noop.Float64ObservableGauge{}

	tests := []struct {
		name          string
		clientObjects client.ObjectList
		list          client.ObjectList
		err           error
		gauge         metric.Float64ObservableGauge
	}{
		{
			name:          "failed to create wrapper",
			list:          &apilifecycle.KeptnAppList{},
			clientObjects: &apilifecycle.KeptnAppList{},
			err:           controllererrors.ErrCannotWrapToListItem,
			gauge:         nil,
		},
		{
			name: "no endtime set",
			list: &apilifecycle.KeptnAppVersionList{},
			clientObjects: &apilifecycle.KeptnAppVersionList{
				Items: []apilifecycle.KeptnAppVersion{
					{
						Status: apilifecycle.KeptnAppVersionStatus{},
					},
				},
			},
			err:   nil,
			gauge: gauge,
		},
		{
			name: "endtime set",
			list: &apilifecycle.KeptnAppVersionList{},
			clientObjects: &apilifecycle.KeptnAppVersionList{
				Items: []apilifecycle.KeptnAppVersion{
					{
						Spec: apilifecycle.KeptnAppVersionSpec{
							KeptnAppSpec: apilifecycle.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
						Status: apilifecycle.KeptnAppVersionStatus{
							EndTime:   metav1.Time{Time: metav1.Now().Time.Add(5 * time.Second)},
							StartTime: metav1.Time{Time: metav1.Now().Time},
						},
					},
				},
			},
			err:   nil,
			gauge: gauge,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apilifecycle.AddToScheme(scheme.Scheme)
			require.Nil(t, err)
			client := fake.NewClientBuilder().WithLists(tt.clientObjects).Build()
			err = ObserveDeploymentDuration(context.TODO(), client, tt.list, gauge, noop.Observer{})
			require.ErrorIs(t, err, tt.err)
		})

	}
}

func TestMetrics_ObserveActiveInstances(t *testing.T) {
	tests := []struct {
		name          string
		clientObjects client.ObjectList
		list          client.ObjectList
		err           error
	}{
		{
			name:          "failed to create wrapper",
			list:          &apilifecycle.KeptnAppList{},
			clientObjects: &apilifecycle.KeptnAppList{},
			err:           controllererrors.ErrCannotWrapToListItem,
		},
		{
			name: "no endtime set - active instances",
			list: &apilifecycle.KeptnAppVersionList{},
			clientObjects: &apilifecycle.KeptnAppVersionList{
				Items: []apilifecycle.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "namespace",
						},
						Spec: apilifecycle.KeptnAppVersionSpec{
							KeptnAppSpec: apilifecycle.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
						Status: apilifecycle.KeptnAppVersionStatus{},
					},
				},
			},
			err: nil,
		},
		{
			name: "endtime set - no active instances",
			list: &apilifecycle.KeptnAppVersionList{},
			clientObjects: &apilifecycle.KeptnAppVersionList{
				Items: []apilifecycle.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "namespace",
						},
						Spec: apilifecycle.KeptnAppVersionSpec{
							KeptnAppSpec: apilifecycle.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
						Status: apilifecycle.KeptnAppVersionStatus{
							EndTime:   metav1.Time{Time: metav1.Now().Time.Add(5 * time.Second)},
							StartTime: metav1.Time{Time: metav1.Now().Time},
						},
					},
				},
			},
			err: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apilifecycle.AddToScheme(scheme.Scheme)
			require.Nil(t, err)
			client := fake.NewClientBuilder().WithLists(tt.clientObjects).Build()
			gauge := noop.Int64ObservableGauge{}
			require.Nil(t, err)
			err = ObserveActiveInstances(context.TODO(), client, tt.list, gauge, noop.Observer{})
			require.ErrorIs(t, err, tt.err)

		})

	}
}

func TestMetrics_ObserveDeploymentInterval(t *testing.T) {
	tests := []struct {
		name          string
		clientObjects client.ObjectList
		list          client.ObjectList
		previous      client.Object
		err           error
	}{
		{
			name:          "failed to create wrapper",
			list:          &apilifecycle.KeptnAppList{},
			clientObjects: &apilifecycle.KeptnAppList{},
			err:           controllererrors.ErrCannotWrapToListItem,
		},
		{
			name: "no previous version",
			list: &apilifecycle.KeptnAppVersionList{},
			clientObjects: &apilifecycle.KeptnAppVersionList{
				Items: []apilifecycle.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "namespace",
						},
						Spec: apilifecycle.KeptnAppVersionSpec{
							KeptnAppSpec: apilifecycle.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "",
						},
					},
				},
			},
			err: nil,
		},
		{
			name: "previous version - no previous object",
			list: &apilifecycle.KeptnAppVersionList{},
			clientObjects: &apilifecycle.KeptnAppVersionList{
				Items: []apilifecycle.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Namespace: "namespace",
						},
						Spec: apilifecycle.KeptnAppVersionSpec{
							KeptnAppSpec: apilifecycle.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
					},
				},
			},
			err: nil,
		},
		{
			name:     "previous version - object found but no endtime",
			list:     &apilifecycle.KeptnAppVersionList{},
			previous: &apilifecycle.KeptnAppVersion{},
			clientObjects: &apilifecycle.KeptnAppVersionList{
				Items: []apilifecycle.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "appName-version",
							Namespace: "namespace",
						},
						Spec: apilifecycle.KeptnAppVersionSpec{
							KeptnAppSpec: apilifecycle.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "appName-previousVersion",
							Namespace: "namespace",
						},
						Spec: apilifecycle.KeptnAppVersionSpec{
							KeptnAppSpec: apilifecycle.KeptnAppSpec{
								Version: "previousVersion",
							},
							AppName:         "appName",
							PreviousVersion: "",
						},
					},
				},
			},
			err: nil,
		},
		{
			name:     "previous version - object found with endtime",
			list:     &apilifecycle.KeptnAppVersionList{},
			previous: &apilifecycle.KeptnAppVersion{},
			clientObjects: &apilifecycle.KeptnAppVersionList{
				Items: []apilifecycle.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "appName-version",
							Namespace: "namespace",
						},
						Spec: apilifecycle.KeptnAppVersionSpec{
							KeptnAppSpec: apilifecycle.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
						Status: apilifecycle.KeptnAppVersionStatus{
							EndTime:   metav1.Time{Time: metav1.Now().Time.Add(5 * time.Second)},
							StartTime: metav1.Time{Time: metav1.Now().Time},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:      "appName-previousVersion",
							Namespace: "namespace",
						},
						Spec: apilifecycle.KeptnAppVersionSpec{
							KeptnAppSpec: apilifecycle.KeptnAppSpec{
								Version: "previousVersion",
							},
							AppName:         "appName",
							PreviousVersion: "",
						},
						Status: apilifecycle.KeptnAppVersionStatus{
							EndTime:   metav1.Time{Time: metav1.Now().Time.Add(5 * time.Second)},
							StartTime: metav1.Time{Time: metav1.Now().Time},
						},
					},
				},
			},
			err: nil,
		},
		{
			name:     "previous version - object found with endtime and revision",
			list:     &apilifecycle.KeptnAppVersionList{},
			previous: &apilifecycle.KeptnAppVersion{},
			clientObjects: &apilifecycle.KeptnAppVersionList{
				Items: []apilifecycle.KeptnAppVersion{
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:       "appName-version-1",
							Namespace:  "namespace",
							Generation: 1,
						},
						Spec: apilifecycle.KeptnAppVersionSpec{
							KeptnAppSpec: apilifecycle.KeptnAppSpec{
								Version: "version",
							},
							AppName:         "appName",
							PreviousVersion: "previousVersion",
						},
						Status: apilifecycle.KeptnAppVersionStatus{
							EndTime:   metav1.Time{Time: metav1.Now().Time.Add(5 * time.Second)},
							StartTime: metav1.Time{Time: metav1.Now().Time},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:       "appName-previousVersion-2",
							Namespace:  "namespace",
							Generation: 2,
						},
						Spec: apilifecycle.KeptnAppVersionSpec{
							KeptnAppSpec: apilifecycle.KeptnAppSpec{
								Version: "previousVersion",
							},
							AppName:         "appName",
							PreviousVersion: "",
						},
						Status: apilifecycle.KeptnAppVersionStatus{
							EndTime:   metav1.Time{Time: metav1.Now().Time.Add(5 * time.Second)},
							StartTime: metav1.Time{Time: metav1.Now().Time},
						},
					},
					{
						ObjectMeta: metav1.ObjectMeta{
							Name:       "appName-previousVersion-1",
							Namespace:  "namespace",
							Generation: 1,
						},
						Spec: apilifecycle.KeptnAppVersionSpec{
							KeptnAppSpec: apilifecycle.KeptnAppSpec{
								Version: "previousVersion",
							},
							AppName:         "appName",
							PreviousVersion: "",
						},
						Status: apilifecycle.KeptnAppVersionStatus{
							EndTime:   metav1.Time{Time: metav1.Now().Time.Add(5 * time.Second)},
							StartTime: metav1.Time{Time: metav1.Now().Time},
						},
					},
				},
			},
			err: nil,
		},
	}

	gauge := noop.Float64ObservableGauge{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := apilifecycle.AddToScheme(scheme.Scheme)
			require.Nil(t, err)
			fakeClient := fake.NewClientBuilder().WithLists(tt.clientObjects).Build()
			err = ObserveDeploymentInterval(context.TODO(), fakeClient, tt.list, gauge, noop.Observer{})
			require.ErrorIs(t, err, tt.err)
		})

	}
}

func TestGetPredecessor(t *testing.T) {
	now := time.Now()
	appVersions := &apilifecycle.KeptnAppVersionList{
		Items: []apilifecycle.KeptnAppVersion{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "my-app-1.0.0-1",
					Generation: 1,
				},
				Spec: apilifecycle.KeptnAppVersionSpec{
					KeptnAppSpec: apilifecycle.KeptnAppSpec{
						Version: "1.0.0",
					},
					AppName: "my-app",
				},
				Status: apilifecycle.KeptnAppVersionStatus{
					StartTime: metav1.NewTime(now),
					EndTime:   metav1.NewTime(now.Add(10 * time.Second)),
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "my-app-1.0.0-2",
					Generation: 2,
				},
				Spec: apilifecycle.KeptnAppVersionSpec{
					KeptnAppSpec: apilifecycle.KeptnAppSpec{
						Version: "1.0.0",
					},
					AppName: "my-app",
				},
				Status: apilifecycle.KeptnAppVersionStatus{
					StartTime: metav1.NewTime(now.Add(1 * time.Second)),
					EndTime:   metav1.NewTime(now.Add(10 * time.Second)),
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:       "my-app-1.1.0-1",
					Generation: 1,
				},
				Spec: apilifecycle.KeptnAppVersionSpec{
					KeptnAppSpec: apilifecycle.KeptnAppSpec{
						Version: "1.1.0",
					},
					AppName:         "my-app",
					PreviousVersion: "1.0.0",
				},
				Status: apilifecycle.KeptnAppVersionStatus{
					StartTime: metav1.NewTime(now),
					EndTime:   metav1.NewTime(now.Add(10 * time.Second)),
				},
			},
		},
	}

	appVersionsWrapper, err := interfaces.NewListItemWrapperFromClientObjectList(appVersions)
	require.Nil(t, err)

	latestAppVersion, err := interfaces.NewMetricsObjectWrapperFromClientObject(appVersionsWrapper.GetItems()[2])
	require.Nil(t, err)

	require.Equal(t, "1.0.0", latestAppVersion.GetPreviousVersion())

	predecessor := getPredecessor(latestAppVersion, appVersionsWrapper.GetItems())

	require.Equal(t, "", predecessor.GetPreviousVersion())

	expectedPredecessor, err := interfaces.NewMetricsObjectWrapperFromClientObject(appVersionsWrapper.GetItems()[0])
	require.Nil(t, err)

	require.Equal(t, expectedPredecessor, predecessor)
}
