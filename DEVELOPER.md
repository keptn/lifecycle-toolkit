# Contributing

## Linters

This project uses a set of linters to ensure a good code quality.
In order to make proper use of those linters inside an IDE, the following configurations are required.

### Visual Studio Code

In Visual Studio Code the [Golang](https://marketplace.visualstudio.com/items?itemName=aldijav.golangwithdidi) extension is required. 
Adding the following lines to the Golang extension configuration file will enable all linters used in this project.

```
"go.lintTool": "golangci-lint",
"go.lintFlags": [
  "--fast",
  "--fix"
]
```

### GoLand / IntelliJ

In GoLand or IntelliJ, the plugin [Go Linter](https://plugins.jetbrains.com/plugin/12496-go-linter) will be required.
Once installed, the configuration can be found in `Settings` >> `Tools` >> `Go Linter` <br>
