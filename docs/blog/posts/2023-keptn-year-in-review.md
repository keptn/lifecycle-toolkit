---
date: 2024-01-24
authors: [staceypotter]
description: >
  In this blog post we reflect back on what a great year 2023 was for Keptn and the road ahead for this year.
comments: true
---

# 2023 Keptn Year in Review

Hello and Happy New Year from the entire Keptn Team.
We hope you had a wonderful holiday season and that you’re all having a wonderful 2024 so far!
We’d like to take a moment to reflect back on what a great year 2023 was for Keptn and what we're looking forward to in 2024!

The biggest news for the project in 2023 was probably the maturing of the cloud native Keptn,
with the former subproject named Keptn Lifecycle Toolkit officially becoming Keptn in August 2023,
and the end-of-life (EOL) for Keptn v1.
We focused on having Keptn be a 100% cloud native,
opinionated way of delivering Kubernetes apps and it was important for us to implement the
“big pillars” of Keptn v1: supporting the metrics component of the Keptn v1 quality gates feature,
observability (new), and deployments by defining KeptnTasks & Evaluations to execute any container
the user provides and extending deployments (job executor service in Keptn v1).
<!-- more -->

We’re happy to report that Keptn has now reached a general level of maturity and stabilization that
we now have beta APIs.
If you’d like to learn more about the transition from v1 to Keptn, the history
of Keptn, or migrating, we recommend the following articles:

- [September 2023: Keptn Lifecycle Toolkit is now Keptn!](https://medium.com/keptn/keptn-lifecycle-toolkit-is-now-keptn-e0812217bf46)
- [July 2022: Keptn reaches incubating status in CNCF!](https://medium.com/keptn/keptn-reaches-the-incubating-status-in-the-cncf-67291e2dda7)
- [Migrating to Keptn](https://keptn.sh/stable/docs/migrate/keptn/)

## Latest Release

### Keptn v0.9.0

This release introduces a new way of scheduling your Keptn workloads by making use of the new
Kubernetes Scheduling Gates feature on Kubernetes clusters with version v1.27 or greater.
We also revamped our Helm chart which is now split into subcharts per operator.
Please find our new Keptn v0.9.0 release [here](https://github.com/keptn/lifecycle-toolkit/releases/tag/klt-v0.9.0)
and the corresponding Helm chart version 0.3.0 [here](https://github.com/keptn/lifecycle-toolkit-charts/releases/tag/keptn-0.3.0).

More details:
🚀 Helm chart revamp: We restructured the helm charts so that each operator has its own helm chart,
with an umbrella chart implemented.
This supports users who want to install or upgrade just specific
components rather than the entire product while still allowing users to install all of Keptn.
The Keptn Helm chart was renamed from klt to keptn.
Also, labels for some components were changed to better describe them.
This will unfortunately need a re-install of the chart, since upgrades are not possible with changed labels.

🌟Scheduling Gates: With Kubernetes versions v1.27 and greater, this feature allows Keptn to use the
default Kubernetes scheduler rather than needing its own custom scheduler.
You can rely on the default scheduler to do its job and Keptn uses the new scheduling gates for blocking
workloads for pre-deployment tasks.
See the [Kubernetes doc: Pod Scheduling Readiness](https://kubernetes.io/docs/concepts/scheduling-eviction/pod-scheduling-readiness/).

#### Lifecycle-Operator v0.8.3 ([release notes](https://github.com/keptn/lifecycle-toolkit/releases/tag/lifecycle-operator-v0.8.3))

#### Metrics-Operator v0.9.0 ([release notes](https://github.com/keptn/lifecycle-toolkit/releases/tag/metrics-operator-v0.9.0))

## Keptn 2023 Talks & Presentations

<!-- markdownlint-disable MD013 -->
<!-- markdownlint-disable MD039 -->
<!-- markdownlint-disable max-one-sentence-per-line -->
We are so grateful to our community who continue to give talks/presentations about our project.
Let us know if we missed yours, for now here's a collection from 2023:

- Dec 2023: [GitOpsCon Europe (Virtual): Working Your Way Through GitOps with Keptn - Rakshit Gondwal, Cloud Native Computing Foundation & Shivang Shandilya, Open Source](https://youtu.be/CQhXfzYVAwY?feature=shared)
- Nov 2023: KubeCon CloudNativeCon North America (Chicago):
    - [Lightning Talk: From Novice to Keptn Contributor: Empowering My Journey in Cloud Native Communities by Yash Pimple](https://youtu.be/TyZS5mH6vM0?feature=shared)
    - [Keptn Lifecycle Toolkit Updates and Deep Dive by Anna Reale](https://youtu.be/H3UxOwS06iI?feature=shared)
- Oct 2023: [Keptn Metrics Operator by Rakshit Gondwal](https://youtu.be/K9O2Xi8P6Y0?feature=shared&t=549)
- Jul 2023: [Keptn Demo: Python & Container Functions Runtime by Florian Bacher](https://youtu.be/fkuo6CAJ1l8?feature=shared)
- Jun 2023: [How GitHub Codespaces Makes Contributing to Keptn & Keptn Docs Much Easier! By Adam Gardner](https://youtu.be/sFNzOhZw7Eg)
- May 2023: cdCon + GitOpsCon (Vancouver):
    - Comparisons of Open Source GitOps Tooling - Akshay Yadav, Orange Business Services & Monika Yadav, NorthCap University (not recorded)
    - Extending Observability to the Application Lifecycle with ArgoCD, Flux, and Keptn by Ana Margarita Medina & Adam Gardner (not recorded)
- Apr 2023: [Introducing Keptn Lifecycle Toolkit by Andi Grabner](https://youtu.be/449HAFYkUlY?feature=shared)
- Apr 2023: KubeCon CloudNativeCon Europe (Amsterdam): [Navigating the Delivery Lifecycle with Keptn by Giovanni Liva, Ana Medina, Brad McCoy, Meha Bhalodiya](https://youtu.be/Ezd6hGnRL84?feature=shared)
- Mar 2023: [Keptn Lifecycle Toolkit: Installation and KeptnTask Creation in Minutes by Adam Gardner](https://youtu.be/Hh01bBwZ_qM?feature=shared)
- Feb 2023: [Is it Observable: What is Keptn Lifecycle Toolkit with Giovanni Liva](https://youtu.be/Uvg4uG8AbFg?feature=shared)
<!-- markdownlint-enable max-one-sentence-per-line -->
<!-- markdownlint-enable MD039 -->
<!-- markdownlint-enable MD013 -->

## Community

At the end of 2023 we transitioned away from our own Keptn Slack workspace to
[CNCF Slack](http://cloud-native.slack.com/)
utilizing the [#keptn channel](https://cloud-native.slack.com/messages/keptn/) to
better align with the Cloud Native ecosystem & other CNCF projects.
We also saw a [new website](https://keptn.sh/) and [docs](https://keptn.sh/stable/docs/)
(the new [docs folder in GitHub](https://github.com/keptn/lifecycle-toolkit/tree/main/docs)
has been updated as well) and look forward to major improvements to both in 2024.

During the 2023 year, we had 221 PRs Opened, 203 PRs Merged, and 176 Active GitHub Members of those 110 were new members,
and 66 were returning members.
We consider “active members” to be persons who have actively engaged with the Keptn
project through actions such as opened or commented on an Issue, opened a PR, had a PR merged,
commented – replied – or created a discussion, forked a repo.

A very special thanks to these new 2023 contributors:

- [Rakshit Gondwal](https://github.com/rakshitgondwal)
- [Prakriti Mandal](https://github.com/prakrit55)
- [Sudipto Baral](https://github.com/sudiptob2)
- [Geoffrey Israel](https://github.com/geoffrey1330)
- [Yash Pimple](https://github.com/YashPimple)
- [Shivang Shandilya](https://github.com/shivangshandilya)
- [Sambhav Gupta](https://github.com/sambhavgupta0705)
- [Utkarsh Umre](https://github.com/UtkarshUmre)
- [Victor Anene](https://github.com/Vickysomtee)
- [Arya Soni](https://github.com/aryasoni98)

## Google Summer of Code

[Google summer of Code (GSoC)](https://summerofcode.withgoogle.com/) is a global, online program focused on
bringing new contributors into open source software development.
GSoC Contributors work with an open source organization on a 12+ week programming project under the guidance of mentors.

We were thrilled to participate in Google Summer of Code last year.
Our primary Mentor this year was Florian Bacher, and our GSoC Contributor was Rakshit Gondwal.
You can read the project proposal here:
[Timeframe for Metrics proposal](https://summerofcode.withgoogle.com/archive/2023/projects/e7z3n3kH).
All in all, we were quite pleased with the results, and because of Rakshit’s continued contributions upon GSoC
completion we were happy to welcome him on as an official Keptn Approver in August 2023.
Rakshit’s GSoC project summary can be found [here](https://github.com/rakshitgondwal/gsoc-2023)
if you’d like to read about his experience.

<!-- markdownlint-disable MD028 -->
> Working with Keptn and my mentors was a great experience.
> Right from the start, I felt welcomed and supported by the community.
> It's amazing how everyone is ready to help out.
> If my main mentor wasn't available, there was always someone else who stepped in to guide me through any issues I faced.
> What I appreciated a lot was that I never felt stressed during the coding phase.
> This made a big difference in helping me learn and work better.
>
> *- Rakshit Gondwal, Google Summer of Code Contributor & Keptn Approver -*

> Working alongside Rakshit during GSoC to enhance the metrics capabilities of Keptn has been a great experience.
> Rakshit's commitment to the project, problem-solving skills, and willingness to learn made it a blast to work with
> him on this project.
> His contributions have significantly improved Keptn's functionality, and I anticipate lasting benefits from Rakshit's work.
>
> *- Florian Bacher, Rakshit’s GSoC Mentor and Keptn Maintainer -*
<!-- markdownlint-enable MD028 -->

## OpenTelemetry CI/CD Working Group

In Keptn, we consider observability a first-class citizen and we want to support users by helping to make
their Kubernetes deployments efficient.
Keptn helps compute DORA metrics by tracking the application from the moment it hits a development cluster
until it hits production with [KEP87](https://github.com/keptn/enhancement-proposals/blob/main/text/0087-klt-traceid-propagation.md).
However, this is only a portion of the full picture, and we would like to support the community with an
Open Standard that allows different tools to work together tracking code from the moment it is merged into main.
For this reason, we support the new
[OTel CI/CD working group](https://github.com/open-telemetry/community/pull/1822/files).

## The Road Ahead in 2024

We’re looking forward to Context Propagation (see [GitHub Issue 1394](https://github.com/keptn/lifecycle-toolkit/issues/1394))
which will enable users to have full traceability of their application's deployment across multiple stages using OpenTelemetry,
full GitOps benefits, and enable users to act upon complex metadata, including attributes such as the current name and version
of the deployed workload/application, as well as the OTel traceParent of the current deployment phase.
Be on the lookout for documentation guides, how-tos, blogs, and a demo for KubeCon EU (Paris) in March!

We started hosting New Contributor Meetings in 2023 as well as rebooted our Keptn Online User Group Meetings this year
and are kicking off the season with
[An Introduction to Keptn on January 24 at 10am ET / 4pm CET by Florian Bacher](https://community.cncf.io/events/details/cncf-keptn-community-presents-keptn-online-user-group-meeting-an-introduction-to-keptn/).
Join our [CNCF Keptn Community Chapter](https://community.cncf.io/keptn-community/) to get notified of all upcoming events!

As always, thank you for being part of the Keptn Community, and be sure to join us in all of these places online:

<!-- markdownlint-disable MD013 -->
- [CNCF slack](http://cloud-native.slack.com/) [#keptn channel](https://cloud-native.slack.com/messages/keptn/) (if you’re not already a member you can [invite yourself here](https://communityinviter.com/apps/cloud-native/cncf))
- [CNCF Keptn Community Chapter](https://community.cncf.io/keptn-community/) (event listing)
- [Google Calendar](https://calendar.google.com/calendar/u/0?cid=ZHluYXRyYWNlLmNvbV9hYmpyaDF1a2YxOGloNDc3dGIxZWthZzJhZ0Bncm91cC5jYWxlbmRhci5nb29nbGUuY29t)
- [YouTube Channel](https://www.youtube.com/channel/UCHMn9HyAMeb81bRlaOuZyuQ)
- [Twitter/X](https://twitter.com/keptnProject)
<!-- markdownlint-enable MD013 -->
