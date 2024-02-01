---
comments: true
---

# Components

## Keptn Components

Keptn consists of the following main components:

* [Keptn Lifecycle Operator](./lifecycle-operator/index.md)
* [Keptn Metrics Operator](./metrics-operator.md)
* [Keptn Scheduler](./scheduling.md)
* [Keptn Certificate Manager](./certificate-operator.md)

The architectural diagram:

```mermaid
graph TD;

A[Lifecycle Operator]
B[Metrics Operator] -- provide metrics --> A
C[Cert manager] -- watch certificate --> A
C[Cert manager] -- watch certificate --> B

style A fill:#d8e6f4,stroke:#fff,stroke-width:px,color:#006bb8
style B fill:#d8e6f4,stroke:#fff,stroke-width:px,color:#006bb8
style C fill:#d8e6f4,stroke:#fff,stroke-width:px,color:#006bb8
```
