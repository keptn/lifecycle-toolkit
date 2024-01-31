#  Add metadata to your application

This guide will walk you through the usage of Context Metadata in Keptn.

After reading this guide you will be able to: 
- add metadata to applications, workloads
and tasks in two ways.
- add metadata information to traces

## Introduction

`KeptnAppContext` is a custom resource in Keptn that allows you to add metadata
and links to traces for a specific application.
This enables you to enrich your Keptn resources and your traces with additional
information, making it easier to understand and analyze
the performance of your applications.

Common uses of this feature include:

- adding commit id or references related to your GitOps and CI-CD tooling
- referencing different stages and actors 
such as who committed this change? In what repo? What is the associated ticket? Was it a hot fix? )

To exploit this feature you need to install Keptn and deploy an application,
for instance, you can follow [step #1-#3 here](../getting-started/observability.md#step-1-install-keptn)

To collect traces you will require Jaeger.
To visualise and inspect the traces, you can either use the Jaeger UI or Grafana.

- Install
  [Grafana](https://grafana.com/grafana/)
  following the instructions in [Grafana Setup](https://grafana.com/docs/grafana/latest/setup-grafana/)
  or the visualization tool of your choice.
- Install
  [Jaeger](https://www.jaegertracing.io/)
  or a similar tool for traces following the instructions in
  [Jaeger Setup](https://www.jaegertracing.io/docs/1.50/getting-started/).

## Include metadata in workload traces

To enrich the workload traces with custom metadata, you can use the
`keptn.sh/metadata` annotation in your deployment.
The values specified in the annotation
are added as key-value attributes to the workload trace.

Modify your deployment file adding`keptn.sh/metadata: "stage=dev"` annotation.

To see the changes Keptn must redeploy: increment the `app.kubernetes.io/version` value
(ex. if you are following our getting started guide change the version
from `0.0.2` to `0.0.3`) or change the `keptn.sh/version` value
if you used the Keptn specific labels in your deployment yaml file.

This way, after the re-deployment, the workload trace will contain the `stage=dev` attribute.

## Include metadata in application traces

Similar to the previous step, custom metadata can also be added to application traces via the
[KeptnAppContext](../reference/api-reference/lifecycle/v1beta1/index.md#keptnappcontext) custom resource.
In the `.spec.metadata` field you can define multiple key-value pairs, which are propagated
to the application trace as attributes in the same manner as for workloads.

The key-value pairs that are added to the application trace are also added
to each workload trace that is part of the application.
If the same key is specified for both
application and workload metadata attributes,
values specified for the workload take precedence.

A `KeptnAppContext` custom resource looks like the following:

```yaml
apiVersion: lifecycle.keptn.sh/v1beta1
kind: KeptnAppContext
metadata:
  name: keptndemoapp
  namespace: keptndemo
spec:
  metadata:
    commit-id: "1234"
    author: "myUser"
```

After applying the `KeptnAppContext` to your cluster, you need to increment the version of your
application by modifying your deployment file and changing the
value of the`app.kubernetes.io/version` field (or `keptn.sh/version` if you used the Keptn specific labels earlier).
(Alternatively you could apply the context resource before or together with the workloads)

After deploying the `KeptnAppContext` resource and re-deploying the application,
Keptn triggers another deployment of your application with the new context metadata, 
and all traces will contain commit-id=1234,author=myUser.
In other words, you should be able to see the application trace as well as the workload trace
contain the defined metadata as key-value attributes.


## Conclusion

Congratulations! You've learned how to use `KeptnAppContext` to add
metadata to applications and their traces.
This can be valuable for understanding the context of your traces and
establishing connections between
different versions and stages of your application.

What's next? Explore more about [traces on Keptn](./otel.md).
The paragraph on 
[linking traces between different application](./otel.md#advanced-tracing-configurations-in-keptn-linking-traces) 
also uses KeptnAppContext.  


