---
comments: true
---

# Coding the docs

Keptn documentation is written using the Markdown language,
with each page written in a separate file.
The following documents document the language:

* [Markdown Guide](https://www.markdownguide.org/getting-started/#flavors-of-markdown)
  discusses Markdown structure and background.
* [Basic Syntax](https://www.markdownguide.org/basic-syntax/#links)
  summarizes the standard Markdown syntax
  that is support by almost all Markdown variants.
* [Markdown Cheat Sheet](https://www.markdownguide.org/cheat-sheet/)
  is a handy reference for the most commonly used Markdown elements.

Markdown supports many variants and the build tools we use
impose a few special requirements that are discussed here.

## Top of source requirements

The top of each documentation source file should look like:

```markdown
---
comments: true
---

# Coding the docs

Beginning of information about the topic.

```

The elements are:

* The `comments` block.
  This allows readers to post comments to the published page.

* A level 1 header (`# title`)  with the title of the page
  as it is displayed in the main canvas of the docs..
  This must be preceded and followed by a single blank line.

    The title displayed in the left frame
    is determined by the title in the `mkdocs.yml` file.
    Be sure that these two titles match.

* Text that introduces the information for the page.
  Do not use stacked headers, with a level 2 header (`## title`)
  immediately following the level 1 header.

## Comments

To comment a line in the documentation,
prepend the `<!--` string at the beginning of the line,
as in:

```markdown
<!-- This is a comment
```

## Indentation of nested lists and code blocks

Paragraphs and code blocks that are nested under a list item
must be indented two spaces from the text of the list item.
If they are not,
the indented material is rendered as flush-left
and ordered lists do not increment the list item number correctly.

For example, the formatting of the bullet list in the preceding section is:

```markdown
* The `comments` block.
  This allows readers to post comments to the published page.

* A level 1 header (`# title`)  with the title of the page
  as it is displayed in the main canvas of the docs..
  This must be preceded and followed by a single blank line.

    The title displayed in the left frame
    is determined by the title in the `mkdocs.yml` file.
    Be sure that these two titles match.

* Text that introduces the information for the page.
  Do not use stacked headers, with a level 2 header (`## title`)
  immediately following the level 1 header.
```

## External links and internal cross-references

Use the standard Markdown conventions for links:

```markdown
[display-string](target-link)
```

The syntax of the `target-link` is different
for external links and internal documentation cross-references.

We recommend putting the link code on a separate line in the source code.
The markdownlint tool limits the number of characters on a line.
Links are exempt from this check
but markdownlint fails the line if it includes text before or after the link.
This is not absolutely necessary if the link target is short
but this convention prevents problems.

### External links to and from documentation

Links to and from the documentation set
from outside the `NAV` path defined in the `mkdocs.yml` file
use the full URL as displayed in the browser address bar
(including the trailing `/` character)
for the page for the `target-link`.

This syntax is used for:

* Links from a documentation page to an external page
* Links **to** files in the same repository as the docmentation source
  but outside the documentation `NAV` path
* Links **from** files in the same repository as the documentation source
  but outside the documentation `NAV` path,
  such as `README.md` and `CONTRIBUTING.md` files

    Links using a relative path to files outside the `NAV` path
    resolve correctly but the targeted documentation page
    does not include the contents block in the left frame.

An example of the coding for an external link is:

```markdown
The Kubernetes
[Pod](https://kubernetes.io/docs/concepts/workloads/pods/)
documentation
```

### Internal cross references in the documentation set

Internal cross-references between pages in the documentation set
(which is the documentation `NAV` path as defined in the `mkdocs.yml` file)
use a `target-link` that is a modified version
of the URL displayed for the page in the rendered documentation.

We suggest that you copy/paste the portion of the URL
that follows `docs/docs` as the base for your `target-link`.
You must then make the following modifications:

* Specify the path name of the targeted file
  relative to `docs/docs` directory
  using the shell convention where `../` represents
  the parent directory
* Add the `.md` suffix to the file name
* Remove the trailing / from the string
* When referencing a sub-section of a page,
  remove the `/` character between the page tag
  and the `#` character that tags the referenced subsection.
* When referencing a section of the docs,
  add the `index.md` or `_index.md` filename to the path

Some examples may clarify this.

#### Cross reference a file in another directory

The full URL for the `Analysis` CRD reference page is:

```shell
https://keptn.sh/stable/docs/reference/crd-reference/analysis/
```

To cross-reference this page
from any page in the `docs/guide` directory
(or other pages at that level), the code is:

```shell
See the
[Analysis](../reference/crd-reference/analysis.md)
CRD reference page.
```

To form this cross-reference::

* Copy/paste the part of the URL after `docs` as a base
* Insert `../` to go up one directory from `guides` to `docs`,
  before the path that goes down the `reference/crd-reference` path
  to identify the file
* Add the `.md` suffix to `Analysis` to form the actual source file name.
* Remove the trailing `/` of the URL

#### Cross-reference a sub-section of another page

To get a link to the `Examples` subsection of the `Analysis` reference page,
view the page in your browser and select `Examples`
from the contents listing in the right frame.
This gives you the following URL:

```shell
https://keptn.sh/stable/docs/reference/crd-reference/analysis/#examples
```

To link to that sub-section, the code is:

```shell
See
[Examples](../reference/crd-reference/analysis.md#examples
```

You see that the `/` in the URL before `#examples` has been removed.

#### Cross-reference another file in the same directory

Another CRD reference page (which is in the same directory)
can reference the `Analysis` reference page
like this:

```shell
[Analysis](analysis.md)
```

#### Cross-reference another section

The URL of the `Installation` section is:

```shell
https://keptn.sh/stable/docs/installation/
```

To cross-reference this section from a file in the `guides` section
(or other file at that level),
use the relative file to the directory
and specify the `index.md` file for the section:

```shell
Follow the instructions in the
[Installation](installation/index.md)
section.
```
