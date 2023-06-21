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
Additionally, it can be used to set the time interval in which automatic app discovery
searches for workloads to put into the same auto-generated `KeptnApp`.
When the parameter is not set, the default value is 30 seconds.

A `KeptnConfig` looks like the following:

```yaml
apiVersion: options.keptn.sh/v1alpha1
kind: KeptnConfig
metadata:
  name: keptnconfig-sample
spec:
  OTelCollectorUrl: 'otel-collector:4317'
  keptnAppCreationRequestTimeoutSeconds: 30
```
