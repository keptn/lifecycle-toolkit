package eventfilter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

const (
	testName1 = "test-name-1"
	testName2 = "test-name-2"

	testNamespace1 = "test-namespace-1"
	testNamespace2 = "test-namespace-2"
)

func TestForObjectNameAndNamespace(t *testing.T) {
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName1,
			Namespace: testNamespace1,
		},
	}
	predicate := ForObjectNameAndNamespace(testName1, testNamespace1)

	assert.True(t, predicate.Generic(event.GenericEvent{
		Object: deployment,
	}))

	deployment = &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName2,
			Namespace: testNamespace1,
		},
	}

	assert.False(t, predicate.Generic(event.GenericEvent{
		Object: deployment,
	}))

	deployment = &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName1,
			Namespace: testNamespace2,
		},
	}

	assert.False(t, predicate.Generic(event.GenericEvent{
		Object: deployment,
	}))

	deployment = &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName2,
			Namespace: testNamespace2,
		},
	}

	assert.False(t, predicate.Generic(event.GenericEvent{
		Object: deployment,
	}))
}

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

func TestForName(t *testing.T) {
	deployment := &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName1,
			Namespace: testNamespace1,
		},
	}

	assert.True(t, hasName(deployment, testName1))
	assert.False(t, hasName(deployment, testName2))

	deployment = &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName2,
			Namespace: testNamespace1,
		},
	}

	assert.False(t, hasName(deployment, testName1))
	assert.True(t, hasName(deployment, testName2))

	deployment = &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName1,
			Namespace: testNamespace2,
		},
	}

	assert.True(t, hasName(deployment, testName1))
	assert.False(t, hasName(deployment, testName2))

	deployment = &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      testName2,
			Namespace: testNamespace2,
		},
	}

	assert.False(t, hasName(deployment, testName1))
	assert.True(t, hasName(deployment, testName2))
}
