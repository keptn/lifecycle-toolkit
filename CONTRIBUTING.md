# Contributing at Lifecycle-toolkit

We are thrilled to have you join us as a contributor! The Keptn Lifecycle Toolkit is a community-driven project and we greatly value collaboration. There are various ways to contribute to the Lifecycle Toolkit, and all contributions are highly valued. Please, explore the options below to learn more about how you can contribute.

# How to Contribute
## Prerequisites
Make sure you have the following prerequisites installed on your operating system before you start contributing:

- [Docker](https://docs.docker.com/get-docker/) to build a new version of the containers.
- A Kubernetes cluster >= Kubernetes 1.24.If you don’t have one, we recommend Kubernetes-in-Docker(kind) to set up your local development environment.
- [kubectl](https://kubernetes.io/docs/tasks/tools/) installed on your system.
- [kustomize](https://kustomize.io/) for customizing Kubernetes resource configurations and generating manifests.

### Visual Studio Code requirements:

This project uses a set of linters to ensure good code quality.
In order to make proper use of those linters inside an IDE, the following configuration is required.

- Install the [golangci-lint](https://golangci-lint.run/usage/integrations/) extension to check for syntax, style, performance, and security issues in Go.
- Install [Golang](https://marketplace.visualstudio.com/items?itemName=aldijav.golangwithdidi) extension is required.

Adding the following lines to the `Golang` extension configuration file will enable all linters used in this project.

```
"go.lintTool": {
    "type": "string",
    "default": "golangci-lint",
    "description": "GolangGCI Linter",
    "scope": "resource",
    "enum": [
        "golangci-lint",
    ]
},
"go.lintFlags": {
    "type": "array",
    "items": {
        "type": "string"
    },
    "default": ["--fast", "--fix"],
    "description": "Flags to pass to GCI Linter",
    "scope": "resource"
},
```
Further information can also be found in the [`golangci-lint` documentation](https://golangci-lint.run/usage/integrations/).

<!-- - GoLand / IntelliJ -->
- Install GoLand or IntelliJ, which is an Integrated Development Environment (IDE) for the Go programming language, along with the plugin [Go Linter](https://plugins.jetbrains.com/plugin/12496-go-linter).

The plugin can be installed via `Settings` >> `Plugins` >> `Marketplace`, search for `Go Linter` and install it.
Once installed, make sure that the plugin is using the `.golangci.yml` file from the root directory.

The configuration of `Go Linter` can be found in the `Tools` section of the settings.
