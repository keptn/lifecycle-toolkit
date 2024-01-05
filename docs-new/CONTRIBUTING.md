# Contribute to the Keptn documentation

This is the root of the source code for
the [Keptn documentation](https://lifecycle.keptn.sh/docs/),
which is part of the
[Keptn](https://keptn.sh) website.

The Keptn documentation is authored with
[markdown](https://www.markdownguide.org/basic-syntax/)
and rendered using
[MkDocs](https://www.mkdocs.org/).

We welcome and encourage contributions of all levels.
You can make modifications in various ways:

- Use the GitHub editor;
  this works well for small modifications.
- Use GitHub Codespaces.
   See
  [Codespaces](./docs/contribute/general/codespace.md).
- If you are making significant changes,
  you may find it better to fork and clone the repository
  and make changes using the text editor or IDE of your choice.
  See [Working with git](docs/contribute/general/git/index.md).

  You can run the website locally
  to check the rendered documentation.
  and then push your changes to the repository as a pull request.

See the
[Contributing guide](./docs/contribute/index.md)
for more information about tools and practices to use
when contributing to the Keptn project.

If you need help getting started,
feel free to ask for help on the `#keptn` channel on the [CNCF Slack](https://cloud-native.slack.com).
We were all new to this once and are happy to help you!

## Building the Documentation Locally

To build and deploy the documentation in a container, execute

```shell
make docs-serve
```

This sets up a container, installs all needed dependencies,
builds the documentation, and serves it.

The URL on which your local documentation website is deployed
is displayed in the logs.
By default this should be `http://0.0.0.0:8000/`

For more details, see
[Build documentation locally](docs/contribute/docs/local-building.md).
