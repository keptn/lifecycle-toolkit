package kubeobjects

import (
	"context"
	"github.com/go-logr/logr/testr"
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/controllers/common/fake"
)

func TestKubeQuery(t *testing.T) {
	fakeClient := fake.NewClient()
	_ = newKubeQuery(context.TODO(), fakeClient, fakeClient, testr.New(t))
}
