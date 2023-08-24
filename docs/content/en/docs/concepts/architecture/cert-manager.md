---
title: Keptn Certificate Manager
description: Learn how the cert-manager works
layout: quickstart
weight: 100
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

### Keptn Cert Manager

The Keptn Cert Manager automatically configures TLS certificates to
[secure communication with the Kubernetes API](https://kubernetes.io/docs/concepts/security/controlling-access/#transport-security).
You can instead
[configure your own certificate manager](https://lifecycle.keptn.sh/docs/install/cert-manager/)
for this purpose.

The Lifecycle Toolkit includes a Mutating Webhook
that requires TLS certificates to be mounted as a volume in its pod.
In version 0.6.0 and later, the certificate creation
is handled automatically by
the [klt-cert-manager](https://github.com/keptn/lifecycle-toolkit/blob/main/klt-cert-manager/README.md).

How it works:

* The certificate is created as a secret
in the `keptn-lifecycle-toolkit-system` namespace
with a renewal threshold of 12 hours.
* If the certificate expires,
the [klt-cert-manager](https://github.com/keptn/lifecycle-toolkit/blob/main/klt-cert-manager/README.md)
renews it.
* The Lifecycle Toolkit operator waits for a valid certificate to be ready.
* When the certificate is ready,
  it is mounted on an empty dir volume in the operator.

`klt-cert-manager` is a customized certificate manager
that is installed with the Lifecycle Toolkit by default.
It is included to simplify installation for new users
and because it is much smaller than most standard certificate managers.
However, KLT is compatible with most certificate managers
and can be configured to use another certificate manager if you prefer.
See [Use Keptn with cert-manager.io](../../operate/cert-manager.md)
for instructions.

## Invalid certificate errors

When a certificate is left over from an older version,
the webhook or the operator may generate errors
because of an invalid certificate.
To solve this, delete the certificate and restart the operator.

The KLT cert-manager certificate is stored as a secret in the `klt` namespace.
To retrieve it:

```shell
kubectl get secrets -n keptn-lifecycle-toolkit-system
```

This returns something like:

```shell
NAME                        TYPE                 DATA   AGE
klt-certs                   Opaque               5      4d23h
```

Specify the `NAME` of the KLT certificate (`klt-certs` in this case)
to delete the KLT certificate:

```shell
kubectl delete secret klt-certs -n keptn-lifecycle-toolkit-system
```
