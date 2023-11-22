---
title: Uninstall
description: How to uninstall Keptn
weight: 10
---

If you installed the previous version of Keptn using `helm`,
you can uninstall it together with all CRDs, webhooks and
custom resources with using the following command:

```shell
helm uninstall keptn -n keptn-lifecycle-toolkit-system
```

If your Keptn instance is not installed in the
`keptn-lifecycle-toolkit-system` namespace, please substitute
it with your custom one.

> **Warning**
Please be aware that uninstalling Keptn from your cluster
will cause loss of all your Keptn data.
