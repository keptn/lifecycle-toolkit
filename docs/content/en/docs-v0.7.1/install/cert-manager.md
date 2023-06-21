---
title: Use your own cert-manager (optional)
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

KLT, however, works well with standard cert-managers.
The KLT cert-manager can also coexist with another cert-manager.
If you are already using a different cert-manager,
you can continue to use that cert-manager for other components
and use the KLT cert-manager just for KLT activities
or you can configure KLT to use that cert-manager.

If you want KLT to use your cert-manager,
you must configure it *before* you install KLT.
The steps are:

* Install the cert-manager of your choice
  if it is not already installed.
* Modify the `Deployment` manifest of each KLT operator component.
* Add the `Certificate` CRD for the cert-manager you are using.

## Modify the KLT manifest

You must modify the KLT manifest for each KLT operator component
to make it aware of the cert-manager you are using.
These instructions implement
[cert-manager.io](https://cert-manager.io/);
the process is similar for other cert-managers.

To configure KLT to use your cert-manager,
change the `Deployment` manifest of each KLT operator component
and **replace** the following `volumes` definition

   ```yaml
   - emptyDir: {}
     name: certs-dir
   ```

   with

   ```yaml
   - name: cert
     secret:
       defaultMode: 420
       secretName: webhook-server-cert
   ```

Each manifest must have the following special annotation:

```yaml
cert-manager.io/inject-ca-from=klt-serving-cert/keptn-lifecycle-toolkit-system
```

The value of the annotation must match the
`name/namespace` of the cert-manager CRD discussed below.

## Add the CRD for your cert-manager

This is the CRD for `cert-manager.io`:

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: klt-serving-cert 
  namespace: keptn-lifecycle-toolkit-system
spec:
  dnsNames:
  - lifecycle-webhook-service.keptn-lifecycle-toolkit-system.svc
  - lifecycle-webhook-service.keptn-lifecycle-toolkit-system.svc.cluster.local
  issuerRef:
    kind: Issuer
    name: klt-selfsigned-issuer
  secretName webhook-server-cert
```

Note the following about these fields:

* The `apiVersion` field refers to the API for the cert-manager.
* The `metadata` section includes two fields.
  The value of these fields must match the annotations
  used in the KLT operator manifests.
* The value of the `secretName` field
  must match the value of the `secretName` field used
  in the `volumes` definition section of the KLT operator manifests above.

See the [CA Injector](https://cert-manager.io/docs/concepts/ca-injector/)
documentation for more details.
