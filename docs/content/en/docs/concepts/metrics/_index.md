---
title: Metrics
description: Learn what Keptn Metrics are and how to use them
icon: concepts
layout: quickstart
weight: 10
hidechildren: true # this flag hides all sub-pages in the sidebar-multicard.html
---

### Keptn Metric
A `KeptnMetric` is a CRD representing a metric. The metric will be collected from the provider specified in the specs.provider.name field. The query is a string in the provider-specific query language, used to obtain a metric. Providing the metrics as CRD into a K8s cluster will facilitate the reusability of this data across multiple components. Furthermore, this allows using multiple observability platforms for different metrics. Please note, there is a limitation that `KeptnMetric` resource needs to be created only in `keptn-lifecycle-toolkit-system` namespace.

A `KeptnMetric` looks like the following:

```yaml
apiVersion: metrics.keptn.sh/v1alpha1
kind: KeptnMetric
metadata:
  name: keptnmetric-sample
  namespace: keptn-lifecycle-toolkit-system
spec:
  provider:
    name: "<your-keptn-evaluation-provider-crd-name>"
  query: "sum(kube_pod_container_resource_limits{resource='cpu'})"
  fetchIntervalSeconds: 5
```

Keptn metrics can be exposed as OTel metrics via port `9999` of the KLT operator. To expose them, the env variable `EXPOSE_KEPTN_METRICS` in the operator manifest needs to be set to `true`. The default value of this variable is `true`. To access the metrics, use the following command:

```
kubectl port-forward deployment/klc-controller-manager 9999 -n keptn-lifecycle-toolkit-system
```

and access the metrics via your browser with:

```
http://localhost:9999/metrics
```