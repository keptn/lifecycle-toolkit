---
title: KeptnConfig
description: Learn what Keptn Configs are and how to use them
icon: concepts
layout: quickstart
weight: 10
hidechildren: true # this flag hides all sub-pages in the sidebar-multicard.html
---


### Keptn Config

A `KeptnConfig` CRD defines configuration values for the Keptn Lifecycle Toolkit.
Currently, it can be used to configure the URL of the OpenTelemetry collector.

A `KeptnConfig` looks like the following:

```yaml
apiVersion: options.keptn.sh/v1alpha1
kind: KeptnConfig
metadata:
  name: keptnconfig-sample
spec:
  OTelCollectorUrl: 'otel-collector:4317'
```
