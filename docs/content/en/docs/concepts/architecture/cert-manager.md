---
title: Keptn Certificate Manager
description: Learn how the cert-manager works
icon: concepts
layout: quickstart
weight: 100
hidechildren: true # this flag hides all sub-pages in the sidebar-multicard.html
---

### Keptn Cert Manager

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

When a certificate is left over from an older version,
the webhook or the operator may generate errors
because of an invalid certificate.
To solve this, delete the certificate and restart the operator.

`klt-cert-manager` is a customized certificate manager
that is installed with the Lifecycle Toolkit by default.
It is included to simplify installation for new users
and because it is much smaller than most standard certificate managers.
However, KLT is compatible with most certificate managers
and can be configured to use another certificate manager if you prefer.
See [Use your own cert-manager](../../install/cert-manager)
for instructions.
