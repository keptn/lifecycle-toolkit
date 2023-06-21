---
title: Lifecycle Toolkit Components
linktitle: Components
description: Basic understanding of Keptn Lifecycle Toolkit Components
weight: 80
cascade:
---

### Keptn Lifecycle Toolkit Components

The Keptn Lifecycle Toolkit consists of two main components:

* Keptn Lifecycle Operator, which splits into two separate operators
in Release 0.7.0 and later:
  * Lifecycle-Operator
  * Metrics-Operator
* Keptn Lifecycle Scheduler

```mermaid
graph TD;
    KLTComponents-->Operators;
    KLTComponents-->Scheduler
   Operators-->Lifecycle-Operator
   Operators-->Metrics-Operator
style KLTComponents fill:#006bb8,stroke:#fff,stroke-width:px,color:#fff
style Operators fill:#d8e6f4,stroke:#fff,stroke-width:px,color:#006bb8
style Scheduler fill:#d8e6f4,stroke:#fff,stroke-width:px,color:#006bb8
style Lifecycle-Operator fill:#d8e6f4,stroke:#fff,stroke-width:px,color:#006bb8
style Metrics-Operator fill:#d8e6f4,stroke:#fff,stroke-width:px,color:#006bb8
```
