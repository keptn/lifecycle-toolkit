package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

const sliContent = `
spec_version: "1.0"
indicators:
  throughput: "builtin:service.requestCount.total:merge(0):count?scope=tag(keptn_project:$PROJECT),tag(keptn_stage:$STAGE),tag(keptn_service:$SERVICE),tag(keptn_deployment:$DEPLOYMENT)"
  response_time_p95: "builtin:service.response.time:merge(0):percentile(95)?scope=tag(keptn_project:$PROJECT),tag(keptn_stage:$STAGE),tag(keptn_service:$SERVICE),tag(keptn_deployment:$DEPLOYMENT)"`

const expectedOutput = `---
kind: AnalysisValueTemplate
apiVersion: metrics.keptn.sh/v1alpha3
metadata:
  name: throughput
  creationTimestamp: null
spec:
  provider:
    name: dynatrace
    namespace: keptn
  query: builtin:service.requestCount.total:merge(0):count?scope=tag(keptn_project:$PROJECT),tag(keptn_stage:$STAGE),tag(keptn_service:$SERVICE),tag(keptn_deployment:$DEPLOYMENT)
---
kind: AnalysisValueTemplate
apiVersion: metrics.keptn.sh/v1alpha3
metadata:
  name: response_time_p95
  creationTimestamp: null
spec:
  provider:
    name: dynatrace
    namespace: keptn
  query: builtin:service.response.time:merge(0):percentile(95)?scope=tag(keptn_project:$PROJECT),tag(keptn_stage:$STAGE),tag(keptn_service:$SERVICE),tag(keptn_deployment:$DEPLOYMENT)
`

func TestConvertSLI(t *testing.T) {
	// no provider nor namespace
	res, err := convertSLI([]byte(sliContent), "", "")
	require.NotNil(t, err)
	require.Equal(t, "", res)

	// invalid file content
	res, err = convertSLI([]byte("invalid"), "dynatrace", "keptn")
	require.NotNil(t, err)
	require.Equal(t, "", res)

	// happy path
	res, err = convertSLI([]byte(sliContent), "dynatrace", "keptn")
	require.Nil(t, err)
	require.Equal(t, expectedOutput, res)
}
