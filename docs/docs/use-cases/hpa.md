---
comments: true
---

# Scaling Workloads with HPA based on Keptn metrics

Kubernetes provides many built-in capabilities to ensure
that enough replicas are running in order to meet the current demand of
your [workloads](https://kubernetes.io/docs/concepts/workloads/).
One of these is the
[HorizontalPodAutoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)
(HPA).

An HPA can make use of the Keptn Metrics
[custom API](https://kubernetes.io/docs/reference/external-api/custom-metrics.v1beta2/)
to scale the number of replicas of [workloads](https://kubernetes.io/docs/concepts/workloads/) based on the current
load.
It does this by using metrics such as CPU throttling, memory consumption, or response time.

## Installation of Keptn Metrics Operator

To use an HPA with the custom metrics API, the
Keptn Metrics Operator must be installed on the cluster.
For more information about installation please refer to the
[installation guide](../installation/index.md).

> **Note**
  The Keptn Lifecycle Operator does not need to be installed for this use-case.

## Installation of metrics provider (optional)

If you do not have a metrics provider installed on your cluster yet, please do so.

For this tutorial we are going to use [Prometheus](https://prometheus.io/).
For more information about how to install Prometheus into your cluster, please
refer to the [Prometheus documentation](https://prometheus.io/docs/prometheus/latest/installation/).

## Deploy sample application

First, we need to deploy our application to the cluster.
For this we are going to
use a single service `podtato-head` application.

```yaml
{% include "./assets/hpa/sample-app.yaml" %}
```

Please create a `podtato-kubectl` namespace and apply the above manifest
to your cluster and continue with the next steps.
After applying, please make sure that the application is up and running:

```shell
$ kubectl get pods -n podtato-kubectl
podtato-head-entry-58d6485d9b-ld9x2         1/1     Running     (2m ago)
```

## Create KeptnMetric and KeptnMetricsProvider resources

To be able to react on the metrics of our application, we need to create
[KeptnMetrics](../reference/crd-reference/metric.md) and
[KeptnMetricsProvider](../reference/crd-reference/metricsprovider.md) resources.
These metrics are
exposed via the custom metrics API, which gives us the possibility to configure
the HPA to react on the values of these metrics:

```yaml
{% include "./assets/hpa/keptnmetric.yaml" %}
```

For more information about the `KeptnMetric` and `KeptnMetricsProvider` resources,
please refer to the [CRD documentation](../reference/api-reference/metrics/v1/index.md).

After a few seconds we should be able to see values for the `cpu-throttling` metric:

```shell
$ kubectl describe  keptnmetrics.metrics.keptn.sh cpu-throttling -n podtato-kubectl
Name:         cpu-throttling
Namespace:    podtato-kubectl
API Version:  metrics.keptn.sh/v1
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

## Set up the HorizontalPodAutoscaler

Now that we are able to retrieve the value of our metric, and have it stored in
our cluster in the status of our `KeptnMetric` resource, we can configure
a `HorizontalPodAutoscaler` to make use of this information and therefore scale
our application automatically:

```yaml
{% include "./assets/hpa/hpa.yaml" %}
```

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
  ----            ------  ------               -------
  AbleToScale     True    ScaleDownStabilized  recent recommendations were higher than current one, applying the highest recent recommendation
  ScalingActive   True    ValidMetricFound     the HPA was able to successfully calculate a replica count from KeptnMetric metric cpu-throttling
  ScalingLimited  True    TooManyReplicas      the desired replica count is more than the maximum replica count
Events:
  Type    Reason             Age                  From                       Message
  ----    ------             ----                 ----                       -------
  Normal  SuccessfulRescale  7m18s (x5 over 16h)  horizontal-pod-autoscaler  New size: 4; reason: KeptnMetric metric cpu-throttling above target
  Normal  SuccessfulRescale  6m18s                horizontal-pod-autoscaler  New size: 7; reason: KeptnMetric metric cpu-throttling above target
  Normal  SuccessfulRescale  6m3s (x4 over 16h)   horizontal-pod-autoscaler  New size: 10; reason: KeptnMetric metric cpu-throttling above target
```

If we retrieve the pods of our application, we can see that, instead of
a single instance at the beginning, there are currently 10 instances running:

```shell
$ kubectl get pods -n podtato-kubectl
NAME                                      READY   STATUS    RESTARTS   AGE
podtato-head-entry-795b4bf76c-22vl8       1/1     Running   0          4m50s
podtato-head-entry-795b4bf76c-4mqz5       1/1     Running   0          4m50s
podtato-head-entry-795b4bf76c-g5bcr       1/1     Running   0          5m5s
podtato-head-entry-795b4bf76c-h22pq       1/1     Running   0          6m5s
podtato-head-entry-795b4bf76c-kgcgb       1/1     Running   0          4m50s
podtato-head-entry-795b4bf76c-kkt82       1/1     Running   0          5m5s
podtato-head-entry-795b4bf76c-lmfnx       1/1     Running   0          6m5s
podtato-head-entry-795b4bf76c-pnq2f       1/1     Running   0          15m
podtato-head-entry-795b4bf76c-r5dx4       1/1     Running   0          5m5s
podtato-head-entry-795b4bf76c-vwdj7       1/1     Running   0          6m5s
```
