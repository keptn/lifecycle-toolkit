---
comments: true
---

# Markdownlint

Keptn uses
[markdownlint](https://github.com/DavidAnson/markdownlint)
to ensure consistent styling within all our markdown files,
including files outside the documentation NAV path
such as `README.md` files..
Specifically, we are using
[markdownlint-cli](https://github.com/igorshubovych/markdownlint-cli).

This page tells how to use  markdownlint.

>
We use `GNU make` to ensure the same functionality locally and within our CI builds.
This allows easier debugging and problem resolution.

## Markdownlint execution

To verify that your markdown code conforms to the rules,
run the following on your local branch:

```shell
make markdownlint
```

To use the auto-fix option, run:

```shell
make markdownlint-fix
```

## Markdownlint configuration

We use the default configuration values for `markdownlint`.

This means:

[.markdownlint-cli2.yaml](https://github.com/keptn/lifecycle-toolkit/blob/main/.markdownlint-cli2.yaml)
contains the rule configuration

We use the default values, so tools like
[markdownlint for VSCode](https://marketplace.visualstudio.com/items?itemName=DavidAnson.vscode-markdownlint)
can be used without additional configuration.
