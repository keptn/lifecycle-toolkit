# klt-cert-manager

The Keptn certificate manager ensures that the webhooks of an operator can obtain a valid certificate
to access the Kubernetes API server.

## Description

This `klt-cert-manager` operator should only be installed when paired with the Lifecycle Toolkit operator versions 0.6.0
or above.
The TLS certificate is mounted as a volume in the LT operator pod and is renewed every 12 hours or every time the LT
operator deployment changes.
The `klt-cert-manager` retrieves all `MutatingWebhookConfigurations`, `ValidatingWebhookConfigurations` and
`CustomResourceDefinitions` based on a label selector that can be defined using the following environment variables:

- `LABEL_SELECTOR_KEY`: Label key used or identifying resources for certificate injection.
Default: `keptn.sh/inject-cert`
- `LABEL_SELECTOR_VALUE`: Label value used for identifying resources for certificate injection.
Default: `true`.

Using these label selectors, `MutatingWebhookConfigurations`, `ValidatingWebhookConfigurations` and
`CustomResourceDefinitions` can be enabled for certificate injection by adding the required labels to their metadata:

````yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  labels:
    keptn.sh/inject-cert: true
  name: keptnconfigs.options.keptn.sh
````

## Using the klt-cert-manager library

The functionality provided by this operator can also be added to other operators by using the `klt-cert-manager` as
a library.
To do this, add the library as a dependency to your application:

```shell
go get github.com/keptn/lifecycle-toolkit/klt-cert-manager
```

Then, in your operator's setup logic, an instance of the `KeptnWebhookCertificateReconciler` can be
created and registered to your operator's controller manager:

```golang
package main

import (
    "flag"
    "log"
    "os"
    "sigs.k8s.io/controller-runtime/pkg/webhook/admission"

    "github.com/keptn/lifecycle-toolkit/klt-cert-manager/controllers/keptnwebhookcontroller"
	"github.com/keptn/lifecycle-toolkit/klt-cert-manager/pkg/certificates"
	certCommon "github.com/keptn/lifecycle-toolkit/klt-cert-manager/pkg/common"
    "github.com/keptn/lifecycle-toolkit/klt-cert-manager/pkg/webhook"
    // +kubebuilder:scaffold:imports
)

func main() {
    // operator setup ... 
    certificateReconciler := keptnwebhookcontroller.NewReconciler(keptnwebhookcontroller.CertificateReconcilerConfig{
        Client:    mgr.GetClient(),
        Scheme:    mgr.GetScheme(),
        Log:       ctrl.Log.WithName("KeptnWebhookCert Controller"),
        Namespace: "my-namespace",
        MatchLabels: map[string]string{
            "inject-cert": "true",
        },
    })
    if err = certificateReconciler.SetupWithManager(mgr); err != nil {
        setupLog.Error(err, "unable to create controller", "controller", "Deployment")
        os.Exit(1)
    }
    //...
    // register mutating/validating webhooks
    webhookBuilder := webhook.NewWebhookBuilder().
        SetNamespace(env.PodNamespace).
        SetPodName(env.PodName).
        SetConfigProvider(cmdConfig.NewKubeConfigProvider()).
		SetManagerProvider(
			webhook.NewWebhookManagerProvider(
				mgr.GetWebhookServer().CertDir, "tls.key", "tls.crt"),
		).
		SetCertificateWatcher(
			certificates.NewCertificateWatcher(
				mgr.GetAPIReader(),
				mgr.GetWebhookServer().CertDir,
				env.PodNamespace,
				certCommon.SecretName,
				setupLog,
			),
		)

    setupLog.Info("starting webhook and manager")
    if err := webhookBuilder.Run(mgr, map[string]*admission.Webhook{
            "/webhook-path": &webhook.Admission{},
        }); err != nil {
        setupLog.Error(err, "problem running manager")
        os.Exit(1)
    }
}
```

Using this approach will require the following `ClusterRole` permissions to be bound to your operator's ServiceAccount:

```yaml
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: certificate-operator-role
rules:
- apiGroups:
  - admissionregistration.k8s.io
  resources:
  - mutatingwebhookconfigurations
  verbs:
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - admissionregistration.k8s.io
  resources:
  - validatingwebhookconfigurations
  verbs:
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - apiextensions.k8s.io
  resources:
  - customresourcedefinitions
  verbs:
  - get
  - list
  - patch
  - update
  - watch
- apiGroups:
  - apps
  resources:
  - deployments
  verbs:
  - get
  - list
  - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: certificate-operator-role
  namespace: my-operator-system
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
```

The required permissions can also be reduced by only allowing access to a specific set of resources, by providing
the `KeptnWebhookCertificateReconciler` with a list of resources that should by enabled for certificate injection,
instead of specifying the `MatchLabels`:

```golang
package main

import (
    "flag"
    "log"
    "os"
    
    "github.com/keptn/lifecycle-toolkit/klt-cert-manager/controllers/keptnwebhookcontroller"
    // +kubebuilder:scaffold:imports
)

func main() {
    // operator setup ... 
    certificateReconciler := keptnwebhookcontroller.NewReconciler(keptnwebhookcontroller.CertificateReconcilerConfig{
        Client:        mgr.GetClient(),
        Scheme:        mgr.GetScheme(),
        Log:           ctrl.Log.WithName("KeptnWebhookCert Controller"),
        Namespace:     "my-namespace",
        WatchResources: &keptnwebhookcontroller.ObservedObjects{
            MutatingWebhooks:          []string{"my-mwh-1", "my-mwh-2"},
            ValidatingWebhooks:        []string{"my-vwh-1", "my-vwh-2"},
            CustomResourceDefinitions: []string{"my-crd-1", "my-crd-2"},
            Deployments:               []string{"my-operator-deployment"},
        },
    })
    if err = certificateReconciler.SetupWithManager(mgr); err != nil {
        setupLog.Error(err, "unable to create controller", "controller", "Deployment")
        os.Exit(1)
    }
    //...
}
```

Using this configuration, you can limit the required permissions as follows:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRole
metadata:
  name: certificate-operator-role
rules:
- apiGroups:
  - admissionregistration.k8s.io
  resources:
  - mutatingwebhookconfigurations
  verbs:
  - get
  - patch
  - update
  - watch
  resourceNames:
    - my-mwh-1
    - my-mwh-2
- apiGroups:
  - admissionregistration.k8s.io
  resources:
  - validatingwebhookconfigurations
  verbs:
  - get
  - patch
  - update
  - watch
  resourceNames:
    - my-vwh-1
    - my-vwh-2
- apiGroups:
  - apiextensions.k8s.io
  resources:
  - customresourcedefinitions
  verbs:
  - get
  - patch
  - update
  - watch
  resourceNames:
    - my-crd-1
    - my-crd-2
- apiGroups:
  - apps
  resources:
  - deployments
  verbs:
  - get
  - list
  - watch
---
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: certificate-operator-role
  namespace: my-operator-system
rules:
- apiGroups:
  - ""
  resources:
  - secrets
  verbs:
  - create
  - delete
  - get
  - list
  - patch
  - update
  - watch
```

## Getting Started

You’ll need a Kubernetes cluster to run against.
You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for
testing, or run against a remote cluster.

> **Note**
Your controller will automatically use the current context in your kubeconfig file (i.e. whatever
cluster `kubectl cluster-info` shows).

### Running on the cluster

1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

1. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/cert-manager:tag
```

1. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/cert-manager:tag
```

### Uninstall CRDs

To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller

UnDeploy the controller to the cluster:

```sh
make undeploy
```

## Contributing

### How it works

This project aims to follow the
Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/)
which provides a reconcile function responsible for synchronizing resources untile the desired state is reached on the
cluster

### Test It Out

1. Install the CRDs into the cluster:

```sh
make install
```

1. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

> **Note**
You can also run this in one step by running: `make install run`

### Modifying the API definitions

If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

> **Note**
Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)
