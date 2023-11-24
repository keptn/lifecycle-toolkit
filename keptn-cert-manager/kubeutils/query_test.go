package kubeutils

import (
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/keptn/lifecycle-toolkit/keptn-cert-manager/fake"
)

func TestKubeQuery(t *testing.T) {
	fakeClient := fake.NewClient()
	_ = newKubeQuery(fakeClient, fakeClient, testr.New(t))
}
