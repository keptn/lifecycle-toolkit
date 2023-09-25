# Metric Analysis

This example shows how the `Analysis` feature of the `metrics-operator` can be used to 
define goals for metric values and evaluate them.

## Difference between `Analysis` and `Evaluations`

`Evaluations` provide an easy way of checking whether the current value a `KeptnMetric` fulfills
a requirement, i.e. is below or above a certain threshold. This is ideal for doing simple checks, such as
validating that a cluster currently has enough resources for the pod of a deployment to be scaled up.
In other cases however, it may not be enough to just check the current value of a metric (e.g. the number of 
available resources), but one may want to evaluate the value of a metric for a specific timeframe. For example,
one may execute load tests for a service after it has been deployed, and then verify if performance-related metrics,
such as response time or error rate during the time where the tests have been executed meet certain criteria.
That is where the `Analysis` capabilities of the `metrics-operator` come in. 

## Defining Metrics

The first step using this feature is to figure out what it is we want to analyze. In our example, we would like to
analyze the response time and the error rate for a service. Let's assume those metrics are retrieved from Prometheus,
so first we're going to create a `KeptnMetricsProvider`. Just like when using `KeptnMetrics`, the provider will
tell Keptn where to retrieve the values we are interested in from.

To create the provider, execute the following commands from this directory:

```shell
kubectl create namespace analysis-demo
```

```shell
cat <<EOF | kubectl apply -f - 
apiVersion: metrics.keptn.sh/v1alpha3
kind: KeptnMetricsProvider
metadata:
  name: my-provider
  namespace: analysis-demo
spec:
  type: prometheus
  targetServer: "http://prometheus-k8s.monitoring.svc.cluster.local:9090"
EOF
```

Next, we will define the queries for the metrics we are interested in. This is done by creating an
`AnalysisValueTemplate` for each of the metrics. Both of these `AnalysisValueTemplates` will refer to the provider
we just created and therefore retrieve the values from the Prometheus server the provider is pointing to. As the name
`AnalysisValueTemplate` suggests, this resource benefits from the support of the go templating syntax, with which we can 
include placeholders, e.g. for the name of the workload, in the query. This way, the same `AnalysisValueTemplate`
can be reused for different workloads, if the query for retrieving the value relevant for them only differs in e.g. their label
selectors. With that being said, let's apply the `AnalysisValueTemplates` to our cluster.
First, we are going to create the template for the response time:

```shell
cat <<EOF | kubectl apply -f - 
apiVersion: metrics.keptn.sh/v1alpha3
kind: AnalysisValueTemplate
metadata:
  name: response-time
spec:
  provider:
    name: my-mocked-provider
  query: 'histogram_quantile(0.95, sum by(le) (rate(http_request_duration_seconds_bucket{handler="{{.workload}}"})))'
EOF
```


