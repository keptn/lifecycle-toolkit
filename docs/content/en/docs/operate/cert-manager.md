---
title: Use Keptn with cert-manager.io (optional)
description: Replace the default KLT cert-manager
weight: 30
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

The Keptn Lifecycle Toolkit includes
a light-weight, customized cert-manager
that is used to register Webhooks to the [KubeAPI](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/).
Bundling the cert-manager simplifies the installation for new users
and provides the functionality KLT needs
without the overhead of other cert-managers.
For a description of the architecture, see
[Keptn Certificate Manager](../concepts/architecture/cert-manager.md).

KLT also works well with `cert-manager.io`.
If you are already using `cert-manager.io`,
you can continue to use it for other components
and use the KLT cert-manager just for KLT activities
or you can disable the KLT cert-manager
and configure KLT to use `cert-manager.io`.

If you want KLT to use `cert-manager.io`,
you must configure it *before* you install KLT.
The steps are:

* Install `cert-manager.io` if it is not already installed.
* Add the `Certificate` and `Issuer` CRs for `cert-manager.io`.
* (optional) Install Keptn without the built-in `klt-cert-manager` via Helm

## Add the CR(s) for cert-manager.io

These are the CRs for `cert-manager.io` to be applied to your cluster:

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: klt-certs
  namespace: <your-namespace>
spec:
  dnsNames:
  - lifecycle-webhook-service.<your-namespace>.svc
  - lifecycle-webhook-service.<your-namespace>.svc.cluster.local
  - metrics-webhook-service.<your-namespace>.svc
  - metrics-webhook-service.<your-namespace>.svc.cluster.local
  issuerRef:
    kind: Issuer
    name: klt-selfsigned-issuer
  secretName: klt-certs
---
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: klt-selfsigned-issuer
  namespace: <your-namespace>
spec:
  selfSigned: {}
```

Note the following about these fields:

* The `apiVersion` field refers to the API for the cert-manager.
* The value of the `.spec.secretName` field as well as the `.metadata.name` of the `Certificate` CR
  must be `klt-certs`.
* Substitute the namespace placeholders with your namespace, where Keptn is installed.

See the [CA Injector](https://cert-manager.io/docs/concepts/ca-injector/)
documentation for more details.
