package eventfilter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
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
