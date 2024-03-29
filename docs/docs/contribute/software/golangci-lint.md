---
comments: true
---

# Golangci-lint configuration

Kept uses the `Golangci-lint` linter to ensure good code quality.
This page describes the configuration required
to make proper use of those linters inside an IDE.

Further information can be found in
the [`golangci-lint` documentation](https://golangci-lint.run/welcome/integrations/).

## Visual Studio Code

In Visual Studio Code the
[Golang](https://marketplace.visualstudio.com/items?itemName=aldijav.golangwithdidi)
extension is required.

Adding the following lines to the `Golang` extension
configuration file enables all linters used in this project.

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

## GoLand / IntelliJ requirements

* Install either the **GoLand** or **IntelliJ**  Integrated Development Environment
(IDE) for the Go programming language, plus the [Go Linter](https://plugins.jetbrains.com/plugin/12496-go-linter) plugin.

* The plugin can be installed via `Settings` >> `Plugins` >> `Marketplace`,
search for `Go Linter` and install it.
Once installed, make sure that the plugin is using the `.golangci.yml`
file from the root directory.

* The configuration of `Go Linter` can be found in the `Tools` section
of the settings.

If you are on Windows, you need to install **make** for the above process to complete.

> **Note**
When using the make command on Windows, you may receive an `unrecognized command` error for a command that is installed.
This usually indicates that `PATH` for the binary is not set correctly).
