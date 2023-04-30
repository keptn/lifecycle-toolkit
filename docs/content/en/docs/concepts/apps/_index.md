---
title: Apps
description: Learn what Keptn Apps are and how to use them
icon: concepts
layout: quickstart
weight: 10
hidechildren: true # this flag hides all sub-pages in the sidebar-multicard.html
---

An App contains information about all workloads and checks associated with an application.
It will use the following structure for the specification of the pre/post deployment and pre/post evaluations checks
that should be executed at app level:

```yaml
apiVersion: lifecycle.keptn.sh/v1alpha2
kind: KeptnApp
metadata:
  name: podtato-head
  namespace: podtato-kubectl
spec:
  version: "1.3"
  workloads:
  - name: podtato-head-left-arm
    version: 0.1.0
  - name: podtato-head-left-leg
    version: 1.2.3
  # The following are optional
  preDeploymentTasks:
  - pre-deployment-hello
  preDeploymentEvaluations:    
  - my-pre-deploy-evaluation
  postDeploymentTasks:
  - post-deployment-hello
  postDeploymentEvaluations:
  - my-post-deploy-evaluation
```

While changes in the workload version will affect only workload checks, a change in the app version will also cause a
new execution of app level checks.
