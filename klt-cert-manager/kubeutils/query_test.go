package kubeutils

import (
	"context"
	"testing"

	"github.com/go-logr/logr/testr"

	"github.com/keptn/lifecycle-toolkit/klt-cert-manager/fake"
)

func TestKubeQuery(t *testing.T) {
	fakeClient := fake.NewClient()
	_ = newKubeQuery(context.TODO(), fakeClient, fakeClient, testr.New(t))
}
