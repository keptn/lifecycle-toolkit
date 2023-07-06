package common

import (
	"fmt"
	"strings"
	"testing"

	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3"
	"github.com/keptn/lifecycle-toolkit/operator/apis/lifecycle/v1alpha3/common"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
)

func TestEventSender_SendK8sEvent(t *testing.T) {
	fakeRecorder := record.NewFakeRecorder(100)
	eventSender := NewEventSender(fakeRecorder)

	eventSender.SendK8sEvent(common.PhaseAppDeployment, "pre-event", &v1alpha3.KeptnAppVersion{
		ObjectMeta: v1.ObjectMeta{
			Name:      "app",
			Namespace: "ns",
		},
	}, "reason-short", "reason-long", "ver1")

	event := <-fakeRecorder.Events

	require.Equal(t, strings.Contains(event, fmt.Sprintf("%s: reason-long / Namespace: ns, Name: app, Version: ver1", common.PhaseAppDeployment.LongName)), true)

}
