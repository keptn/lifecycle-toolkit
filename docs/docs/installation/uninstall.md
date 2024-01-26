---
comments: true
---

# Uninstall

> **Warning**
Please be aware that uninstalling Keptn from your cluster
will cause loss of all your Keptn data.

If you installed the previous version of Keptn using `helm`,
you can uninstall it together with all CRDs, webhooks and
custom resources with using the following command:

```shell
helm uninstall keptn -n keptn-system
```

If your Keptn instance is not installed in the
`keptn-system` namespace, please substitute
it with your custom one.
