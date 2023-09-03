---
title: Migrating to Keptn
description: Notes to help you migrate from Keptn v1 to Keptn
weight: 900
hidechildren: false # this flag hides all sub-pages in the sidebar-multicard.html
---

Keptn uses a different paradigm than that used for Keptn v1
and so migration from Keptn v1 will not be a straight-forward process.
In this section, we will assemble information to help people
who want to move from Keptn v1 as it becomes available.

> **Note**
This section is under development.
Information that is published here has been reviewed for technical accuracy
but the format and content is still evolving.
We hope you will contribute your experiences
and questions that you have.

These instructions mostly assume that you want to utilize
all Keptn features.
Note, however, that you can install and use some functionality
such as Keptn Metrics and Observability
without implementing all Keptn features.

This section currently includes the following topics:

* [Migration strategy](strategy) --
  Overview of the recommended migration strategy

* [Set up metrics and observability](metrics-observe) --
  When you identify the data-sources being used
  and provide Keptn with information about your OpenTelemetry collector,
  Keptn begins to accumulate information that you can view
  for your deployment.

* [Migrate Quality Gates to Keptn Evaluations](evaluations)

* [Migrate CI/CD functionality](cicd)

* [Migrate JES services to Keptn Tasks](jes)

* [Migrate remediation services to Day 2 monitoring](day2)
