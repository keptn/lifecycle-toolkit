---
title: Evaluations
description: Learn what Keptn Evaluations are and how to use them
icon: concepts
layout: quickstart
weight: 10
hidechildren: true # this flag hides all sub-pages in the sidebar-multicard.html
---


### Keptn Evaluation Definition

A `KeptnEvaluationDefinition` is a CRD used to define evaluation tasks that can be run by the Keptn Lifecycle Toolkit
as part of pre- and post-analysis phases of a workload or application.

A Keptn evaluation definition looks like the following:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha3
kind: KeptnEvaluationDefinition
metadata:
  name: my-prometheus-evaluation
  namespace: example
spec:
  source: prometheus
  objectives:
    - keptnMetricRef:
        name: available-cpus
        namespace: example
      evaluationTarget: ">1"
    - keptnMetricRef:
        name: cpus-throttling
        namespace: example
      evaluationTarget: "<0.01"
```

A `KeptnEvaluationDefinition` references one or more [`KeptnMetric`s](../metrics/).
If multiple `KeptnMetric`s are used, the Keptn Lifecycle Toolkit will consider the
evaluation successful if **all** metrics are respecting their `evaluationTarget`.
