---
comments: true
---

# CRD name

Copy this template to create a new CRD reference page.

1. Replace the variable text in metadata with information for this page
1. Delete these instructions from your file
1. Populate the page with appropriate content

## Synopsis

```yaml
apiVersion: <library>
kind: <kind>
metadata:
  name: <name>
spec:
  ...
```

## Fields

<!-- Detailed description of each field/parameter -->

<!-- * **apiVersion** -- API version being used -->
<!-- * **kind** -- Resource type. -->
<!--    Must be set to `<xxx>` -->
<!-- * **metadata** -->
<!--   * **name** -- Unique name of this <resource>. -->
<!--     Names must comply with the -->
<!-- markdownlint-disable-next-line line-length -->
<!--     [Kubernetes Object Names and IDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#dns-subdomain-names) -->
<!--     specification. -->
<!-- * **spec** -->
<!-- ... -->

Fields should be clearly marked as to whether they are
required or optional.
If they are optional,
text should explain the behavior when that field is not populated.

## Usage

<!-- How this CRD is "activated".  For example, which event uses this CRD -->
<!-- Instructions and guidelines for when and how to customize a CRD -->

## Examples

Include code snippets that illustrate
how this resource is populated.
Code examples should use be embedded links to example source code
so that they will be updated automatically when changes are made to the example.

## Files

* Link to subsection for this resource in the "API Reference"

## Differences between versions

## See also

* Links to related reference pages
* Links to related User Guides and other documentation
