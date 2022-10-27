# Keptn Lifecycle Controller

![build](https://img.shields.io/github/workflow/status/keptn/lifecycle-controller/CI)
![goversion](https://img.shields.io/github/go-mod/go-version/keptn/lifecycle-controller?filename=operator%2Fgo.mod)
![version](https://img.shields.io/github/v/release/keptn/lifecycle-controller)
![status](https://img.shields.io/badge/status-not--for--production-red)
[![GitHub Discussions](https://img.shields.io/github/discussions/keptn/lifecycle-controller)](https://github.com/keptn/lifecycle-controller/discussions)

The purpose of this repository is to demonstrate and experiment with
a prototype of a _**Keptn Lifecycle Controller**_.
The goal of this prototype is to introduce a more “cloud-native” approach for pre- and post-deployment, as well as the concept of application health checks.
It is an experimental project, under the umbrella of the [Keptn Application Lifecycle working group](https://github.com/keptn/wg-app-lifecycle).

## Deploy the latest release

**Known Limitations**
* Kubernetes >=1.24 is needed to deploy the Lifecycle Controller
* The Lifecycle Controller is currently not compatible with [vcluster](https://github.com/loft-sh/vcluster)

**Installation**

The lifecycle controller includes a Mutating Webhook which requires TLS certificates to be mounted as a volume in its pod. The certificate creation
is handled automatically by [cert-manager](https://cert-manager.io). To install **cert-manager**, follow their [installation instructions](https://cert-manager.io/docs/installation/).

When *cert-manager* is installed, you can run

<!---x-release-please-start-version-->

```
kubectl apply -f https://github.com/keptn/lifecycle-controller/releases/download/v0.3.0/manifest.yaml
```

<!---x-release-please-end-->

to install the latest release of the lifecycle controller.

The lifecycle controller uses the OpenTelemetry collector to provide a vendor-agnostic implementation of how to receive,
process and export telemetry data. To install it, follow their [installation instructions](https://opentelemetry.io/docs/collector/getting-started/).
We also provide some more information about this in our [observability example](./examples/observability/).

## Goals

The Keptn Lifecycle Controller aims to support Cloud Native teams with:

- Pre-requisite evaluation before deploying workloads and applications
- Finding out when an application (not workload) is ready and working
- Checking the Application Health in a declarative (cloud-native) way
- Standardized way for pre- and post-deployment tasks
- Provide out-of-the-box Observability of the deployment cycle

![](./assets/operator-maturity.jpg)

The Keptn Lifecycle Controller could be seen as a general purpose and declarative [Level 3 operator](https://operatorframework.io/operator-capabilities/) for your Application.
For this reason, the Keptn Lifecycle Controller is agnostic to deployment tools that are used and works with any GitOps solution.

## Documentation

- [How to use the Lifecycle Controller?](./docs/user-guide/Overview.md) - How to define workloads for Keptn
- [Usage Examples](./examples/) - Learn how to use Keptn Lifecycle Controller with basic and real-world examples.
- [User Guide](./docs/user-guide/) - Documentation for those working on projects operated by Keptn Lifecycle Controller
- [Administrator Guide](./docs/admin-guide/) - Documentation for administrators who install and maintain Keptn Lifecycle Controller
- [Architecture](./docs/architecture/) - Keptn Lifecycle Controller architecture, overview of key components, and how it works.
- [Specification](./docs/spec/) - Keptn Lifecycle Controller specifications.
- [Keptn Functions Runtime](./functions-runtime/) - Function Runtime used for defining handlers in Keptn

## License

Please find more information in the [LICENSE](LICENSE) file.

## Contributing

See [CONTRIBUTING.md](./CONTRIBUTING.md).

## Thanks to all the people who have contributed 💜

<a href="https://github.com/keptn/lifecycle-controller/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=keptn/lifecycle-controller" />
</a>

Made with [contrib.rocks](https://contrib.rocks).
