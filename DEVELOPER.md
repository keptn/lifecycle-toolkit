# Contributing

## Linters

This project uses a set of linters to ensure a good code quality.
In order to make proper use of those linters inside an IDE, the following configurations are required.<br>
Further information can also be found at the `golangci-lint` documentation [here](https://golangci-lint.run/usage/integrations/)

### Visual Studio Code

In Visual Studio Code the [Golang](https://marketplace.visualstudio.com/items?itemName=aldijav.golangwithdidi) extension is required.<br>
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

### GoLand / IntelliJ

In GoLand or IntelliJ, the plugin [Go Linter](https://plugins.jetbrains.com/plugin/12496-go-linter) will be required.<br>
The plugin can bi installed via `Settings` >> `Plugins` >> `Marketplace`, search for `Go Linter` and install it.<br>
Once installed, make sure that the plugin is using the `.golangci.yml` file from the root directory.<br>
The configuration of `Go Linter` can be found in the `Tools` section of the settings.
