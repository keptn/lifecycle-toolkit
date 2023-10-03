---
title: Keptn Components
linktitle: Components
description: Basic understanding of Keptn Components
weight: 20
---

### Keptn Components

Keptn consists of two main components:

* Keptn Lifecycle Operator, which splits into two separate operators
in Release 0.7.0 and later:
  * Lifecycle-Operator
  * Metrics-Operator
* Keptn Lifecycle Scheduler

```mermaid
graph TD;
    KeptnComponents-->Operators;
    KeptnComponents-->Scheduler
   Operators-->Lifecycle-Operator
   Operators-->Metrics-Operator
style KeptnComponents fill:#006bb8,stroke:#fff,stroke-width:px,color:#fff
style Operators fill:#d8e6f4,stroke:#fff,stroke-width:px,color:#006bb8
style Scheduler fill:#d8e6f4,stroke:#fff,stroke-width:px,color:#006bb8
style Lifecycle-Operator fill:#d8e6f4,stroke:#fff,stroke-width:px,color:#006bb8
style Metrics-Operator fill:#d8e6f4,stroke:#fff,stroke-width:px,color:#006bb8
```
