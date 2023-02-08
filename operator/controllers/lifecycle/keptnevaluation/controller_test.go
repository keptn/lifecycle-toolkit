package keptnevaluation

import (
	"context"
	"testing"

	"github.com/go-logr/logr/testr"
	klcv1alpha2 "github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha2"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/providers"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const KltNamespace = "klt-namespace"

func TestKeptnEvaluationReconciler_fetchDefinitionAndProvider(t *testing.T) {

	metricEvalDef, DTEvalDef, PromEvalDef, EvalDef := setupEvalDefinitions()
	DTProv, PromProv := setupProviders()
	client := fake.NewClient(metricEvalDef, DTEvalDef, PromEvalDef, EvalDef, DTProv, PromProv)

	r := &KeptnEvaluationReconciler{
		Client: client,
		Scheme: client.Scheme(),
		Log:    testr.New(t),
	}

	tests := []struct {
		name                 string
		namespacedDefinition types.NamespacedName
		wantDef              *klcv1alpha2.KeptnEvaluationDefinition
		wantProv             *klcv1alpha2.KeptnEvaluationProvider
		wantErr              bool
	}{
		{
			name: "keptn metrics",
			namespacedDefinition: types.NamespacedName{
				Namespace: KltNamespace,
				Name:      "myKeptn",
			},
			wantDef:  metricEvalDef,
			wantProv: providers.GetDefaultMetricProvider(KltNamespace),
		},
		{
			name: "DT metrics",
			namespacedDefinition: types.NamespacedName{
				Namespace: KltNamespace,
				Name:      "myDT",
			},
			wantDef:  DTEvalDef,
			wantProv: DTProv,
		},

		{
			name: "Prometheus metrics",
			namespacedDefinition: types.NamespacedName{
				Namespace: KltNamespace,
				Name:      "myProm",
			},
			wantDef:  PromEvalDef,
			wantProv: PromProv,
		},

		{
			name: "Unexisting Evaluation Def",
			namespacedDefinition: types.NamespacedName{
				Namespace: KltNamespace,
				Name:      "whatever",
			},
			wantDef:  nil,
			wantProv: nil,
			wantErr:  true,
		},
		{
			name: "Unexisting Provider",
			namespacedDefinition: types.NamespacedName{
				Namespace: KltNamespace,
				Name:      "mydef",
			},
			wantDef:  nil,
			wantProv: nil,
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got, got1, err := r.fetchDefinitionAndProvider(context.TODO(), tt.namespacedDefinition)
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchDefinitionAndProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantDef != nil {
				require.Equal(t, got.Name, tt.wantDef.Name)
			} else {
				require.Nil(t, got)
			}

			if tt.wantProv != nil {
				require.Equal(t, got1.Name, tt.wantProv.Name)
			} else {
				require.Nil(t, got1)
			}

		})
	}
}

func setupEvalDefinitions() (*klcv1alpha2.KeptnEvaluationDefinition, *klcv1alpha2.KeptnEvaluationDefinition, *klcv1alpha2.KeptnEvaluationDefinition, *klcv1alpha2.KeptnEvaluationDefinition) {
	metricEvalDef := &klcv1alpha2.KeptnEvaluationDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: KltNamespace,
			Name:      "myKeptn",
		},
		Spec: klcv1alpha2.KeptnEvaluationDefinitionSpec{
			Source:     providers.KeptnMetricProviderName,
			Objectives: nil,
		},
		Status: klcv1alpha2.KeptnEvaluationDefinitionStatus{},
	}

	DTEvalDef := &klcv1alpha2.KeptnEvaluationDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: KltNamespace,
			Name:      "myDT",
		},
		Spec: klcv1alpha2.KeptnEvaluationDefinitionSpec{
			Source:     providers.DynatraceProviderName,
			Objectives: nil,
		},
		Status: klcv1alpha2.KeptnEvaluationDefinitionStatus{},
	}

	PromEvalDef := &klcv1alpha2.KeptnEvaluationDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: KltNamespace,
			Name:      "myProm",
		},
		Spec: klcv1alpha2.KeptnEvaluationDefinitionSpec{
			Source:     providers.PrometheusProviderName,
			Objectives: nil,
		},
		Status: klcv1alpha2.KeptnEvaluationDefinitionStatus{},
	}

	EvalDef := &klcv1alpha2.KeptnEvaluationDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: KltNamespace,
			Name:      "mdef",
		},
		Spec: klcv1alpha2.KeptnEvaluationDefinitionSpec{
			Source:     "dunno",
			Objectives: nil,
		},
		Status: klcv1alpha2.KeptnEvaluationDefinitionStatus{},
	}

	return metricEvalDef, DTEvalDef, PromEvalDef, EvalDef
}

func setupProviders() (*klcv1alpha2.KeptnEvaluationProvider, *klcv1alpha2.KeptnEvaluationProvider) {
	DTProv := &klcv1alpha2.KeptnEvaluationProvider{
		ObjectMeta: metav1.ObjectMeta{
			Name:      providers.DynatraceProviderName,
			Namespace: KltNamespace,
		},
	}

	PromProv := &klcv1alpha2.KeptnEvaluationProvider{
		ObjectMeta: metav1.ObjectMeta{
			Name:      providers.PrometheusProviderName,
			Namespace: KltNamespace,
		},
	}

	return DTProv, PromProv
}
