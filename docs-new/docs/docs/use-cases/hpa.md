# Keptn + HorizontalPodAutoscaler

## Scaling Workloads based on Keptn metrics

Kubernetes provides many built-in capabilities to ensure
that enough replicas are running in order to meet the current demand of your [workloads](https://kubernetes.io/docs/concepts/workloads/).
One of these is the
[HorizontalPodAutoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)
(HPA).

An HPA can make use of the Keptn Metrics
[custom API](https://kubernetes.io/docs/reference/external-api/custom-metrics.v1beta2/)
to scale the number of replicas of [workloads](https://kubernetes.io/docs/concepts/workloads/) based on the current
load.
It does this by using metrics such as CPU throttling, memory consumption, or response time.

### Installation of Keptn Metrics Operator

To use an HPA with the custom metrics API, the
Keptn Metrics Operator must be installed on the cluster.
For more information about installation please refer to the
[installation guide](../installation/_index.md).

> **Note**
The Keptn Lifecycle Operator does not need to be installed for this use-case.

### Installation of metrics provider (optional)

If you do not have a metrics provider installed on your cluster yet, please do so.

For this tutorial we are going to use [Prometheus](https://prometheus.io/).
For more information about how to install Prometheus into your cluster, please
refer to the [Prometheus documentation](https://prometheus.io/docs/prometheus/latest/installation/).

### Deploy sample application

First, we need to deploy our application to the cluster.
For this we are going to
use a single service `podtato-head` application.

{{< embed path="docs/assets/hpa/sample-app.yaml" >}}

Please create a `podtato-kubectl` namespace and apply the above manifest
to your cluster and continue with the next steps.
After applying, please make sure that the application is up and running:

```shell
$ kubectl get pods -n podtato-kubectl
podtato-head-entry-58d6485d9b-ld9x2         1/1     Running     (2m ago)
```

### Create KeptnMetric and KeptnMetricsProvider custom resources

To be able to react on the metrics of our application, we need to create
`KeptnMetrics` and `KeptnMetricsProvider` custom resources.
These metrics are
exposed via the custom metrics API, which gives us the possibility to configure
the HPA to react on the values of these metrics:

{{< embed path="docs/assets/hpa/keptnmetric.yaml" >}}

For more information about the `KeptnMetric` and `KeptnMetricsProvider` custom resources,
please refer to the [CRD documentation](../reference/api-reference/metrics/v1beta1/).

After a few seconds we should be able to see values for the `cpu-throttling` metric:

```shell
$ kubectl describe  keptnmetrics.metrics.keptn.sh cpu-throttling -n podtato-kubectl
Name:         cpu-throttling
Namespace:    podtato-kubectl
API Version:  metrics.keptn.sh/v1beta1
Kind:         KeptnMetric
Spec:
  Fetch Interval Seconds:  60
  Provider:
    Name:  prometheus-provider
  Query:  avg(rate(container_cpu_cfs_throttled_seconds_total{container="server", namespace="podtato-kubectl"}))
Status:
  Raw Value: <omitted for readability>
  Value:         1.63
```

Here we can see that the value of the `cpu-throttling` metric is `1.63`

### Set up the HorizontalPodAutoscaler

Now that we are able to retrieve the value of our metric, and have it stored in
our cluster in the status of our `KeptnMetric` custom resource, we can configure
a `HorizontalPodAutoscaler` to make use of this information and therefore scale
our application automatically:

{{< embed path="docs/assets/hpa/hpa.yaml" >}}

As we can see in this example, we are now referring to the `KeptnMetric`
we applied earlier, and tell the HPA to scale up our application, until our
target value of `5` for this metric is reached, or the number of replicas
has reached a maximum of `10`.

If the load of the application is high enough, we will be able to see
the automatic scaling of our application:

```shell
$ kubectl describe  horizontalpodautoscalers.autoscaling -n podtato-kubectl podtato-hpa
Name:                                                             podtato-hpa
Namespace:                                                        podtato-kubectl
Reference:                                                        Deployment/podtato-head-entry
Metrics:                                                          ( current / target )
  "cpu-throttling" on KeptnMetric/cpu-throttling (target value):  30.5 / 5
Min replicas:                                                     1
Max replicas:                                                     10
Deployment pods:                                                  10 current / 10 desired
Conditions:
  Type            Status  Reason               Message
  