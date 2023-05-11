## Linters requirements

This project uses a set of linters to ensure good code quality.
In order to make proper use of those linters inside an IDE,
the following configuration is required.

### Golangci-lint

Further information can also be found in
the [`golangci-lint` documentation](https://golangci-lint.run/usage/integrations/).

#### Visual Studio Code

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

#### GoLand / IntelliJ requirements

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

### Markdownlint

We are using [markdownlint](https://github.com/DavidAnson/markdownlint) to ensure consistent styling
within our Markdown files.
Specifically we are using [markdownlint-cli](https://github.com/igorshubovych/markdownlint-cli).

We are using `GNU MAKE` to ensure the same functionality locally and within our CI builds.
This should allow easier debugging and problem resolution.

#### Markdownlint execution

To verify that your markdown code conforms to the rules, run the following on your local branch:

```shell
make markdownlint
```

To use the auto-fix option, run:

```shell
make markdownlint-fix
```

#### Markdownlint Configuration

We use the default configuration values for `markdownlint`.

This means:

* [.markdownlint-cli2.yaml](./.markdownlint-cli2.yaml) contains the rule configuration
* [.markdownlintignore](./.markdownlintignore) list files that markdown-lint ignores,  using `.gitignore` conventions

We use the default values, so tools like
[markdownlint for VSCode](https://marketplace.visualstudio.com/items?itemName=DavidAnson.vscode-markdownlint)
can be used without additional configuration.
