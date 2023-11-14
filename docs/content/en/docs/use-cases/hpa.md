---
title: Keptn + HorizontalPodAutoscaler
description: Using the HorizontalPodAutoscaler
weight: 20
---

Use the Kubernetes Custom Metrics API
to refer to `KeptnMetric` via the
[Kubernetes HorizontalPodAutoscaler](https://kubernetes.io/docs/tasks/run-application/horizontal-pod-autoscale/)
(HPA),
as in the following example:

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: podtato-head-entry
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
          name: keptnmetric-sample
        describedObject:
          apiVersion: metrics.keptn.sh/v1beta1
          kind: KeptnMetric
          name: keptnmetric-sample
        target:
          type: Value
          value: "10"
```

See the
[Scaling Kubernetes Workloads based on Dynatrace Metrics](https://www.linkedin.com/pulse/scaling-kubernetes-workloads-based-dynatrace-metrics-keptnproject/)
blog post
for a detailed discussion of doing this with Dynatrace metrics.
The same approach could be used to implement HPA with other data providers.
