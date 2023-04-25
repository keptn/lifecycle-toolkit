---
title: Requirements
description: Supported software versions and information about resources required
icon: concepts
layout: quickstart
weight: 15
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

## Supported Kubernetes versions

The Keptn Lifecycle Controller requires Kubernetes v1.24.0 or later.

## Resource requirements

## cert-manager

KLT includes a lightweight cert-manager
that is used for installation and Webhooks.
You can configure a different cert-manager
before you install KLT.
See [Implement your own cert-manager](cert-manager.md)
for instructions.
