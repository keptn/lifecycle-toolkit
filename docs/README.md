# Choosing where to add documentation

If the change to the docs needs to reflect the next version of KLT, please edit them here, following the instructions
below.
For already existing documentation versions directly edit them
from <https://github.com/keptn-sandbox/lifecycle-toolkit-docs> or from <https://lifecycle.keptn.sh/>.

## Adding documentation to the dev repo

To verify your changes to the dev documentations you can use the makefile:

```shell
cd  lifecycle-toolkit/docs
make clone
make build
make server
```

Note: If the above command is not working try with `sudo` command.

After the server is running on <http://localhost:1314>.
Any modification in the docs folder will be reflected on the server under the dev revision.
You can modify the content in realtime to verify the correct behaviour of links and such.

### Markdown linting

To check your markdown files for linter errors, run the following from the repo root:

```shell
make markdownlint
```

To use the auto-fix option, run:

```shell
make markdownlint-fix
```
