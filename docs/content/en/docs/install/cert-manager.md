---
title: Implement your own cert-manager (optional)
description: Replace the default KLT cert-manager
weight: 30
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

The Keptn Lifecycle Toolkit includes
a light-weight, customized cert-manager
that is used for installation and to implement Webhooks.
Bundling the cert-manager simplifies the installation for new users
and provides the functionality KLT needs
without the overhead of other cert-managers.

However, KLT works well with standard cert-managers.
You can redefine the cert-manager that KLT uses *before* you install KLT.

The steps are:

* Install the cert-manager of your choice.
* Modify the `Deployment` manifest of each KLT component.
* Add the `Certificate` CRD for the cert-manager you are using.

## Modify the KLT manifest

You must modify the KLT manifest for each KLT component
to make it aware of the cert-manager you are using.
The instructions here are for implementing
[cert-manager.io](https://cert-manager.io/);
the process is similar for other cert-managers.

To do this, change the `Deployment` manifest of each KLT component
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

The manifests must have the following special annotation:

```yaml
cert-manager.io/inject-ca-from=klt-serving-cert/keptn-lifecycle-toolkit-system
```

The value of the annotation must match the
`name/namespace` of the cert-manager CRD discussed below.

## Add the CRD for your cert-manager

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
  :secretName webhook-server-cert // this has to match the name of the "secretName" field in the volume definition of step 1
```

Note the following about these fields:

* The `apiVersion` field refers to the API for the cert-manager.
* The `metadata` section includes two fields.
  The value of these fields must match the annotations
* The value of the `secretName` field
  must match the value of the `secretName` field used
  in the `volumes` definition section of the KLT manifests above.

See the [CA Injector](https://cert-manager.io/docs/concepts/ca-injector/)
documentation for more details.
