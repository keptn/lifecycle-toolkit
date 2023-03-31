package keptnappcreationrequest

import (
	"context"
	"testing"

	"github.com/go-logr/logr"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	controllerruntime "sigs.k8s.io/controller-runtime"
	k8sfake "sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestKeptnAppCreationRequestReconciler_Reconcile(t *testing.T) {
	fakeClient := k8sfake.NewClientBuilder().WithObjects().Build()

	err := klcv1alpha3.AddToScheme(fakeClient.Scheme())
	require.Nil(t, err)

	kacr := &klcv1alpha3.KeptnAppCreationRequest{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-kacr",
			Namespace: "my-namespace",
		},
		Spec: klcv1alpha3.KeptnAppCreationRequestSpec{
			AppName: "my-app",
		},
	}

	err = fakeClient.Create(context.TODO(), kacr)
	require.Nil(t, err)

	r := &KeptnAppCreationRequestReconciler{
		Client: fakeClient,
		Scheme: fakeClient.Scheme(),
		Log:    logr.Logger{},
	}

	_, err = r.Reconcile(context.Background(), controllerruntime.Request{
		NamespacedName: types.NamespacedName{
			Namespace: kacr.Namespace,
			Name:      kacr.Name,
		},
	})

	require.Nil(t, err)
}
