---
title: KeptnConfigs
description: Learn what Keptn Configs are and how to use them
icon: concepts
layout: quickstart
weight: 10
hidechildren: true # this flag hides all sub-pages in the sidebar-multicard.html
---


### Keptn Config
A `KeptnConfig` is a CRD used to define configuration values of the Keptn Lifecycle Toolkit.
In the current state, there is a possibility to configure url of OTel collector.

A `KeptnConfig` looks like the following:

```yaml
apiVersion: options.keptn.sh/v1alpha1
kind: KeptnConfig
metadata:
  name: keptnconfig-sample
spec:
  OTelCollectorUrl: 'otel-collector:4317'
```
