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
apiVersion: lifecycle.keptn.sh/v1alpha2
kind: KeptnEvaluationDefinition
metadata:
  name: my-prometheus-evaluation
spec:
  source: prometheus
  objectives:
    - name: query-1
      query: "xxxx"
      evaluationTarget: <20
    - name: query-2
      query: "yyyy"
      evaluationTarget: >4
```

### Keptn Evaluation Provider

A `KeptnEvaluationProvider` is a CRD used to define evaluation provider, which will provide data for the
pre- and post-analysis phases of a workload or application.

A Keptn evaluation provider looks like the following:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha2
kind: KeptnEvaluationProvider
metadata:
  name: prometheus
spec:
  targetServer: "http://prometheus-k8s.monitoring.svc.cluster.local:9090"
  secretName: prometheusLoginCredentials
```
