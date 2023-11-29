# Redeploy/Restart an Application

A [KeptnApp](../reference/crd-reference/app.md) can fail
because of an unsuccessful pre-deployment evaluation
or pre-deployment task.
For example, this happens if the target value of a
[KeptnEvaluationDefinition](../reference/crd-reference/evaluationdefinition.md)
resource is misconfigured
or a pre-deployment evaluation checks the wrong URL.

After you fix the configuration
that caused the pre-deployment evaluation or task to fail,
you can increment the `spec.revision` value
and apply the updated `KeptnApp` manifest
to create a new revision of the `KeptnApp`
without modifying the `version`.

Afterwards, all related `KeptnWorkloadVersions`
automatically refer to the newly created revision of the `KeptnAppVersion`
to determine whether they are allowed
to enter their respective deployment phases.

To illustrate this, consider the following example:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: restartable-apps
  annotations:
    keptn.sh/lifecycle-toolkit: "enabled"
