---
title: Contribution guidelines for documentation
description: Guidelines for contributing towards Keptn Lifecycle Toolkit
weight: 300
---

The [Contribution Guidelines](../../general/contrib-guidelines-gen) page
contains guidelines that are relevant
for both documentation and software contributions.
This page lists additional guidelines
that are relevant only to documentation

## Guidelines for contributing

* Keep your language clean and crisp.
  We do not have a *Style Guide* for Keptn but the
  [Google developer documentation style guide](https://developers.google.com/style)
  is a good general reference.

* Always build the documentation locally to check the formatting
  and verify that all links are working properly.
  See [Build Documentation Locally](../local-building)
  for details.

* Always run the following to fix most markdown issues in your PR
  and identify issues that can not be fixed automatically:

  ```shell
  make markdownlint-fix
  ```

  See [Markdownlint](../linter-requirements/#markdownlint)
  for details.
