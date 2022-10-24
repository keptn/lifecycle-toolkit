# Keptn Application Controller

<details>
<summary>Table of Contents</summary>

<!-- run with: npx markdown-toc --no-first-h1 --no-stripHeadingTags -i file.md -->

<!-- toc -->

- [Reconciliation loop](#reconciliation-loop)
- [Telemetry](#telemetry)

<!-- tocstop -->

</details>

The Keptn Application Controller watches two CRDs:

- [KeptnApp][]
- [KeptnAppVersion][]

## Reconciliation loop

In the reconciliation loop the Application controller MUST look for the [KeptnApp][] that triggered the loop, called `ka`.
Afterwards, the controller MUST look for a [KeptnAppVersion][] in the namespace of `ka` with name tbd...



## Telemetry

The controller should create two OpenTelemetry spans: 

- `appversion_deployment` which exposes the deployment progress,
- and `reconcile_app` which exposes the LFC progress in getting to the specified state.
These spans MUST be annotated with the following attributes:

- `keptn.deployment.app.name` with the name of the KeptnApp that is being deployed
- ...
tbd

Furthermore, for each important step the controller MUST emit a [K8s event](https://kubernetes.io/docs/reference/kubernetes-api/cluster-resources/event-v1/)
namely:

- when the [KeptnAppVersion][] cannot be created;
- when the [KeptnAppVersion][] cannot is successfully created.



[KeptnApp]: ../crds/v1alpha1/KeptnApp.md
[KeptnAppVersion]: ../crds/v1alpha1/KeptnAppVersion.md
