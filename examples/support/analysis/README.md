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
  name: response-time-p95
  namespace: analysis-demo
spec:
  provider:
    name: my-provider
  query: 'histogram_quantile(0.95, sum by(le) (rate(http_server_request_latency_seconds_bucket{job="{{.workload}}"}[1m])))'
EOF
```

Our second metric will be the error rate, which we will retrieve with the following template:

```shell
cat <<EOF | kubectl apply -f - 
apiVersion: metrics.keptn.sh/v1alpha3
kind: AnalysisValueTemplate
metadata:
  name: error-rate
  namespace: analysis-demo
spec:
  provider:
    name: my-provider
  query: 'rate(http_requests_total{status_code="500", job="{{.workload}}"}[1m]) or on() vector(0)'
EOF
```

### Setting up Prometheus

**Option 1: Using real data from podtato-head**

The queries provided to these `AnalysisValueTemplates` rely on actual data coming from Prometheus, particularly
the data for the **podtato-head** application that has been deployed in the [sample-app example](../../sample-app/README.md).
Also, this example assumes you have prometheus installed in your cluster, which can be done by following the instructions in
the [observability-example](../observability/README.md).

When you have done so, you also need to make sure that the podtato-head application is monitored by Prometheus by
creating a Prometheus `ServiceMonitor`. Ths can be done by applying the manifest in `./config/service-monitor.yaml`:

```shell
kubectl apply -f ./config/service-monitor.yaml
```

**Option 2: Using Mockserver**

If you do not want to go through the process of deploying an actual application and setting up Prometheus, you can also use
[MockServer](https://www.mock-server.com) to fake the monitoring data. To deploy the mock server, use the following command:

```shell
kubectl apply -f ./config/mock-server.yaml
```

Once the MockServer is up and running, we will adjust the `KeptnMetricsProvider` to retrieve data from there,
rather than from a real Prometheus instance:

```shell
cat <<EOF | kubectl apply -f - 
apiVersion: metrics.keptn.sh/v1alpha3
kind: KeptnMetricsProvider
metadata:
  name: my-provider
  namespace: analysis-demo
spec:
  type: prometheus
  targetServer: "http://mockserver.analysis-demo.svc.cluster.local:1080"
EOF
```

## Defining goals for the metrics 

Now that we have defined our metrics, it is time to describe what we expect from these values. This is done in an
`AnalysisDefinition`, which looks can be applied using the following command:

```shell
cat <<EOF | kubectl apply -f - 
apiVersion: metrics.keptn.sh/v1alpha3
kind: AnalysisDefinition
metadata:
  name: my-analysis-definition
  namespace: analysis-demo
spec:
  objectives:
    - analysisValueTemplateRef:
        name: response-time-p95
      target:
        failure:
          greaterThan:
            fixedValue: 500m
        warning:
          greaterThan:
            fixedValue: 300m
      weight: 1
      keyObjective: false
    - analysisValueTemplateRef:
        name: error-rate
      target:
        failure:
          greaterThan:
            fixedValue: 0
      weight: 1
      keyObjective: true
  totalScore:
    passPercentage: 60
    warningPercentage: 50
EOF
```

## Executing an Analysis

Finally, now that we have all of our metrics and our goals defined, we can apply an instance of the `Analysis` CRD
to perform an analysis for a specific timeframe

```shell
cat <<EOF | kubectl apply -f - 
apiVersion: metrics.keptn.sh/v1alpha3
kind: Analysis
metadata:
  name: analysis-sample
  namespace: analysis-demo
spec:
  timeframe:
    from: 2023-09-25T13:10:00Z
    to: 2023-09-25T13:15:00Z
  args:
    "workload": "podtato-head-frontend"
  analysisDefinition:
    name: my-analysis-definition
EOF
```

Once applied, the status of the analysis can be checked with:

```shell
kubectl get analyses -n analysis-demo
```

This command should yield a list of all analyses within our `analysis-demo` namespace, together with the current status of the
analysis. In our case, we will receive one analysis which has already completed and has passed:

**NOTE:** This example uses real data from an actual service monitored by Prometheus, so it could very well be that
the result of your analysis might be a different one, as those values heavily depend on the environment this example
is executed in.

```shell
NAME              ANALYSISDEFINITION       STATE       WARNING   PASS
analysis-sample   my-analysis-definition   Completed             true
```

To get further details on the analysis, we can also retrieve the complete yaml representation:

```shell
kubectl get analyses -n analysis-demo analysis-sample -oyaml 
```

This will return something like the following, and includes the `status.raw` field, which contains a detailed
breakdown on the result of each objective evaluation.

```yaml
apiVersion: metrics.keptn.sh/v1alpha3
kind: Analysis
metadata:
  name: analysis-sample
  namespace: analysis-demo
spec:
  analysisDefinition:
    name: my-analysis-definition
  args:
    workload: podtato-head-frontend
  timeframe:
    from: "2023-09-25T13:10:00Z"
    to: "2023-09-25T13:15:00Z"
status:
  pass: true
  raw: '{"objectiveResults":[{"result":{"failResult":{"operator":{"greaterThan":{"fixedValue":"500m"}},"fulfilled":false},"warnResult":{"operator":{"greaterThan":{"fixedValue":"300m"}},"fulfilled":false},"warning":false,"pass":true},"objective":{"analysisValueTemplateRef":{"name":"response-time-p95"},"target":{"failure":{"greaterThan":{"fixedValue":"500m"}},"warning":{"greaterThan":{"fixedValue":"300m"}}},"weight":1},"value":0.00475,"score":1},{"result":{"failResult":{"operator":{"greaterThan":{"fixedValue":"0"}},"fulfilled":false},"warnResult":{"operator":{},"fulfilled":false},"warning":false,"pass":true},"objective":{"analysisValueTemplateRef":{"name":"error-rate"},"target":{"failure":{"greaterThan":{"fixedValue":"0"}}},"weight":1,"keyObjective":true},"value":0,"score":1}],"totalScore":2,"maximumScore":2,"pass":true,"warning":false}'
  state: Completed
```
