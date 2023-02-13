# Deploying an application using the Keptn Lifecycle Controller and ArgoCD

In this example, we will show you how to install our sample application *podtatohead* using the Keptn Lifecycle
Controller and [ArgoCD](https://argo-cd.readthedocs.io/en/stable/).

## TL;DR

* You can install ArgoCD and Keptn-lifecycle-toolkit using: `make install`
* Install argo CLI according to the instructions [here](https://argo-cd.readthedocs.io/en/stable/cli_installation/)
* Afterward, you can fetch the secret for the ArgoCD CLI using: `make argo-get-password`
* Then you can port-forward the ArgoUI using: `make port-forward-argocd`
  * Alternatively, you can access Argo using the CLI, configure it using `make argo-configure-cli`
* Deploy the PodTatoHead Demo Application: `make argo-install-podtatohead`
* Watch the progress on your ArgoUI: `http://localhost:8080`. Use the `admin` user and the password
  from `make argo-get-password`.

## Prerequisites

This tutorial assumes, that you already installed the Keptn Lifecycle Controller (
see <https://github.com/keptn/lifecycle-toolkit>). The installation instructions can be
found [here](https://github.com/keptn/lifecycle-toolkit#deploy-the-latest-release). Furthermore, you have to install
ArgoCD, as in the following their [installation instructions](https://argo-cd.readthedocs.io/en/stable/getting_started/)
.

### Install ArgoCD

If you don't have an already existing installation of ArgoCD, you
can [install](https://raw.githubusercontent.com/argoproj/argo-cd/v2.4.14/manifests/install.yaml) it using the following
commands:

```shell
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/v2.4.14/manifests/install.yaml
```

With these commands, ArgoCD will be installed in the `argocd` namespace.

After that, you can find the password for ArgoCD using the following command:

```shell
kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d
```

### Port-Forward ArgoCD and access the UI

To access the ArgoCD UI, you can port-forward the ArgoCD service using the following command:

```shell
kubectl port-forward svc/argocd-server -n argocd 8080:443
```

Then you can access the UI using <http://localhost:8080>.

## Installing the Demo Application

To install the demo application, you can use the following command (
apply [this manifest](https://raw.githubusercontent.com/keptn/lifecycle-toolkit/main/examples/support/argo/config/app.yaml)):

```shell
kubectl apply -f https://raw.githubusercontent.com/keptn-sandbox/lifecycle-toolkit-examples/main/support/argo/config/app.yaml
```

You will see that the application will be deployed using ArgoCD. You can watch the progress on the ArgoCD UI and should
see the following:
![img.png](assets/argo-screen.png)

In the meanwhile you can watch the progress of the deployment using:
```kubectl get pods -n podtato-kubectl```

* See that the pods are pending until the pre-deployment tasks have passed
* Pre-Deployment Tasks are started
* Pods get scheduled

```kubectl get keptnworkloadinstances -n podtato-kubectl```

* Get the current status of the workloads
* See in which phase your workload deployments are at the moment

```kubectl get keptnappversions -n podtato-kubectl```

*Get the current status of the application

* See in which phase your application deployment is at the moment

After some time all resources should be in a succeeded state. In the Argo-UI you will see that the application is in
sync.

<!-- markdownlint-disable-next-line MD033 MD013 -->
<img referrerpolicy="no-referrer-when-downgrade" src="https://static.scarf.sh/a.png?x-pxid=858843d8-8da2-4ce5-a325-e5321c770a78" />
