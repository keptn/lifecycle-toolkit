package eventfilter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
