---
comments: true
---

# Markdownlint

We are using [markdownlint](https://github.com/DavidAnson/markdownlint) to ensure consistent styling
within our Markdown files.
Specifically we are using [markdownlint-cli](https://github.com/igorshubovych/markdownlint-cli).

>
We are using `GNU make` to ensure the same functionality locally and within our CI builds.
This should allow easier debugging and problem resolution.

## Markdownlint execution

To verify that your markdown code conforms to the rules, run the following on your local branch:

```shell
make markdownlint
```

To use the auto-fix option, run:

```shell
make markdownlint-fix
```

## Markdownlint Configuration

We use the default configuration values for `markdownlint`.

This means:

[.markdownlint-cli2.yaml](https://github.com/keptn/lifecycle-toolkit/blob/main/.markdownlint-cli2.yaml)
contains the rule configuration

We use the default values, so tools like
[markdownlint for VSCode](https://marketplace.visualstudio.com/items?itemName=DavidAnson.vscode-markdownlint)
can be used without additional configuration.
