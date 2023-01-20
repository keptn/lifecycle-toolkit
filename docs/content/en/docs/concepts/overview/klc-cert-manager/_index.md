---
title: How it works
icon: concepts
layout: quickstart
weight: 5
hidechildren: true # this flag hides all sub-pages in the sidebar-multicard.html
---

### Keptn Cert Manager

The Lifecycle Toolkit includes a Mutating Webhook which requires TLS certificates to be mounted as a volume in its pod. Since version 0.6.0, the certificate creation
is handled automatically by [klt-cert-manager](https://github.com/keptn/lifecycle-toolkit/blob/main/klt-cert-manager/README.md).

The certificate is created as a secret in the keptn-lifecycle-toolkit-system namespace with a renewal threshold of 12 hours. The Lifecycle Toolkit operator will be waiting for a valid certificate to be ready.
The certificate will be mounted on an empty dir volume in the operator.

#### FAQ
In case of certificates left over from a older version it could happen that the webhook or the operator errors due to invalid certificate. To solve this it is enough to delete de certificate and restart the operator.

