# AutoScaling with HPA and KeptnMetrics

This example shows how `KeptnMetrics` can be used as a reference for a `HorizontalPodAutoscaler` to decide when to
scale workloads up or down.
To demonstrate this, the following steps will be covered in this example:

1. Deploy a sample application
2. Create a `KeptnMetric` to monitor the throttled CPU of all pods serving our application
3. Generate some load to put the application under stress
4. Deploy a `HorizontalPodAutoscaler` to scale up/down to meet our goal for the CPU-throttling metric

## Deploying the Application

To deploy the application, you can run the following command:

```shell
make 1-deploy-app
```

This will create a namespace called `podtato-metrics`, and a deployment called `podtato-head-entry` that is
accessible via a `ClusterIP` service.
After executing the command, you should see one
pod running the application we just created:

```shell
$kubectl get pods -n podtato-metrics
NAME                                  READY   STATUS    RESTARTS   AGE
podtato-head-entry-795b4bf76c-bjlfw   1/1     Running   0          1m
```

## Creating the Metric

To create the metric that keeps track of the throttled CPU of our application, run the following command:

```shell
make 2-create-metric
```

This will create a `KeptnMetricsProvider` pointing to Prometheus
(`http://prometheus-k8s.monitoring.svc.cluster.local:9090`), and a `KeptnMetric` that retrieves
the value of the following query:

```shell
avg(rate(container_cpu_cfs_throttled_seconds_total{container="server", namespace="podtato-metrics"}[1m]))
```

To verify that the metric is wired up properly, we can retrieve it via `kubectl`:

```shell
$kubectl get keptnmetrics.metrics.keptn.sh -n podtato-metrics cpu-throttling
NAME             PROVIDER     QUERY                                                                                                       VALUE
cpu-throttling   prometheus   avg(rate(container_cpu_cfs_throttled_seconds_total{container="server", namespace="podtato-metrics"}[1m]))   0.01433336027598159
```

## Generating Load

Now that we have our application up and running, and can retrieve the `KeptnMetric` value,
it is time to generate some load.
To do so, we will create a `Job` that regularly
sends a request to our application.
The Job can be created using the following command:

```shell
make 3-generate-load
```

Once the Job is running, we will see that our `KeptnMetric`'s value will increase after some time:

```shell
$kubectl get keptnmetrics.metrics.keptn.sh -n podtato-metrics cpu-throttling
NAME             PROVIDER     QUERY                                                                                                       VALUE
cpu-throttling   prometheus   avg(rate(container_cpu_cfs_throttled_seconds_total{container="server", namespace="podtato-metrics"}[1m]))   0.25475392739204
```

## Deploying the HorizontalPodAutoscaler

Now, to meet the demand of our application, we will deploy a `HorizontalPodAutoscaler` that will
observe the value of the `cpu-throttling` metric, and scale the demo application up or down, based on the target
we have specified for our metric.
In our case, we want to ensure that the value of our metric stays
below `0.05`, and we are willing to scale up to 10 replicas of our demo application:

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: podtato-metrics-hpa
  namespace: podtato-metrics
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: podtato-head-entry
  minReplicas: 1
  maxReplicas: 10
  metrics:
    - type: Object
      object:
        metric:
          name: cpu-throttling
        describedObject:
          apiVersion: metrics.keptn.sh/v1alpha2
          kind: KeptnMetric
          name: cpu-throttling
        target:
          type: Value
          value: "0.05"
```

The HPA can be deployed using the following command:

```shell
make 4-deploy-hpa
```

Once the HPA has been deployed, we should immediately see it scaling up the replica count of our application,
since the metric value was already above our target.
You can also verify this by inspecting the current
state of the `podtato-metrics-hpa` autoscaler:

```shell
$ make get-hpa-status
kubectl describe  horizontalpodautoscalers.autoscaling -n podtato-metrics podtato-metrics-hpa
Warning: autoscaling/v2beta2 HorizontalPodAutoscaler is deprecated in v1.23+, unavailable in v1.26+; use autoscaling/v2 HorizontalPodAutoscaler
Name:                                                             podtato-metrics-hpa
Namespace:                                                        podtato-metrics
Labels:                                                           <none>
Annotations:                                                      <none>
CreationTimestamp:                                                Mon, 03 Apr 2023 11:18:58 +0200
Reference:                                                        Deployment/podtato-head-entry
Metrics:                                                          ( current / target )
  "cpu-throttling" on KeptnMetric/cpu-throttling (target value):  57m / 50m
Min replicas:                                                     1
Max replicas:                                                     10
Deployment pods:                                                  6 current / 7 desired
Conditions:
  Type            Status  Reason              Message
  ----            ------  ------              -------
  AbleToScale     True    SucceededRescale    the HPA controller was able to update the target scale to 7
  ScalingActive   True    ValidMetricFound    the HPA was able to successfully calculate a replica count from KeptnMetric metric cpu-throttling
  ScalingLimited  False   DesiredWithinRange  the desired count is within the acceptable range
Events:
  Type     Reason                        Age                 From                       Message
  ----     ------                        ----                ----                       -------
  Normal   SuccessfulRescale             70m                 horizontal-pod-autoscaler  New size: 9; reason: KeptnMetric metric cpu-throttling above target
```

As a consequence of that, you should eventually see a decrease of the `cpu-throttling` value:

```shell
$kubectl get keptnmetrics.metrics.keptn.sh -n podtato-metrics cpu-throttling
NAME             PROVIDER     QUERY                                                                                                       VALUE
cpu-throttling   prometheus   avg(rate(container_cpu_cfs_throttled_seconds_total{container="server", namespace="podtato-metrics"}[1m]))   0.036489273639926
```
