---
title: Keptn Certificate Manager
icon: concepts
layout: quickstart
weight: 5
hidechildren: true # this flag hides all sub-pages in the sidebar-multicard.html
---

### Keptn Cert Manager

The Lifecycle Toolkit includes a Mutating Webhook which requires TLS certificates to be mounted as a volume in its pod.
In version 0.6.0 and later, the certificate creation
is handled automatically by
the [klt-cert-manager](https://github.com/keptn/lifecycle-toolkit/blob/main/klt-cert-manager/README.md).

The certificate is created as a secret in the `keptn-lifecycle-toolkit-system` namespace with a renewal threshold of 12
hours.
If it expires, the [klt-cert-manager](https://github.com/keptn/lifecycle-toolkit/blob/main/klt-cert-manager/README.md)
renews it.
The Lifecycle Toolkit operator waits for a valid certificate to be ready.
The certificate is mounted on an empty dir volume in the operator.

When a certificate is left over from an older version, the webhook or the operator may generate errors because of an
invalid certificate. To solve this, delete the certificate and restart the operator.
