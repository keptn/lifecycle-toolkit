# Contributing to the Keptn Lifecycle Toolkit

We are thrilled to have you join us as a contributor!
The Keptn Lifecycle Toolkit is a community-driven project and we greatly value collaboration.
There are various ways to contribute to the Lifecycle Toolkit, and all contributions are highly valued.
Please, explore the options below to learn more about how you can contribute.

## Prerequisites

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
