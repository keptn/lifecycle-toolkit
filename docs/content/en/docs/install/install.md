---
title: Install KLT
description: Install the Keptn Lifecycle Toolkit
weight: 10
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

Keptn Lifecycle Toolkit works whether or not you use a GitOps strategy.
The following is an imperative walkthrough.

If you prefer a GitOps / declarative-based approach follow [this demo instead](https://example.com).

## Prerequisites

- A Kubernetes cluster > 1.24 (we recommend [Kubernetes kind](https://kind.sigs.k8s.io/docs/user/quick-start/))
  (`kind create cluster`)
- [Helm](https://helm.io) CLI available

## Objectives

- Install Keptn Lifecycle Toolkit on your cluster
- Annotate a namespace and deployment to enable Keptn Lifecycle Toolkit
- View DORA Metrics
- Install Grafana and Observability tooling to view DORA metrics

## System Overview

By the end of this page, here is what will be built.
This system will be built in stages.

![system overview](/docs/install/assets/install01.png)

## The Basics: A Deployment, Keptn and DORA Metrics

Let's start with the basics.
Here is what we will now build.

A deployment will occur.
Keptn will monitor the deployment and generate:

- An OpenTelemetry trace per deployment
- DORA metrics

![the basics](/docs/install/assets/install02.png)

Notice though that the metrics and traces have nowhere to go.
That will be fixed in a subsequent step.

## Step 1: Install Keptn Lifecycle Toolkit

Install Keptn Lifecycle Toolkit using Helm:

```shell
helm repo add klt https://charts.lifecycle.keptn.sh
helm repo update
helm upgrade --install keptn klt/klt -n keptn-lifecycle-toolkit-system --create-namespace --wait
```

Keptn will need to know where to send OpenTelemetry traces.
Of course, Jaeger is not yet installed so traces have nowhere to go (yet),
but creating this configuration now means the system is preconfigured.

Save this file as `collectorconfig.yaml`:

```yaml
---
apiVersion: options.keptn.sh/v1alpha1
kind: KeptnConfig
metadata:
  name: keptnconfig-sample
  namespace: keptn-lifecycle-toolkit-system
spec:
  OTelCollectorUrl: 'jaeger-collector.keptn-lifecycle-toolkit-system.svc.cluster.local:4317'
  keptnAppCreationRequestTimeoutSeconds: 30
```

Apply the file and restart KLT to pick up the new config:

```shell
kubectl apply -f collectorconfig.yaml
kubectl rollout restart deployment -n keptn-lifecycle-toolkit-system -l control-plane=lifecycle-operator
kubectl rollout status deployment -n keptn-lifecycle-toolkit-system -l control-plane=lifecycle-operator --watch
kubectl rollout restart deployment -n keptn-lifecycle-toolkit-system -l component=scheduler
kubectl rollout status deployment -n keptn-lifecycle-toolkit-system -l component=scheduler --watch
```

## Create Namespace for Demo Application

Save this file as `namespace.yaml`.
The annotation means that Keptn Lifecycle Toolkit is active for workloads in this namespace.

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: keptndemo
  annotations:
    keptn.sh/lifecycle-toolkit: enabled
```

Create the namespace:

```shell
kubectl apply -f namespace.yaml
```

## Deploy Demo Application

It is time to deploy the demo application.

Save this manifest as `app.yaml`:

```shell
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  namespace: keptndemo
  labels:
    app.kubernetes.io/name: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app.kubernetes.io/name: nginx
  template:
    metadata:
      labels:
        app.kubernetes.io/part-of: keptndemoapp
        app.kubernetes.io/name: nginx
        app.kubernetes.io/version: 0.0.1
    spec:
      containers:
      - name: nginx
        image: nginx:1.14.2
        ports:
        - containerPort: 80
---
apiVersion: apps/v1
kind: Service
apiVersion: v1
kind: Service
metadata:
  name: nginx
  namespace: keptndemo
spec:
  selector:
    app.kubernetes.io/name: nginx
  ports:
    - protocol: TCP
      port: 8080
      targetPort: 80
```

Now apply:

```shell
kubectl apply -f app.yaml
```

Keptn looks for these 3 labels:

- `app.kubernetes.io/part-of`
- `app.kubernetes.io/name`
- `app.kubernetes.io/version`

These are [Kubernetes recommended labels](https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/#labels)
but if you want to use different labels, you can swap for:

- `keptn.sh/app` instead of `app.kubernetes.io/part-of`
- `keptn.sh/workload` instead of `app.kubernetes.io/name`
- `keptn.sh/version` instead of `app.kubernetes.io/version`

## Explore Keptn

Keptn is now aware of your deployments and is generating DORA statistics about them.

Keptn has created a CRD to track your application.
The name of which is based on the `part-of` label.

It may take up to 30 seconds to create the `KeptnApp` so run the following command until you see the `keptnappdemo` CRD.

```shell
kubectl -n keptndemo get keptnapp
```

Expected output:

```shell
NAME           AGE
keptndemoapp   2s
```

Keptn also creates a new application version every time you increment the `version` label..

The `PHASE` will change as the deployment progresses.
A successful deployment is shown as `PHASE=Completed`

```shell
kubectl -n keptndemo get keptnappversion
```

Expected output:

```shell
NAME                      APPNAME        VERSION   PHASE
keptndemoapp-0.0.1-***    keptndemoapp   0.0.1     Completed
```

Keptn can run tasks and SLO evaluations before and after deployment.
You haven't configured this yet, but you can see the full lifecycle for a `keptnappversion` by running:

```shell
kubectl -n keptndemo get keptnappversion -o wide
```

Keptn applications are a collection of workloads.
By default, Keptn will build `KeptnApp` CRDs based on the labels you provide.

In the example above, the `KeptnApp` called `keptndemoapp` contains one workload (based on the `name` label):

## View your application

Port-forward to expose your app on `http://localhost:8080`:

```shell
kubectl -n keptndemo port-forward svc/nginx 8080
```

You should see the "Welcome to nginx" page.

![nginx demo app](/docs/install/assets/nginx.png)

## View DORA Metrics

Keptn is generating DORA metrics and OpenTelemetry traces for your deployments.

These metrics are exposed via the Keptn lifecycle operator `/metrics` endpoint on port `2222`.

To see these raw metrics, port-forward to the lifecycle operator metrics service:

```shell
kubectl -n keptn-lifecycle-toolkit-system port-forward service/keptn-klt-lifecycle-operator-metrics-service 2222
```

Access metrics in Prometheus format on `http://localhost:2222/metrics`.
Look for metrics starting with `keptn_`.

![keptn prometheus metrics](/docs/install/assets/keptnprommetrics.png)

## Make DORA metrics more user friendly

It is much more user friendly to provide dashboards for metrics, logs and traces.
So let's install new Observability components to help us:

- Cert manager: Jaeger requires cert-manager
- Jaeger: Store and view DORA deployment traces
- Prometheus: Store DORA metrics
- OpenTelemetry collector: Scrape metrics from the above DORA metrics endpoint & forward to Prometheus
- Grafana (and some prebuilt dashboards): Visualise the data

![add observability](/docs/install/assets/install01.png)

## Install Cert Manager

Jaeger requires Cert Manager, so install it now:

```shell
kubectl apply -f https://github.com/cert-manager/cert-manager/releases/download/v1.12.2/cert-manager.crds.yaml
helm repo add jetstack https://charts.jetstack.io
helm repo update
helm install cert-manager --namespace cert-manager --version v1.12.2 jetstack/cert-manager --create-namespace --wait
```

## Install Jaeger

Save this file as `jaeger.yaml`:

```shell
apiVersion: jaegertracing.io/v1
kind: Jaeger
metadata:
  name: jaeger
spec:
  strategy: allInOne
```

Install Jaeger to store and visualise the deployment traces generated by Keptn:

```shell
kubectl create namespace observability
kubectl apply -f https://github.com/jaegertracing/jaeger-operator/releases/download/v1.46.0/jaeger-operator.yaml -n observability
kubectl wait --for=condition=available deployment/jaeger-operator -n observability --timeout=300s
kubectl apply -f jaeger.yaml -n keptn-lifecycle-toolkit-system
kubectl wait --for=condition=available deployment/jaeger -n keptn-lifecycle-toolkit-system --timeout=300s
```

Port-forward to access Jaeger:

```shell
kubectl -n keptn-lifecycle-toolkit-system port-forward svc/jaeger-query 16686
```

Jaeger is available on `http://localhost:16686`

## Install Grafana dashboards

Create some Keptn Grafana dashboards that will be available when Grafana is installed and started:

```shell
kubectl create ns monitoring
kubectl apply -f https://raw.githubusercontent.com/keptn/lifecycle-toolkit/main/examples/support/observability/config/prometheus/grafana-config.yaml
kubectl apply -f https://raw.githubusercontent.com/keptn/lifecycle-toolkit/main/examples/support/observability/config/prometheus/grafana-dashboard-keptn-applications.yaml
kubectl -n monitoring label cm/grafana-dashboard-keptn-applications grafana_dashboard="1"
kubectl apply -f https://raw.githubusercontent.com/keptn/lifecycle-toolkit/main/examples/support/observability/config/prometheus/grafana-dashboard-keptn-overview.yaml
kubectl -n monitoring label cm/grafana-dashboard-keptn-overview grafana_dashboard="1"
kubectl apply -f https://raw.githubusercontent.com/keptn/lifecycle-toolkit/main/examples/support/observability/config/prometheus/grafana-dashboard-keptn-workloads.yaml
kubectl -n monitoring label cm/grafana-dashboard-keptn-workloads grafana_dashboard="1"
```

## Install Grafana datasources

This file will configure Grafana to look at the Jaeger service and the Prometheus service on the cluster.

Save this file as `datasources.yaml`:

```yaml
apiVersion: v1
kind: Secret
type: Opaque
metadata:
  labels:
    grafana_datasource: "1"
  name: grafana-datasources
  namespace: monitoring
stringData:
  datasources.yaml: |-
    {
        "apiVersion": 1,
        "datasources": [
            {
                "access": "proxy",
                "editable": false,
                "name": "prometheus",
                "orgId": 1,
                "type": "prometheus",
                "url": "http://observability-stack-kube-p-prometheus.monitoring.svc:9090",
                "version": 1
            },
            {
                "orgId":1,
                "name":"Jaeger",
                "type":"jaeger",
                "typeName":"Jaeger",
                "typeLogoUrl":"public/app/plugins/datasource/jaeger/img/jaeger_logo.svg",
                "access":"proxy",
                "url":"http://jaeger-query.keptn-lifecycle-toolkit-system.svc.cluster.local:16686",
                "user":"",
                "database":"",
                "basicAuth":false,
                "isDefault":false,
                "jsonData":{"spanBar":{"type":"None"}},
                "readOnly":false
            }
        ]
    }
```

Now apply it:

```shell
kubectl apply -f datasources.yaml
```

## Install kube prometheus stack

This will install:

- Prometheus
- Prometheus Configuration
- Grafana & default dashboards

Save this file as `values.yaml`:

```yaml
grafana:
  adminPassword: admin
  sidecar.datasources.defaultDatasourceEnabled: false
prometheus:
  prometheusSpec:
    additionalScrapeConfigs:
      - job_name: "scrape_klt"
        scrape_interval: 5s
        static_configs:
          - targets: ['keptn-klt-lifecycle-operator-metrics-service.keptn-lifecycle-toolkit-system.svc.cluster.local:2222']
```

```shell
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
helm upgrade --install observability-stack prometheus-community/kube-prometheus-stack --version 48.1.1 --namespace monitoring --values=values.yaml --wait
```

## Access Grafana

```shell
kubectl -n monitoring port-forward svc/observability-stack-grafana 80
```

- Grafana username: `admin`
- Grafana password: `admin`

View the Keptn dashboards at: `http://localhost/dashboards`

Remember that Jaeger and Grafana weren't installed during the first deployment
so expect the dashboards to look a little empty.

## Deploy v0.0.2 and populate Grafana

By triggering a new deployment, Keptn will track this deployment and the Grafana dashboards will actually have data.

Modify your `app.yaml` and change the `app.kubernetes.io/version` from `0.0.1` to `0.0.2`.

Apply your update:

```shell
kubectl apply -f app.yaml
```

After about 30 seconds you should now see two `keptnappversions`:

```shell
kubectl -n keptndemo get keptnappversion
```

Expected output:

```shell
NAME                          APPNAME        VERSION   PHASE
keptndemoapp-0.0.1-6b86b273   keptndemoapp   0.0.1     Completed
keptndemoapp-0.0.2-d4735e3a   keptndemoapp   0.0.2     AppDeploy
```

Wait until the `PHASE` of `keptndemoapp-0.0.2` is `Completed`.
This signals that the deployment was successful and the pod is running.

View the Keptn Applications Dashboard and you should see the DORA metrics and an OpenTelemetry per trace:

![keptn applications dashboard](/docs/install/assets/keptnapplications.png)

![deployment trace](/docs/install/assets/deploymenttrace.png)

## More control over KeptnApp

You may have noticed that the `KeptnApp` Custom Resources are created automatically by KLT.

The lifecycle toolkit automatically groups workloads into `KeptnApp` by looking for matching `part-of` annotations.
Any workloads with the same `part-of` annotation is said to be `part-of` the same `KeptnApp`.

However, you can override this automatic behaviour by creating a custom `KeptnApp` CRD.
In this way, you are in full control of what constitutes a Keptn Application.
See [Define a Keptn Application](../implementing/integrate/#define-keptnapp-manually) for more information.

## What's next?

Keptn can run pre and post deployment tasks and SLO evaluations automatically.

Continue to Keptn learning journey by [adding deployment tasks](https://example.com).
