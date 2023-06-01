---
title: Migration strategy
description: General guidelines for migrating your deployment to KLT
weight: 10
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

> **Note**
This section is under development.
Information that is published here has been reviewed for technical accuracy
but the format and content is still evolving.
We hope you will contribute your experiences
and questions that you have.

No two migrations are alike but these are some general guidelines
for how to approach the project.

1. Create a new Kubernetes cluster that is not running anything else
   and use that to build out your deployment environment.
1. Install and configure the deployment tool(s) you want to use
   to deploy the components of your software.
   You can use different deployment tools for different components.
1. Install
