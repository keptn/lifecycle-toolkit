---
title: Metrics
description: Learn what Keptn Metrics are and how to use them
icon: concepts
layout: quickstart
weight: 10
hidechildren: true # this flag hides all sub-pages in the sidebar-multicard.html
---

### Keptn Metric
A `KeptnMetric` is a CRD used to define SLI provider with a query and to store metric data fetched from the provider. Providing the metrics as CRD into a K8s cluster will facilitate the reusability of this data across multiple components. Furthermore, this allows using multiple observability platforms for different metrics.

A `KeptnMetric` looks like the following:

```yaml
apiVersion: metrics.keptn.sh/v1alpha1
kind: KeptnMetric
metadata:
  name: keptnmetric-sample
  namespace: keptn-lifecycle-toolkit-system
spec:
  provider:
    name: "<your-keptn-evaluation-provider>"
  query: "<your query>"
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