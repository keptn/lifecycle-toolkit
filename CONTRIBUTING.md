# Contributing to the Keptn Lifecycle Toolkit

We are thrilled to have you join us as a contributor!
The Keptn Lifecycle Toolkit is a community-driven project and we greatly value collaboration.
There are various ways to contribute to the Lifecycle Toolkit, and all contributions are highly valued.
Please, explore the options below to learn more about how you can contribute.

## Prerequisites

## Related Technologies

You should understand some related technologies
to effectively use and contribute to the Keptn lifecycle-toolkit.
This section provides links to some materials that can help your learning.
The information has been gathered from the community and is subject to alteration.
If you have suggestions about additional content that should be included in this list, please submit an issue.

### Kubernetes

- **Understand the basics of Kubernetes**
  - [ ] [Kubernetes official documentation](https://kubernetes.io/docs/concepts/overview/)
  - [ ] [Kubernetes For Beginner](https://youtu.be/X48VuDVv0do)
- **Kubernetes Architecture**
  - [ ] [Philosophy](https://youtu.be/ZuIQurh_kDk)
  - [ ] [Kubernetes Deconstructed: Understanding Kubernetes by Breaking It Down](https://www.youtube.com/watch?v=90kZRyPcRZw)
- **CRD**
  - [ ] [Custom Resouce Definition (CRD)](https://www.youtube.com/watch?v=xGafiZEX0YA)
  - [ ] [Kubernetes Operator simply explained in 10 mins](https://www.youtube.com/watch?v=ha3LjlD6g7g)
  - [ ] [Writing Kubernetes Controllers for CRDs](https://www.youtube.com/watch?v=7wdUa4Ulwxg)
- **Kube-builder Tutorial**
  - [ ] [book.kubebuilder.io](https://book.kubebuilder.io/introduction.html)
- **Isitobservable**
  - [ ] Keptn has tight integrations with Observability tools and therefore knowing how to _Observe a System_ is important.
  - [ ] [Isitobservable website](https://isitobservable.io/)
  - [ ] [Is it Observable? with Henrik Rexed](https://www.youtube.com/watch?v=aMwk2qo0v40)

### Understanding SLO, SLA, SLIs

- **Overview**
  - [ ] [Overview](https://www.youtube.com/watch?v=tEylFyxbDLE)
  - [ ] [The Art of SLOs (Service Level Objectives)](https://www.youtube.com/watch?v=E3ReKuJ8ewA)

### Operator SDK

- **Go-based Operators**
  - [ ] [Go operator tutorial from RedHat](https://docs.okd.io/latest/operators/operator_sdk/golang/osdk-golang-tutorial.html)

## Linters

This project uses a set of linters to ensure good code quality.
In order to make proper use of those linters inside an IDE, the following configuration is required.

Further information can also be found in
the [`golangci-lint` documentation](https://golangci-lint.run/usage/integrations/).

### Visual Studio Code

In Visual Studio Code the [Golang](https://marketplace.visualstudio.com/items?itemName=aldijav.golangwithdidi)
extension is required.

Adding the following lines to the `Golang` extension configuration file will enable all linters used in this project.

```json
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

### GoLand / IntelliJ

In GoLand or IntelliJ, the plugin [Go Linter](https://plugins.jetbrains.com/plugin/12496-go-linter) will be required.

The plugin can be installed via `Settings` >> `Plugins` >> `Marketplace`, search for `Go Linter` and install it.
Once installed, make sure that the plugin is using the `.golangci.yml` file from the root directory.

The configuration of `Go Linter` can be found in the `Tools` section of the settings.
