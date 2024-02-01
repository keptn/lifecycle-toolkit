---
comments: true
---

# Spell Checker

All PRs that are pushed to a Keptn repository
are run through a spell checker that is based on the
[check-spelling](https://github.com/check-spelling/check-spelling)
GitHub action.

Note, however, that you are still responsible for reading your text carefully.
The tool does not flag if you used a real word that is not the right word,
if you misuse "setup" versus "set up", and so forth.

## Handling Spell Checker errors

If you get a Spell Checker error:

1. Click the `details` link for Spell Checker
   in the checks for your PR.
1. Click on the Summary to get to the job summary.
1. This shows the word(s) that caused Spell Checker to fail.

   - If these are genuine misspellings,
     correct the spelling in your local branch
     or using the GitHub editor
     and push the new commit to resolve the errors.

   - If a word that is flagged is a legitimate word,
     follow the instructions in the report
     to propose adding it to our [dictionary](https://github.com/keptn/lifecycle-toolkit/blob/main/.github/actions/spelling/expect.txt).
     This request will be added to your PR for review
     and, if approved, will be merged when the PR is merged.

   - If your PR includes a file that should not be spell-checked,
     you can add it to the
     [excludes.txt](https://github.com/keptn/lifecycle-toolkit/blob/main/.github/actions/spelling/excludes.txt) file
     as part of your PR.
     It will be reviewed and, if approved,
     merged as part of your PR.

For more details about handling errors that are found, see the
[Check-spelling docs](https://docs.check-spelling.dev/).

## Implementation details

For full technical details about the spell checker, see:

- [check-spelling documentation](https://docs.check-spelling.dev/)
- [check-spelling GitHub repository](https://github.com/check-spelling/check-spelling)

The Keptn spell checker checks both documentation and code for spelling,
based on a set of dictionaries:

- We use general English language and technical terminology
  from dictionaries that are provided and maintained by check-spelling.
  The configuration is specified in files in the
  [.github/actions/spelling](https://github.com/keptn/lifecycle-toolkit/tree/main/.github/actions/spelling)
  directory.
  [expect.txt](https://github.com/keptn/lifecycle-toolkit/tree/main/.github/actions/spelling/expect.txt)
  lists Keptn terms for both documentation and code.
- We also use the specialized technical dictionaries provided by check-spelling specific domains
  or programming languages such as Kubernetes, Go or HTML.
  The dictionaries we use are specified in the
  [.github/workflows/spell-checker.yml](https://github.com/keptn/lifecycle-toolkit/blob/main/.github/workflows/spell-checker.yml)
  file.

Check-spelling supports both American and British spelling
and both are allowed in the Keptn documentation.

Check-spelling provides dictionaries for a number of non-English languages
but we do not currently use those.

Note that Check-spelling does not check for proper capitalization of terms.
All dictionary terms are listed with lowercase letters.
Check-spelling recognizes capitalized versions of these words but,
if capitalized words are listed in a dictionary,
check-spelling rejects non-capitalized forms
which are common in code.
