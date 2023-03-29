---
title: Keptn's Metrics Operator
linktitle: Metrics Operator
description: Basic understanding of KLT Metrics Operator
weight: 80
cascade:
---

### Metrics Operator

**Metrics Operator** collects, processes, and analyzes metrics data from a variety of sources. Once collected, this data, can be used to generate a variety of reports and dashboards that provide insights into the health and performance of the application and infrastructure.

The Metrics Operator consists of two main components:

* Metrics Controller
* Metrics Adapter

```mermaid
graph TD;
    Metrics-Operator-->Metrics-Adapter;
    Metrics-Operator-->Metrics-Controller
style Metrics-Operator fill:#006bb8,stroke:#fff,stroke-width:px,color:#fff
style Metrics-Adapter fill:#d8e6f4,stroke:#fff,stroke-width:px,color:#006bb8
style Metrics-Controller fill:#d8e6f4,stroke:#fff,stroke-width:px,color:#006bb8
```

**Metrics adapter** is used to expose custom metrics from an application to external monitoring and alerting tools. The adapter exposes custom metrics on a specific endpoint where external monitoring and alerting tools can scrape them. It is an important component of the metrics operator as it allows for the collection and exposure of custom metrics, which can be used to gain insight into the behavior and performance of applications running on a Kubenetes cluster.     

**Metrics controller** is used to fetch metrics from a SLI provider. The controller reconciles a `KeptnMetric` CR and updates it's status with the metric value provided by the selected SLI provider. The steps in which the Controller fetches metrics are given below:
- It first fetches the `KeptnMetric` object to reconcile.
- If the object is not found, it returns and lets Kubernetes handle deleting all associated resources.
- If the object is found, the code checks that if the metric has been updated within the configured interval which is defined in the `Spec.FetchIntervalSeconds`. If not, then it skips reconciling and requeues the request for later.
- If the metric should be reconciled, it fetches the provider defined in the `Spec.Provider.Name` field.
- If the provider is not found, it returns and requeues the request for later.
- If the provider is found, it loads the provider and evaluates the query defined in the `Spec.Query` field.
- If the evaluation is succesful, it stores the fetched value in the status of the `KeptnMetric` object.