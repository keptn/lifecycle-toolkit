---
comments: true
---

# Scaling Workloads with KEDA based on Keptn metrics

If you want to use [KEDA](https://keda.sh/) to scale your workloads, you can set it up
to consume metrics from Keptn.
This gives you great flexibility for your setup since Keptn can
consolidate different observability providers for you and KEDA
will simplify the scaling operations for you.

This use case is enabled by the Keptn Metrics
[custom API](https://kubernetes.io/docs/reference/external-api/custom-metrics.v1beta2/)
that the Keptn Metrics Operator provides.

## Before you begin

For this use case presentation, a few components need to be deployed to the
cluster in order to have a full setup:

- Keptn Metrics Operator: Will be used for metrics consolidation.

    For more information about installation please refer to the
    [installation guide](../installation/index.md).

    > The Keptn Lifecycle Operator does not need to be installed for this use-case.

- [KEDA](https://keda.sh/): Will be used for scaling.
- [Prometheus](https://prometheus.io/): Will be used as metrics provider.
  
    For more information about how to install Prometheus into your cluster, please
    refer to the [Prometheus documentation](https://prometheus.io/docs/prometheus/latest/installation/).


## Deploy sample application

First, we need to deploy our application to the cluster.
For this we are going to
use a single service `podtato-head` application.

=== "deployment.yaml"

    ```yaml
    {% include "./assets/keda/sample-app.yaml" %}
    ```

=== "service.yaml"

    ```yaml
    {% include "./assets/keda/sample-service.yaml" %}
    ```

Please create a `podtato-kubectl` namespace and apply the above manifests
to your cluster and continue with the next steps.
After applying, please make sure that the application is up and running:

```shell
$ kubectl get pods -n podtato-kubectl
podtato-head-entry-58d6485d9b-ld9x2         1/1     Running     (2m ago)
```

## Create KeptnMetric and KeptnMetricsProvider custom resources

To be able to react on the metrics of our application, we need to create
`KeptnMetrics` and `KeptnMetricsProvider` custom resources.
These metrics are
exposed via the Keptn Metrics Operator, which gives us the possibility to configure
KEDA to react on the values of these metrics:

=== "KeptnMetric"

    ```yaml
    {% include "./assets/keda/keptnmetric.yaml" %}
    ```

=== "KeptnMetricsProvider"

    ```yaml
    {% include "./assets/keda/keptnmetricsprovider.yaml" %}
    ```

For more information about the `KeptnMetric` and `KeptnMetricsProvider` custom resources,
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

## Set up the KEDA ScaledObject

Now that we are able to retrieve the value of our metric, and have it stored in
our cluster in the status of our `KeptnMetric` custom resource, we can configure
a `HorizontalPodAutoscaler` to make use of this information and therefore scale
our application automatically:

```yaml
{% include "./assets/keda/scaledobject.yaml" %}
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
