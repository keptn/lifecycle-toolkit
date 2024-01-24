---
comments: true
---

# Keptn + cert-manager.io

Keptn includes
a light-weight, customized cert-manager
that is used to register Webhooks to the [KubeAPI](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/).
Bundling the cert-manager simplifies the installation for new users
and provides the functionality Keptn needs
without the overhead of other cert-managers.
For a description of the architecture, see
[Keptn Certificate Manager](../../components/certificate-operator.md).

Keptn also works well with `cert-manager.io`.
If you are already using `cert-manager.io`,
you can continue to use it for other components
and use the Keptn cert-manager just for Keptn activities
or you can disable the Keptn cert-manager
and configure Keptn to use `cert-manager.io`.

If you want Keptn to use `cert-manager.io`,
you must configure it *before* you install Keptn.
The steps are:

* Install `cert-manager.io` if it is not already installed.
* Add the `Certificate` and `Issuer` CRs for `cert-manager.io`.
* (optional) Install Keptn without the built-in `keptn-cert-manager`
and with injected CA annotations via Helm

## Add the CR(s) for cert-manager.io

These are the CRs for `cert-manager.io` to be applied to your cluster:

```yaml
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: keptn-certs
  namespace: <keptn-namespace>
spec:
  dnsNames:
  - lifecycle-webhook-service.<keptn-namespace>.svc
  - lifecycle-webhook-service.<keptn-namespace>.svc.cluster.local
  - metrics-webhook-service.<keptn-namespace>.svc
  - metrics-webhook-service.<keptn-namespace>.svc.cluster.local
  issuerRef:
    kind: Issuer
    name: keptn-selfsigned-issuer
  secretName: keptn-certs
---
apiVersion: cert-manager.io/v1
kind: Issuer
metadata:
  name: keptn-selfsigned-issuer
  namespace: <keptn-namespace>
spec:
  selfSigned: {}
```

Note the following about these fields:

* The `apiVersion` field refers to the API for the cert-manager.
* The value of the `.spec.secretName` field as well as the `.metadata.name` of the `Certificate` CR
  must be `keptn-certs`.
* Substitute the namespace placeholders with your namespace, where Keptn is installed.

## Injecting CA Annotations

`cert-manager.io` supports specific annotations for
injectable resources depending on the injection source.
To configure these annotations, modify the `global.caInjectionAnnotation` Helm value.
See the [CA Injector](https://cert-manager.io/docs/concepts/ca-injector/) documentation for more details.

Here is an example `values.yaml` file demonstrating the configuration of CA injection
by using the `cert-manager.io/inject-ca-from` annotation:

```yaml
global:
  certManagerEnabled: false # disable Keptn Cert Manager
  caInjectionAnnotations:
    cert-manager.io/inject-ca-from: keptn-system/keptn-certs
```

Refer to the
[Customizing the configuration of components](../index.md#customizing-the-configuration-of-components)
for more details.
