// nolint: dupl
package keptnwebhookcontroller

import (
	"context"
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/require"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	apiv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestLabelSelectorRetriever_GetCRDs(t *testing.T) {
	crd1 := &apiv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "crd1",
			Labels: map[string]string{
				"foo": "bar",
			},
		},
		Spec: apiv1.CustomResourceDefinitionSpec{},
	}
	crd2 := &apiv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "crd2",
			Labels: map[string]string{
				"foo": "foo",
			},
		},
		Spec: apiv1.CustomResourceDefinitionSpec{},
	}
	scheme := runtime.NewScheme()
	err := apiv1.AddToScheme(scheme)
	require.Nil(t, err)

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(crd1, crd2).Build()

	retriever := &LabelSelectorRetriever{
		MatchLabels: map[string]string{
			"foo": "bar",
		},
		Client: fakeClient,
	}

	crds, err := retriever.GetCRDs(context.Background())

	require.Nil(t, err)

	require.NotNil(t, crds)
	require.Len(t, crds.Items, 1)
	require.Equal(t, "crd1", crds.Items[0].Name)
}

func TestLabelSelectorRetriever_GetCRDs_ReturnEmptyListIfNothingMatches(t *testing.T) {
	crd1 := &apiv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "crd1",
		},
		Spec: apiv1.CustomResourceDefinitionSpec{},
	}
	scheme := runtime.NewScheme()
	err := apiv1.AddToScheme(scheme)
	require.Nil(t, err)

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(crd1).Build()

	retriever := &LabelSelectorRetriever{
		MatchLabels: map[string]string{
			"foo": "bar",
		},
		Client: fakeClient,
	}

	crds, err := retriever.GetCRDs(context.Background())

	require.Nil(t, err)

	require.NotNil(t, crds)
	require.Empty(t, crds.Items)
}

func TestLabelSelectorRetriever_GetMutatingWebhooks(t *testing.T) {
	mwh1 := &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mwh1",
			Labels: map[string]string{
				"foo": "bar",
			},
		},
	}
	mwh2 := &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mwh2",
			Labels: map[string]string{
				"foo": "foo",
			},
		},
	}
	scheme := runtime.NewScheme()
	err := admissionregistrationv1.AddToScheme(scheme)
	require.Nil(t, err)

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(mwh1, mwh2).Build()

	retriever := &LabelSelectorRetriever{
		MatchLabels: map[string]string{
			"foo": "bar",
		},
		Client: fakeClient,
	}

	mwhs, err := retriever.GetMutatingWebhooks(context.Background())

	require.Nil(t, err)

	require.NotNil(t, mwhs)
	require.Len(t, mwhs.Items, 1)
	require.Equal(t, "mwh1", mwhs.Items[0].Name)
}

func TestLabelSelectorRetriever_GetMutatingWebhook_ReturnEmptyListIfNothingMatches(t *testing.T) {
	mwh1 := &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mwh1",
		},
	}

	scheme := runtime.NewScheme()
	err := admissionregistrationv1.AddToScheme(scheme)
	require.Nil(t, err)

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(mwh1).Build()

	retriever := &LabelSelectorRetriever{
		MatchLabels: map[string]string{
			"foo": "bar",
		},
		Client: fakeClient,
	}

	mwhs, err := retriever.GetMutatingWebhooks(context.Background())

	require.Nil(t, err)

	require.NotNil(t, mwhs)
	require.Empty(t, mwhs.Items)
}

func TestLabelSelectorRetriever_GetValidatingWebhooks(t *testing.T) {
	vwh1 := &admissionregistrationv1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "vwh1",
			Labels: map[string]string{
				"foo": "bar",
			},
		},
	}
	vwh2 := &admissionregistrationv1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "vwh2",
			Labels: map[string]string{
				"foo": "foo",
			},
		},
	}
	scheme := runtime.NewScheme()
	err := admissionregistrationv1.AddToScheme(scheme)
	require.Nil(t, err)

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(vwh1, vwh2).Build()

	retriever := &LabelSelectorRetriever{
		MatchLabels: map[string]string{
			"foo": "bar",
		},
		Client: fakeClient,
	}

	vwhs, err := retriever.GetValidatingWebhooks(context.Background())

	require.Nil(t, err)

	require.NotNil(t, vwhs)
	require.Len(t, vwhs.Items, 1)
	require.Equal(t, "vwh1", vwhs.Items[0].Name)
}

func TestLabelSelectorRetriever_GetValidatingWebhook_ReturnEmptyListIfNothingMatches(t *testing.T) {
	vwh1 := &admissionregistrationv1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "vwh1",
		},
	}

	scheme := runtime.NewScheme()
	err := admissionregistrationv1.AddToScheme(scheme)
	require.Nil(t, err)

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(vwh1).Build()

	retriever := &LabelSelectorRetriever{
		MatchLabels: map[string]string{
			"foo": "bar",
		},
		Client: fakeClient,
	}

	vwhs, err := retriever.GetValidatingWebhooks(context.Background())

	require.Nil(t, err)

	require.NotNil(t, vwhs)
	require.Empty(t, vwhs.Items)
}

func TestNewResourceRetriever(t *testing.T) {
	type args struct {
		config CertificateReconcilerConfig
	}
	tests := []struct {
		name string
		args args
		want IResourceRetriever
	}{
		{
			name: "label selector retriever",
			args: args{
				config: CertificateReconcilerConfig{
					WatchResources: nil,
					MatchLabels:    nil,
				},
			},
			want: &LabelSelectorRetriever{},
		},
		{
			name: "resource name retriever",
			args: args{
				config: CertificateReconcilerConfig{
					WatchResources: &ObservedObjects{},
					MatchLabels:    nil,
				},
			},
			want: &ResourceNameRetriever{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retriever := NewResourceRetriever(tt.args.config)
			require.IsType(t, tt.want, retriever)
		})
	}
}

func TestResourceNameRetriever_GetCRDs(t *testing.T) {
	crd1 := &apiv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "crd1",
			Labels: map[string]string{
				"foo": "bar",
			},
		},
		Spec: apiv1.CustomResourceDefinitionSpec{},
	}
	crd2 := &apiv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "crd2",
			Labels: map[string]string{
				"foo": "foo",
			},
		},
		Spec: apiv1.CustomResourceDefinitionSpec{},
	}
	scheme := runtime.NewScheme()
	err := apiv1.AddToScheme(scheme)
	require.Nil(t, err)

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(crd1, crd2).Build()

	retriever := &ResourceNameRetriever{
		Client: fakeClient,
		WatchResources: ObservedObjects{
			CustomResourceDefinitions: []string{"crd1", "crd2"},
		},
		Log: testr.New(t),
	}

	crds, err := retriever.GetCRDs(context.TODO())

	require.Nil(t, err)

	require.Len(t, crds.Items, 2)
}

func TestResourceNameRetriever_GetCRDs_ContinueIfOneNotFound(t *testing.T) {
	crd1 := &apiv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "crd1",
			Labels: map[string]string{
				"foo": "bar",
			},
		},
		Spec: apiv1.CustomResourceDefinitionSpec{},
	}
	crd2 := &apiv1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "crd2",
			Labels: map[string]string{
				"foo": "foo",
			},
		},
		Spec: apiv1.CustomResourceDefinitionSpec{},
	}
	scheme := runtime.NewScheme()
	err := apiv1.AddToScheme(scheme)
	require.Nil(t, err)

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(crd1, crd2).Build()

	retriever := &ResourceNameRetriever{
		Client: fakeClient,
		WatchResources: ObservedObjects{
			CustomResourceDefinitions: []string{"crdx", "crd2"},
		},
		Log: testr.New(t),
	}

	crds, err := retriever.GetCRDs(context.TODO())

	require.Nil(t, err)

	require.Len(t, crds.Items, 1)
}

func TestResourceNameRetriever_GetMutatingWebhooks(t *testing.T) {
	mwh1 := &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mwh1",
		},
	}
	mwh2 := &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mwh2",
		},
	}
	scheme := runtime.NewScheme()
	err := admissionregistrationv1.AddToScheme(scheme)
	require.Nil(t, err)

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(mwh1, mwh2).Build()

	retriever := &ResourceNameRetriever{
		Client: fakeClient,
		WatchResources: ObservedObjects{
			MutatingWebhooks: []string{"mwh1", "mwh2"},
		},
		Log: testr.New(t),
	}

	mwhs, err := retriever.GetMutatingWebhooks(context.TODO())

	require.Nil(t, err)

	require.Len(t, mwhs.Items, 2)
}

func TestResourceNameRetriever_GetMutatingWebhooks_ContinueIfOneNotFound(t *testing.T) {
	mwh1 := &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mwh1",
			Labels: map[string]string{
				"foo": "bar",
			},
		},
	}
	mwh2 := &admissionregistrationv1.MutatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "mwh2",
			Labels: map[string]string{
				"foo": "foo",
			},
		},
	}
	scheme := runtime.NewScheme()
	err := admissionregistrationv1.AddToScheme(scheme)
	require.Nil(t, err)

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(mwh1, mwh2).Build()

	retriever := &ResourceNameRetriever{
		Client: fakeClient,
		WatchResources: ObservedObjects{
			MutatingWebhooks: []string{"mwhx", "mwh2"},
		},
		Log: testr.New(t),
	}

	mwhs, err := retriever.GetMutatingWebhooks(context.TODO())

	require.Nil(t, err)

	require.Len(t, mwhs.Items, 1)
}

func TestResourceNameRetriever_GetValidatingWebhooks(t *testing.T) {
	vwh1 := &admissionregistrationv1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "vwh1",
		},
	}
	vwh2 := &admissionregistrationv1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "vwh2",
		},
	}
	scheme := runtime.NewScheme()
	err := admissionregistrationv1.AddToScheme(scheme)
	require.Nil(t, err)

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(vwh1, vwh2).Build()

	retriever := &ResourceNameRetriever{
		Client: fakeClient,
		WatchResources: ObservedObjects{
			ValidatingWebhooks: []string{"vwh1", "vwh2"},
		},
		Log: testr.New(t),
	}

	vwhs, err := retriever.GetValidatingWebhooks(context.TODO())

	require.Nil(t, err)

	require.Len(t, vwhs.Items, 2)
}

func TestResourceNameRetriever_GetValidatingWebhooks_ContinueIfOneNotFound(t *testing.T) {
	mwh1 := &admissionregistrationv1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "vwh1",
		},
	}
	mwh2 := &admissionregistrationv1.ValidatingWebhookConfiguration{
		ObjectMeta: metav1.ObjectMeta{
			Name: "vwh2",
		},
	}
	scheme := runtime.NewScheme()
	err := admissionregistrationv1.AddToScheme(scheme)
	require.Nil(t, err)

	fakeClient := fake.NewClientBuilder().WithScheme(scheme).WithObjects(mwh1, mwh2).Build()

	retriever := &ResourceNameRetriever{
		Client: fakeClient,
		WatchResources: ObservedObjects{
			ValidatingWebhooks: []string{"vwhx", "vwh2"},
		},
		Log: testr.New(t),
	}

	vwhs, err := retriever.GetValidatingWebhooks(context.TODO())

	require.Nil(t, err)

	require.Len(t, vwhs.Items, 1)
}
