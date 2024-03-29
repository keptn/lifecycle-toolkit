---
comments: true
---

# Contribution guidelines for documentation

The [Contribution Guidelines](../general/contrib-guidelines-gen.md) page
contains guidelines that are relevant
for both documentation and software contributions.

This page lists additional guidelines
that are relevant only to documentation.

## Guidelines for contributing

* Keep your language clean and crisp.
  We do not have an official *Style Guide* for Keptn but the
  [Google developer documentation style guide](https://developers.google.com/style)
  is a good general reference.

* Use topic sentences for sections and paragraphs.
  When reading a well-written technical document,
  you should be able to read the first sentence in each paragraph
  and especially in each section to get an idea of what might follow.

    Good oral presentations commonly begin with a "set-up"
    where they describe a problem
    and then proceed to tell how to fix that problem.
    When using oral presentations as source material,
    it is important to rewrite the text
    so that the actual subject of discussion comes first.

* Avoid using FAQ's in documentation.
  In general, they say "here is some miscellaneous information
  that I was too lazy to organize logically for you."
  On rare occasions, they may be appropriate,
  such as if you need a quick reference to a large, complicated document
  and include links to detailed documentation about the topic.

* We are attempting to avoid duplicated information across the doc set,
  especially for information that is expected to change.
  For example, information about supported Kubernetes versions
  and the command sequence to install Keptn should usually be done
  as references to the official installation section of the docs.

    For usability considerations, we have a few exceptions to this rule
    for the main `README.md` file and the [Getting Started Guide](../../getting-started/index.md).

* When you want to display a sample file that exists in the repository,
  use the `include <file-path>` shortcode syntax
  (which automatically pulls the current version of the file into your document)
  rather than copying the text.
  See
  [Displaying sample files](code-docs.md/#comments)
  for details.

* `markdownlint` enforces limits on line length.
  Links to other documents are exempted from this limit
  but, if a line has words before and after the long string,
  `markdownlint` fails.
  A good practice is to just put all links on their own line.
  So, instead of coding:

    ```md
    The [Other section](long-link-to-section) page
    ```

    you should code the following,
    unless the link is so short
    that you are sure it will not violate the line-length rules:

    ```md
    The
    [Other section](long-link-to-section)
    page
    ```

* Always build the documentation locally to check the formatting
  and verify that all links are working properly.
  See [Build Documentation Locally](./local-building.md)
  for details.

* Always run the following to fix most markdown issues in your PR
  and identify issues that can not be fixed automatically:

    ```shell
    make markdownlint-fix
    ```

  See [markdownlint](./markdownlint.md)
  for details.
