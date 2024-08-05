package eventfilter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

const (
	testName1 = "test-name-1"
	testName2 = "test-name-2"

	testNamespace1 = "test-namespace-1"
	testNamespace2 = "test-namespace-2"
)

//nolint:dupl
func TestForNamespace(t *testing.T) {
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName1,
			Namespace: testNamespace1,
		},
	}

	assert.True(t, isInNamespace(deployment, testNamespace1))
	assert.False(t, isInNamespace(deployment, testNamespace2))

	deployment = &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName2,
			Namespace: testNamespace1,
		},
	}

	assert.True(t, isInNamespace(deployment, testNamespace1))
	assert.False(t, isInNamespace(deployment, testNamespace2))

	deployment = &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName1,
			Namespace: testNamespace2,
		},
	}

	assert.False(t, isInNamespace(deployment, testNamespace1))
	assert.True(t, isInNamespace(deployment, testNamespace2))

	deployment = &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName2,
			Namespace: testNamespace2,
		},
	}

	assert.False(t, isInNamespace(deployment, testNamespace1))
	assert.True(t, isInNamespace(deployment, testNamespace2))
}

func Test_matchesName(t *testing.T) {
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-deployment",
		},
	}

	assert.True(t, matchesName(deployment, []string{"my-deployment"}))

	deployment = &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: "my-other-deployment",
		},
	}

	assert.False(t, matchesName(deployment, []string{"my-deployment"}))
}

func Test_matchesLabels(t *testing.T) {
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Labels: map[string]string{
				"app": "test",
				"env": "dev",
			},
		},
	}

	firstSelector := labels.SelectorFromSet(labels.Set{
		"app": "test",
		"env": "dev", // different value for 'env'
	})

	assert.True(t, matchesLabels(deployment, firstSelector))

	secondSelectors := labels.SelectorFromSet(labels.Set{
		"app": "test",
		"env": "prod",
	})

	assert.False(t, matchesLabels(deployment, secondSelectors))

	deploymentNoLabels := &v1.Deployment{}

	assert.False(t, matchesLabels(deploymentNoLabels, firstSelector))
}

func TestForLabelsAndNamespace(t *testing.T) {
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName1,
			Namespace: testNamespace1,
			Labels: map[string]string{
				"app": "test",
				"env": "dev",
			},
		},
	}

	selector := labels.SelectorFromSet(labels.Set{
		"app": "test",
		"env": "dev",
	})

	// when the deployments matched with labels and is present in the required namespace.
	assert.True(t, ForLabelsAndNamespace(selector, testNamespace1).Generic(event.GenericEvent{Object: deployment}))

	// when the namespace doesn't match.
	assert.False(t, ForLabelsAndNamespace(selector, "another-namespace").Generic(event.GenericEvent{Object: deployment}))

	// when the labels don't match.
	deploymentNoLabels := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName1,
			Namespace: testNamespace1,
		},
	}
	assert.False(t, ForLabelsAndNamespace(selector, testNamespace1).Generic(event.GenericEvent{Object: deploymentNoLabels}))
}

func TestForNamesAndNamespace(t *testing.T) {
	deployment1 := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-1",
			Namespace: "test-namespace-1",
		},
	}

	deployment2 := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment-2",
			Namespace: "test-namespace-2",
		},
	}

	names := []string{"deployment-1", "deployment-2"}
	namespace := "test-namespace-1"

	// Test deployments in the same namespace and with matching names
	pred := ForNamesAndNamespace(names, namespace)
	assert.True(t, pred.Generic(event.GenericEvent{Object: deployment1}))
	assert.False(t, pred.Generic(event.GenericEvent{Object: deployment2}))

	// Test deployments in different namespace
	namespace = "test-namespace-2"
	pred = ForNamesAndNamespace(names, namespace)
	assert.False(t, pred.Generic(event.GenericEvent{Object: deployment1}))
	assert.True(t, pred.Generic(event.GenericEvent{Object: deployment2}))

	// Test deployments with mismatched names
	names = []string{"deployment-3", "deployment-4"}
	namespace = "test-namespace-1"
	pred = ForNamesAndNamespace(names, namespace)
	assert.False(t, pred.Generic(event.GenericEvent{Object: deployment1}))
	assert.False(t, pred.Generic(event.GenericEvent{Object: deployment2}))
}
