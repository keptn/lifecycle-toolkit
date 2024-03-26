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

<!-- markdownlint-disable MD046 -->

=== "deployment.yaml"

    ```yaml
    {% include "./assets/keda/sample-app.yaml" %}
    ```

=== "service.yaml"

    ```yaml
    {% include "./assets/keda/sample-service.yaml" %}
    ```

<!-- markdownlint-enable MD046 -->

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

<!-- markdownlint-disable MD046 -->

=== "KeptnMetric"

    ```yaml
    {% include "./assets/keda/keptnmetric.yaml" %}
    ```

=== "KeptnMetricsProvider"

    ```yaml
    {% include "./assets/keda/keptnmetricsprovider.yaml" %}
    ```

<!-- markdownlint-enable MD046 -->

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
  Value:         4.53
```

Here we can see that the value of the `cpu-throttling` metric is `4.53`

## Set up the KEDA ScaledObject

Now that we are able to retrieve the value of our metric, and have it stored in
our cluster in the status of our `KeptnMetric` custom resource, we can configure
a `ScaledObject` to make use of this information in KEDA and therefore scale
our application automatically:

```yaml
{% include "./assets/keda/scaledobject.yaml" %}
```

As we can see in this example, by setting the `url` field,
we are now referring to the `KeptnMetric` and fetch it from the
Keptn Metrics Operator.
KEDA will scale up our application, until our target value or `1` is reached,
or we hit the maximum number of replicas which is `3` in this example.


If the load of the application is high enough, we will be able to see
the automatic scaling of our application:

```shell
$ kubectl describe scaledobject -n podtato-kubectl my-scaledobject
Name:         my-scaledobject
Namespace:    podtato-kubectl
Labels:       deploymentName=podtato-head-entry
              scaledobject.keda.sh/name=my-scaledobject
API Version:  keda.sh/v1alpha1
Kind:         ScaledObject
Spec:
  Max Replica Count:  3
  Scale Target Ref:
    Name:  podtato-head-entry
  Triggers:
    Metadata:
      Target Value:    1
      Unsafe Ssl:      true
      URL:             http://metrics-operator-service.keptn-system.svc.cluster.local:9999/api/v1/metrics/chainsaw-proud-wallaby/test
      Value Location:  value
    Type:              metrics-api
Status:
  Conditions:
    Message:  ScaledObject is defined correctly and is ready for scaling
    Reason:   ScaledObjectReady
    Status:   True
    Type:     Ready
    Message:  Scaling is not performed because triggers are not active
    Reason:   ScalerNotActive
    Status:   False
    Type:     Active
    Message:  No fallbacks are active on this scaled object
    Reason:   NoFallbackFound
    Status:   False
    Type:     Fallback
    Status:   Unknown
    Type:     Paused
  External Metric Names:
    s0-metric-api-value
  Health:
    s0-metric-api-value:
      Number Of Failures:  0
      Status:              Happy
  Hpa Name:                keda-hpa-test-scaledobject
  Last Active Time:        2024-03-26T09:36:36Z
  Original Replica Count:  1
  Scale Target GVKR:
    Group:            apps
    Kind:             Deployment
    Resource:         deployments
    Version:          v1
  Scale Target Kind:  apps/v1.Deployment
Events:
  Type     Reason              Age                From           Message
  ----     ------              ----               ----           -------
  Normal   KEDAScalersStarted  63s                keda-operator  Started scalers watch
  Normal   ScaledObjectReady   63s                keda-operator  ScaledObject is ready for scaling
  Warning  KEDAScalerFailed    33s (x2 over 63s)  keda-operator  error requesting metrics endpoint: valueLocation must point to value of type number or a string representing a Quantity got: ''
  Normal   KEDAScalersStarted  18s (x5 over 63s)  keda-operator  Scaler metrics-api is built.
```

If we retrieve the pods of our application, we can see that, instead of
a single instance at the beginning, there are currently 3 instances running:

```shell
$ kubectl get pods -n podtato-kubectl
NAME                                  READY   STATUS    RESTARTS   AGE
podtato-head-entry-7796c8f786-4cdtc   1/1     Running   0          3m29s
podtato-head-entry-7796c8f786-nsk2c   1/1     Running   0          2m41s
podtato-head-entry-7796c8f786-qj85h   1/1     Running   0          2m41s
```
