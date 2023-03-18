# Contribute to the Keptn documentation

This document provides information about contributing to
the [Keptn Lifecycle Toolkit documentation](https://lifecycle.keptn.sh/docs/),
which is part of the [Keptn](https://keptn.sh) website.

The Keptn Lifecycle Toolkit documentation is authored with
[markdown](https://www.markdownguide.org/basic-syntax/)
and rendered using the Hugo
[Docsy](https://www.docsy.dev/) theme.

We welcome and encourage contributions of all levels.
You can make modifications using the GitHub editor;
this works well for small modifications but,
if you are making significant changes,
you may find it better to fork and clone the repository
and make changes using the text editor or IDE of your choice.
You can also run the Docsy based website locally
to check the rendered documentation
and then push your changes to the repository as a pull request.

If you need help getting started,
feel free to ask for help on the `#help-contributing` or `#keptn-docs` channels on the [Keptn Slack](https://keptn.sh/community/#slack).
We were all new to this once and are happy to help you!

## [Guidelines for contributing](https://github.com/keptn/lifecycle-toolkit/tree/main/docs/content/en/docs/guidelines-for-contributing.md)


## [Building the documentation locally](https://github.com/keptn/lifecycle-toolkit/tree/main/docs/content/en/docs/building-docs-locally.md)


### Building markdown files without Hugo

The Hugo generator described above only renders
the markdown files under the */content/docs* directory.
If you need to render another markdown file
(such as this *CONTRIBUTING.md* file)
to check your formatting, you have the following options:

* If you are using an IDE to author the markdown text,
     use the markdown preview browser for the IDE.
* You can push your changes to GitHub
     and use the GitHub previewer (*View Page*).
* You can install and use the
     [grip](https://github.com/joeyespo/grip/blob/master/README.md) previewer
     to view the rendered content locally.
     When *grip* is installed,
     you can format the specified file locally
     by running the following in its own shell:

     ```console
     grip <file>.md
     ```

     Point your browser at `localhost:6419` to view the formatted file.
     The document updates automatically
     each time you write your changes to disk.

## [Interacting with github](https://github.com/keptn/lifecycle-toolkit/tree/main/docs/content/en/docs/Interacting-with-github.md)

### Developer Certification of Origin (DCO)

Licensing is very important to open source projects. It helps ensure the software continues to be available under the
terms that the author desired.

Keptn uses [Apache License 2.0](https://github.com/keptn/lifecycle-toolkit/blob/main/LICENSE) to strike a balance
between open contributions and allowing you to use the software however you would like to.

The license tells you what rights you have that are provided by the copyright holder.
It is important that the contributor fully understands what rights they are licensing and agrees to them.
Sometimes the copyright holder isn't the contributor, such as when the contributor is doing work on behalf of a company.

To make a good faith effort to ensure these criteria are met,
Keptn requires the Developer Certificate of Origin (DCO) process to be followed.

The DCO is an attestation attached to every contribution made by every developer.
In the commit message of the contribution, the developer simply adds a Signed-off-by statement and
thereby agrees to the DCO, which you can find below or at <http://developercertificate.org/>.

```text
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.

Everyone is permitted to copy and distribute verbatim copies of this
license document, but changing it is not allowed.


Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
    have the right to submit it under the open source license
    indicated in the file; or

(b) The contribution is based upon previous work that, to the best
    of my knowledge, is covered under an appropriate open source
    license and I have the right under that license to submit that
    work with modifications, whether created in whole or in part
    by me, under the same open source license (unless I am
    permitted to submit under a different license), as indicated
    in the file; or

(c) The contribution was provided directly to me by some other
    person who certified (a), (b) or (c) and I have not modified
    it.

(d) I understand and agree that this project and the contribution
    are public and that a record of the contribution (including all
    personal information I submit with it, including my sign-off) is
    maintained indefinitely and may be redistributed consistent with
    this project or the open source license(s) involved.
```

#### DCO Sign-Off Methods

The DCO requires a sign-off message in the following format to appear on each commit in the pull request:

```text
Signed-off-by: Humpty Dumpty <humpty.dumpty@example.com>
```

The DCO text can either be manually added to your commit body,
or you can add either **-s** or **--signoff** to your usual git commit commands.
If you are using the GitHub UI to make a change you can add the sign-off message directly to the commit message
when creating the pull request.
If you forget to add the sign-off you can also amend a previous commit
with the sign-off by running **git commit --amend -s**.
If you've pushed your changes to GitHub already you'll need to force push your branch after this with **git push -f**.

**ATTRIBUTION**:

* <https://probot.github.io/apps/dco/>
* <https://github.com/opensearch-project/common-utils/blob/main/CONTRIBUTING.md>
* <https://code.asam.net/simulation/wiki/-/wikis/docs/project_guidelines/ASAM-DCO?version_id=c510bffb1195dc04deb9db9451112669073f0ba5>
* <https://thesofproject.github.io/latest/contribute/contribute_guidelines.html>

## [Source file structure](https://github.com/keptn/lifecycle-toolkit/tree/main/docs/content/en/docs/source-file-structure.md)

## Guidelines for working on documentation in development versus already released documentation

[This material will be provided when we define the versioning scheme to use]

### Documentation for new features

Most documentation changes should be made to the docs-dev branch,
which means creating a PR in the `lifecycle-toolkit` repository
under the `docs/content/en/docs` directory.
You can view the local build as described above.
We are releasing new versions of the software frequently
so this makes new content available reasonably quickly.

### Documentation for published docs

If a critical problem needs to be solved immediately,
you can modify the documentation source in the sandbox.
In this case, modify the files in the
`keptn-sandbox/lifecycle-toolkit-docs` repository directly.
You can view these changes locally on the `localhost:1314` port.

Note that changes made to the docs in the sandbox
will be overwritten so the same changes should be applied
to the corresponding doc source in the `lifecycle-toolkit` documentation.
