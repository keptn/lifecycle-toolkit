<!-- markdownlint-disable first-line-heading -->
## Deployment Check Orchestration

### orchestrating deployment checks as part of scheduler

To reduce complexity of custom checks use Keptn to:

* Pre-Deploy:
  * validate external dependenicies
  * confiorm images are scanned
  * ...
* Post-Deploy:
  * Execute tests
  * Notify Stakeholders
  * Promote to next stage
  * ...
* Automatically validate against your SLO (Service Level Objectives)
