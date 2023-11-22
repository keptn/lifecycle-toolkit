---
title: Keptn + HorizontalPodAutoscaler
description: Using the HorizontalPodAutoscaler
weight: 20
---

## Scaling Workloads based on Keptn metrics

Kubernetes provides us with a lot of built-in capabilities to ensure
running enough replicas in order to meet the current demand of the Workloads.
One of these is the
[HorizontalPodAutoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)
(HPA).

HPA can make use of the [Keptn Metrics](../getting-started/metrics.md) custom API
in order to scale the number of replicas of the Workloads based on the current
load by using metrics such as CPU throttling, memory consumption or response time.

### Installation of Keptn Metrics Operator

To use HPA with custom metrics API, Keptn Metrics Operator needs to be installed on our cluster.
For more information about installation please refer to the official
[installation guide](../installation/_index.md).

> **Note**
Please be aware that Keptn Lifecycle Operator does not need to be installed for this use-case.

### Installation of metrics provider (optional)

If you do not have metrics provider installed on your cluster yet, please do so.

For this tutorial we are going to use [Prometheus](https://prometheus.io/).
For more information about how to install Prometheus into your cluster, please
refer to the [official Prometheus documentation](https://prometheus.io/docs/prometheus/latest/installation/).

### Deploy sample application

In the next step, we need to deploy our application to the cluster.
For this we are going to
use a single service `podtato-head` application.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: podtato-head-entry
  namespace: podtato-kubectl
  labels:
    app: podtato-head
spec:
  selector:
    matchLabels:
      component: podtato-head-entry
  template:
    metadata:
      labels:
        component: podtato-head-entry
    spec:
      terminationGracePeriodSeconds: 5
      containers:
      - name: server
        image: ghcr.io/podtato-head/entry:0.2.7
        imagePullPolicy: Always
        resources:
          limits:
            cpu: 5m
            memory: 128Mi
          requests:
            cpu: 1m
            memory: 64Mi
        ports:
        - containerPort: 9000
        env:
        - name: PODTATO_PORT
          value: "9000"
---
apiVersion: v1
kind: Service
metadata:
  name: podtato-head-entry
  namespace: podtato-kubectl
  labels:
    app: podtato-head
spec:
  selector:
    component: podtato-head-entry
  ports:
  - name: http
    port: 9000
    protocol: TCP
    nodePort: 30900
    targetPort: 9000
  type: NodePort
```

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
exposed via custom metrics API, what gives us the possibility to configure HPA
to react on the values of these metrics:

```yaml
apiVersion: metrics.keptn.sh/v1beta1
kind: KeptnMetricsProvider
metadata:
  name: prometheus-provider
  namespace: podtato-kubectl
spec:
  type: prometheus
  targetServer: <your-metrics-provider-server>
---
apiVersion: metrics.keptn.sh/v1beta1
kind: KeptnMetric
metadata:
  name: cpu-throttling
spec:
  provider:
    name: prometheus-provider
  query: 'avg(rate(container_cpu_cfs_throttled_seconds_total{container="server", namespace="podtato-kubectl"}))'
  fetchIntervalSeconds: 10
  range:
    interval: "30s"
```

For more information about the `KeptnMetric` and `KeptnMetricsProvider` custom resources,
please refer to the official [CRD documentation](../crd-ref/metrics/v1beta1/_index.md).

After a few seconds we should be able to see values for the `cpu-throttling` metrics:

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

### Setup HorizontalPodAutoscaler

And now that we are able to retrieve the value of our metric, and have it stored in
our cluster in the status of our `KeptnMetric` custom resource, we can configure
a `HorizontalPodAutoscaler` to make use of this information and therefore scale
our application automatically:

```yaml
apiVersion: autoscaling/v
kind: HorizontalPodAutoscaler
metadata:
  name: podtato-hpa
  namespace: podtato-kubectl
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
          apiVersion: metrics.keptn.sh/v1beta1
          kind: KeptnMetric
          name: cpu-throttling
        target:
          type: Value
          value: "5"
```

As we can see in this example, we are now referring to the `KeptnMetric`
we applied earlier, and tell `HPA` to scale up our application, until our
target value of `5` for this metric is reached, or the number of replicas
has reached a maximum of `10`.

Now if the load of the application is high enough, we will be able to see
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

If we retrieve the pods of our application, we can see that instead of
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
