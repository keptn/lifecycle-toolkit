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
[use cert-manager.io](../../installation/configuration/cert-manager.md)
for this purpose.

Keptn includes a Mutating Webhook
that requires TLS certificates to be mounted as a volume in its pod.
In version 0.6.0 and later, the certificate creation
is handled automatically by
the [keptn-cert-manager](https://github.com/keptn/lifecycle-toolkit/blob/main/keptn-cert-manager/README.md).

How it works:

* The certificate is created as a secret
in the `keptn-system` namespace
with a renewal threshold of 12 hours.
* If the certificate expires,
the [keptn-cert-manager](https://github.com/keptn/lifecycle-toolkit/blob/main/keptn-cert-manager/README.md)
renews it.
* The Keptn `lifecycle-operator` waits for a valid certificate to be ready.
* When the certificate is ready,
  it is mounted on an empty dir volume in the operator.

`keptn-cert-manager` is a customized certificate manager
that is installed with Keptn by default.
It is included to simplify installation for new users
and because it is much smaller than most standard certificate managers.
However, Keptn is compatible with most certificate managers
and can be configured to use another certificate manager if you prefer.
See [Use Keptn with cert-manager.io](../../installation/configuration/cert-manager.md)
for instructions.

## Invalid certificate errors

When a certificate is left over from an older version,
the webhook or the operator may generate errors
because of an invalid certificate.
To solve this, delete the certificate and restart the operator.

The Keptn cert-manager certificate is stored as a secret in the
`keptn-system` namespace.
To retrieve it:

```shell
kubectl get secrets -n keptn-system
```

This returns something like:

```shell
NAME                        TYPE                 DATA   AGE
keptn-certs                   Opaque               5      4d23h
```

Specify the `NAME` of the Keptn certificate (`keptn-certs` in this case)
to delete the Keptn certificate:

```shell
kubectl delete secret keptn-certs -n keptn-system
```
