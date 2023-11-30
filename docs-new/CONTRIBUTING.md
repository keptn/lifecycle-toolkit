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

## Building the Documentation Locally

To build and deploy the documentation in a docker container, execute

```shell
make docs-serve
```

This will setup a docker container, install all needed dependencies,
build the documentation and serve it.

The URL on which your local documentation website is deployed will be
displayed in the logs.
By default is should be `http://0.0.0.0:8000/`
