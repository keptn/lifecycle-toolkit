---
title: DCO
description: Working with DCO
weight: 200
---

### Developer Certification of Origin (DCO)

Keptn requires the Developer Certificate of Origin (DCO) process
to be followed for each commit.
With the DCO, you attest that the contribution adheres
to the terms of the Apache License that covers Keptn
and that you are granting ownership of your work to the Keptn project.

Licensing is very important to open source projects.
It helps ensure that the software continues to be available under the
terms that the author desired.
Keptn uses [Apache License 2.0](https://github.com/keptn/lifecycle-toolkit/blob/main/LICENSE),
which strikes a balance between open contributions
and allowing you to use the software however you would like to.

The license tells you what rights the copyright holder gives you.
It is important that the contributor fully understands
what rights they are licensing and agrees to them.
Sometimes the copyright holder is not the contributor,
such as when the contributor is doing work on behalf of a company.

You must add a Signed-off-by statement for each commit you contribute,
thereby agreeing to the DCO.
The text of the DCO is given here
and available at <http://developercertificate.org/>:

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

## DCO Sign-Off Methods

The DCO check runs on each PR to verify
that the commit has been signed off properly.
Your builds will fail and can not be merged if the DCO check fails.

Do any of the following
to implement the DCO signoff on each commit:

* [Add **-s** or **--signoff**](#sign-off-with-git-commit--s)
  to your usual `git commit` commands
* [Manually add text](#manually-add-text-to-commit-description)
  to your commit body
* [Automate DCO](#automate-dco)

## Sign off with git commit -s

Use the **-s** or **--signoff** flag to the `git commit` command
to sign off on a commit.
For example:

```bash
git commit -s -m "my awesome contribution"
```

If you forget to add the sign-off,
run the following command to amend the previous commit
with the sign-off:

```bash
git commit --amend --signoff
```

Use the following command
to sign off the last 2 commits you made:

```bash
git rebase HEAD~2 --signoff
```

## Manually add text to commit description

To sign off on a commit not made with the command line
(such as those made directly with the github editor
or as suggestions made to a PR during review),
you can add text like the following to the commit description block.
You must specify your real name and the email address to use:

  ```text
  Signed-off-by: Humpty Dumpty <humpty.dumpty@example.com>
  ```

## Automate DCO

The DCO  process is sometimes inconvenient but you can automate it
by creating a pre-commit git hook as follows:

1. Create the hook:

    ``` bash
    touch .git/hooks/prepare-commit-msg
    ```

2. Add the following to the `prepare-commit-msg` file:

    ```bash
    SOB=$(git var GIT_AUTHOR_IDENT | sed -n 's/^\(.*>\).*$/Signed-off-by: \1/p')
    grep -qs "^$SOB" "$1" || echo "$SOB" >> "$1"
    ```

3. Give it execution permissions by calling:

    ```bash
    chmod +x ./.git/hooks/prepare-commit-msg
    ```

## ATTRIBUTION

* <https://probot.github.io/apps/dco/>
* <https://github.com/opensearch-project/common-utils/blob/main/CONTRIBUTING.md>
* <https://code.asam.net/simulation/wiki/-/wikis/docs/project_guidelines/ASAM-DCO?version_id=c510bffb1195dc04deb9db9451112669073f0ba5>
* <https://thesofproject.github.io/latest/contribute/contribute_guidelines.html>
