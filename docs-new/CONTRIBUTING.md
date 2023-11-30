# Contribute to the Keptn documentation

This document provides information about contributing to
the [Keptn documentation](https://lifecycle.keptn.sh/docs/),
which is part of the [Keptn](https://keptn.sh) website.

The Keptn documentation is authored with
[markdown](https://www.markdownguide.org/basic-syntax/)
and rendered using
[MkDocs](https://www.mkdocs.org/).

We welcome and encourage contributions of all levels.
You can make modifications using the GitHub editor;
this works well for small modifications but,
if you are making significant changes,
you may find it better to fork and clone the repository
and make changes using the text editor or IDE of your choice.
You can also run the website locally
to check the rendered documentation
and then push your changes to the repository as a pull request.

If you need help getting started,
feel free to ask for help on the `#help-contributing` or `#keptn-docs` channels on the [Keptn Slack](https://keptn.sh/community/#slack).
We were all new to this once and are happy to help you!

## Guidelines for Contributing

Please check [Contribution Guidelines](content/en/contribute/docs/contrib-guidelines-docs/_index.md).

## Building the Documentation Locally

To build and deploy the documentation in a docker container, execute

```shell
make serve
```

This will setup a docker container, install all needed dependencies,
build the documentation and serve it.

The URL on which your local documentation website is deployed will be
displayed in the logs.
By default is should be `http://0.0.0.0:8000/`

## Interacting with github

The documentation source is stored on github.com
and you use the standard github facilities to modify it.
Please check [Working with Git](content/en/contribute/general/git/_index.md).

### Developer Certification of Origin (DCO)

All commits must be accompanied by a DCO sign-off.
 See
[DCO](content/en/contribute/general/dco)
for more information.
