## Choosing where to add documentation

If the change to the docs needs to reflect the next version of KLT, please edit them here, following the instructions below.
For already existing documentation versions directly edit them from https://github.com/keptn-sandbox/lifecycle-toolkit-docs or from https://lifecycle.keptn.sh/.

## Adding documentation to the dev repo

To verify your changes to the dev documentations you can use the makefile:

```
cd  lifecycle-toolkit/docs

make clone
make build
make server
```

After the server is running on http://localhost:1314/docs-dev.
Any modification in the docs folder will be reflected on the server under the dev revision.
You can modify the content in realtime to verify the correct behaviour of links and such.

### Markdown linting
To check your markdown files for linter errors, run the following from the repo root:

```
make markdownlint
```

To use the auto-fix option, run:

```
make markdownlint-fix
```
