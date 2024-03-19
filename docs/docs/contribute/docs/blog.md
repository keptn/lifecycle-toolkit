---
comments: true
---

# Blogs

Blog posts are authored in markdown,
submitted and reviewed as PRs,
published as part of the web site,
and built using the same tools and GitHub practices
as other documentation.
However, you must take a few additional steps when writing a blog post.

Keptn uses the
[blog plugin from mkdocs-material](https://squidfunk.github.io/mkdocs-material/setup/setting-up-a-blog/)
to manage the blog infrastructure.
To integrate your blog with the blog plugin,
you must also do the following:

* Create an entry in the `docs/blog/.authors.yml` file
  for each author of the blog.

* The GitHub user ID used to identify yourself
  in the `docs/.authors.yml` file
  must be added to the `.github/actions/spelling/expect.txt` file
  so that the spell checker knows about it.

* Blog source is added to the `docs/blog/posts` directory.
  Individual blogs are not listed in the `mkdocs.yaml` file
  like other documentation.
  Instead, the blog plugin manages integration of all blogs
  with source in the `blogs` directory.

* Additional metadata is required as part of the source files's front matter.

Each of these requirements is discussed on this page.

## Populate docs/blog/.authors.yml file

The value of the `authors:` field in the blog's front matter
links to an entry in the `blogs/authors.yml` file.
The blog plugin uses this information to render author information
on the blog page.

The basic fields that we require are documented here.
For information about additional fields that are available, see the mkdocs
[authors_file](https://squidfunk.github.io/mkdocs-material/plugins/blog/#config.authors_file).

```yaml
authors:
  ...
  <GitHub-UserID>:
    name: <Fullname>
    description: <Role>
    avatar: <avatar-for-picture>
    url: <Github-URL-for-user>
```

### author.yml fields

* **GitHub-UserID** -- This is the user ID used to log into GitHub.
  it serves as a key for the record.
    * **name:** -- Your first and last name
      as it should appear in the posted blog
    * **description:** Your role as it should appear in the posted blog.
      For example, "Keptn Maintainer", "Keptn Contributor", "Keptn User",
      your role in another project or title at a company
      with the name of the project or company..
    * **avatar:** URL for the picture to use in blog posts.
      To use the same picture you use on GitHub,
      open the image in a new tab and use the URL displayed in the address bar.
    * **url:** -- URL for your record on GitHub.

### authors.yml example

```yaml
authors:
  ...
  sampleuser:
    name: Sample User
    description: Senior Software Developer, Example, Inc.
    avatar: https://avatars.githubusercontent.com/u/...
    url: https://github.com/sampleuser
```

## Update spelling/expect.txt file with your ID

The spell checker will flag your user ID as an unrecognized word.
You can manually add this string to the
`.github/actions/spelling/expect.txt` file
as part of your PR,
although the easiest way to handle this is to push your PR
then resolve the error as discussed on the
[Spell Checker](spell-check.md)
page.

## Blog source code

Your blog should be developed using the standard Git
flow documented in
[Working with Git](../general/git/index.md).
When you have created a local branch:

* Create a .md file in the `docs/blog/posts` directory.
  Give the file a meaningful name;
  remember that many people from different organizations
  may be contributing to this directory.

    You do not need to modify the `mkdocs.yml` file for your blog.

* If your blog has graphics, screen shots, YAML files, etc.
  that will be included,
  create an `assets` subdirectory for those files
  in a subdirectory that has the same name
  as the root name of your .md file.
  For example, if your source file is named `myblog.md`
  (which is not actually detailed enough to be a good file name),
  you need to create a myblog/assets subdirectory.

* Follow the instructions in
  [Blog front matter](#blog-front-matter)
  to provide the metadata required by the blog plugin.

* The text that follows the `#` line until the first `##` line
  is the introduction to your article
  and also the abstract that is displayed on the "Blogs" landing page.

Other coding notes for blogs:

* Blog posts are considered part of the MkDocs NAV path.
  This means that:

    * Use the practices documented in
      [External links and internal cross-references](code-docs.md/#external-links-and-internal-cross-references)
      for your blog.
    * You can use a local build to render your blog locally as you write.
      See
      [Build documentation locally](local-building.md)
      for details.
    * The `readthedocs.build` preview associated with your PR
      contains the rendered version of your blog
      so that you and your reviewers can see it.

## Blog front matter

The blog plugin requires some information to manage the blog.
This is specified as part of the file's metadata.
Here we document the fields that are required.
Additional fields can be added; see the
[blog plugin documentation](https://squidfunk.github.io/mkdocs-material/setup/setting-up-a-blog/#writing-your-first-post)
for more information.

```md
---
date: YYYY-MM-DD
authors: [<GitHub-UserID>]
description: >
  <Brief description of this blog, all in one source line>
categories:
  - <cat-1>
  - <cat-2>
  - ...
comments: true
---

# <blog title>
```

### Blog front matter fields

* **date** -- Date when blog was most recently posted
* **authors** -- Author of this blog,
  identified by the Github User ID that is used as the key
  in the `docs/blog/.authors.yml` file.
  This is used to generate the author information
  that is displayed.
* **description: >** -- Brief description
* **categories:** -- Keywords used to generate entries in the "Categories"
  section of the "Blogs" landing page.
  Set as many categories as appropriate and use reasonable terminology.
  If an existing category matches a category you want for your blog,
  be sure to match that terminology exactly.
  For example, "Installation" is an existing category
  so it would be inappropriate to define "Installing" as a category.
* **comments: true**
* **blog-title** -- The title that is displayed for your blog,
  coded as a level 1 (`#`) header.

### Front matter example

```md
---
date: 2024-02-01
authors: [sampleuser]
description: >
  This blog details how to integrate Keptn with MyTool.
categories:
  - SRE
  - Analysis
  - MyTool
comments: true
---

# Using Keptn with MyTool
```
