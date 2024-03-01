---
date: 2023-12-19
authors: [bacherfl]
description: >
  In this blog post you will learn how to use Keptn for analyzing the performance of an application.
categories:
  - SRE
  - Application Performance
  - Analysis
comments: true
---

# Application Performance Analysis with Keptn

In the dynamic world of DevOps and continuous delivery, keeping applications reliable
and high-performing is a top priority.

Site reliability engineers (SREs) rely on Service Level Objectives (SLOs) to set the standards that the
Service Level Indicators (SLIs) of an application must meet, like response time, error rate,
or any other metric that might be relevant to the application.

The use of SLOs is not a new concept, but integrating them into an application comes with its own set of issues:

- Figuring out which SLIs and SLOs to use— do you get the SLI values from one monitoring source or multiple?
This complexity makes it harder to use them effectively.
- Defining SLO priorities.
Imagine a new version of a service that fixes a concurrency problem but
slows down response time.
In this case, this may be a valid trade-off and the new version should
not be denied due to an increase in the response time, given that the error rate will decrease.
Situations like these call for a way of defining a grading logic
where different priorities can be assigned to SLOs.

- Defining and storing SLOs.
It's crucial to clearly define and store these goals in one central place,
ideally a declarative resource in a GitOps repository, where each change can be easily traced back.

In this article, we'll explore how Keptn tackles these challenges with its new Analysis feature.
We will deploy a demo application onto a Kubernetes cluster to show Keptn helps SREs gather
and make sense of SLOs, making the whole process more straightforward and efficient.

The example application will provide some metrics by itself by serving them via its Prometheus endpoint,
while other data will come from Dynatrace.
<!-- more -->

## Defining data providers

Everything in Keptn is configured via Kubernetes Custom Resources.
We notify Keptn about our monitoring data sources by adding two KeptnMetricsProvider
resources to our Kubernetes cluster - one for our Prometheus instance, the other one for our Dynatrace tenant.

```yaml
{% include "./application-performance-analysis/metric-providers.yaml" %}
```

## Defining SLIs

Now that we have defined our data sources, let's tell Keptn what SLIs we want to monitor
and how to retrieve them from Prometheus and Dynatrace.
This is done by applying `AnalysisValueTemplate` resources to the cluster.
If you have worked with Keptn in the past, you will notice that the structure of these
resources is similar to the `KeptnMetrics` resources
(see [this article](https://www.linkedin.com/pulse/scaling-kubernetes-workloads-based-dynatrace-metrics-keptnproject),
if you would like to learn more about KeptnMetrics and how to use them to automatically scale your workloads).

The difference between `KeptnMetrics` and `AnalysisValueTemplates` is:

- `KeptnMetrics` are monitored and updated continuously, meaning that they always represent the
latest known value of the given metric.
This makes them a good candidate for being observed by a `HorizontalPodAutoscaler` to make scaling decisions.

- `AnalysisValueTemplates` provide the means to get the value of a metric during a concrete time window.
This makes them well-suited for tasks such as analyzing the results of a load test
that has been executed after the deployment of a new version.

In our case, we will create two `AnalysisValueTemplates` resources.
The first one measures the error rate of our workload, using data from Prometheus:

```yaml
{% include "./application-performance-analysis/error-rate.yaml" %}
```

As a second metric, we measure the memory usage of our application using the following `AnalysisValueTemplate`:

```yaml
{% include "./application-performance-analysis/memory-usage.yaml" %}
```

As can be seen in the `spec.query` field of the resource above,
`AnalysisValueTemplate` resources support the [Go templating syntax](https://pkg.go.dev/text/template).
With that, you can include placeholders in the query that are substituted at the time the
concrete values for the metrics are retrieved.
This comes in handy when, e.g., the query for the metrics is the same for different workloads
and only differs slightly, perhaps due to different label selectors being used for different workloads.
This way you do not need to create one `AnalysisValueTemplate` resource per workload
but can reuse one for different workloads, and pass through the value for the
actual workload at the time you perform an Analysis.

## Defining SLOs

The next step is to set up our expectations towards our SLOs, i.e. the
goals we would like them to meet.
This is done via an `AnalysisDefinition` resource like the following:

```yaml
{% include "./application-performance-analysis/analysis-definition.yaml" %}
```

This `AnalysisDefinition` resource has two objectives, which both refer
to the `AnalysisValueTemplate` resources we created previously.
If you closely inspect both, you will notice that they differ in the weights they have been assigned,
meaning that the goal for the error-rate has a higher priority than memory consumption.
In combination with the target scores defined in the totalScore object,
this means that passing the objective for the error-rate is mandatory for an analysis to be successful,
or at least to achieve the warning state.
The latter would be achieved if, for example, the error rate objective is passed,
but the memory consumption exceeds the defined limit of `30M`.
Also, note that even though we use values coming from different data sources,
i.e. Prometheus and Dynatrace, in the `AnalysisDefinition`, we do not need to consider any
implementation-specific details when referring to them.
You only need to provide the name of the `AnalysisValueTemplate`,
and the metrics-operator determines where to retrieve the data based on the information in the KeptnMetricsProviders.

## Executing an Analysis

Now, it is time to trigger an Analysis.
This is done by applying an Analysis resource which looks as follows:

```yaml
{% include "./application-performance-analysis/analysis.yaml" %}
```

Applying this resource causes Keptn to:

- Retrieve the values of the `AnalysisValueTemplate` resource referenced in the
`AnalysisDefinition` that is used for this Analysis instance.
- After all required values have been retrieved, the objectives of the `AnalysisDefinition` are evaluated,
and the overall result is computed.
- This analysis uses the values of the last ten minutes (due to spec.timeframe.recent being set to 10m).
Alternatively, you can also specify a concrete timeframe,
using the `spec.timeframe.from` and `spec.timeframe.to` properties.
- We also provide the argument workload to the analysis, using the `spec.args` property.
Arguments passed to the analysis via this property are used when computing the actual query,
using the templating string of the `AnalysisValueTemplates` resource.
In our case, we use this in the error-rate `AnalysisValueTemplate`, where we set the
query to `rate(http_requests_total{status_code='500', job='{{.workload}}'}[1m]) or on() vector(0)`.

Applying this resource causes Keptn to retrieve the values of the `AnalysisValueTemplate`
resource referenced in the `AnalysisDefinition` that is used for this `Analysis` instance.
After all required values have been retrieved, the objectives of the `AnalysisDefinition` are evaluated,
and the overall result is computed.

This analysis uses the values of the last ten minutes (due to `spec.timeframe.recent` being set to `10m`)
but you can also specify a concrete timeframe,
using the `spec.timeframe.from` and `spec.timeframe.to` properties.

We also provide the argument workload to the analysis, using the spec.args property.
Arguments passed to the analysis via this property are used when computing the actual query,
using the templating string of the `AnalysisValueTemplates` resource.
In our case, we use this in the error-rate `AnalysisValueTemplate`,
and set the query to `rate(http_requests_total{status_code='500', job='{{.workload}}'}[1m]) or on() vector(0)`.

For our `Analysis` with `spec.args.workload` set to `simple-go-service`, the resulting query is:

```shell
rate(http_requests_total{status_code='500', job='simple-go-service'}[1m]) or on() vector(0). 
```

## Inspecting the results

After applying an `Analysis` resource, we can do a quick check of its state using `kubectl`:

```shell
$ kubectl get analysis -n simple-go 
  
NAME               ANALYSISDEFINITION         STATE       WARNING   PASS 
service-analysis   my-analysis-definition     Completed             true 
```

The output of that command tells us if the Analysis has been completed already.
As seen above, this is the case, and we can already see that it has passed.
So now it's time to dive deeper into the results and see what information we get in the status of the resource:

```shell
kubectl get analysis service-analysis -n simple-go –oyaml
```

This command gives us the complete YAML representation of the `Analysis`:

```yaml
{% include "./application-performance-analysis/analysis-status.yaml" %}
```

As you can see, this already gives us a lot more information,
with the meatiest piece being the status.raw field.
This is a JSON representation of the retrieved values and the goals we have set for them.
However, this raw information is not easily digestible for our human eyes, so let's format it using:

```shell
kubectl get analyses service-analysis -n simple-go -o=jsonpath='{.status.raw}' | jq .
```

Giving us the following as a result:

```json
{% include "./application-performance-analysis/analysis-breakdown.json" %}
```

In the JSON object, we see:

- A list of the objectives that we defined earlier in our `AnalysisDefinition`
- The values of the related metrics
- The actual query that was used for retrieving their data from our monitoring data sources

Based on that, each objective is assigned a score,
which is equal to the weight of that objective if the objective has been met.
If not, the objective gets a score of 0.

Note: You can specify `warning` criteria in addition to `failure` criteria.
If that is the case, and the value of a metric does not violate the `failure` criteria,
but the warning criteria, it gets a score that is half of the weight.
This allows you to be even more granular with the grading of your analysis.
In our case, both objectives have been met, so we get the full score
and therefore pass the evaluation with flying colors.

## Summary

To summarize, we have seen how we can define multiple monitoring data sources
and let Keptn fetch the data we are interested in and provide us with a unified way of accessing this data.
In the next step, we created a clear set of criteria we expect from that data to
decide whether the related application is healthy or not.
Finally, we have seen how we can easily perform an analysis and interpret its result.
We did all this by using Kubernetes manifests and placing them in a
GitOps repository next to our application's manifests.

If you would like to try out Keptn and its analysis capabilities yourself,
feel free to head over to the [Keptn docs](https://lifecycle.keptn.sh/docs/)
and follow the guides to [install Keptn](https://lifecycle.keptn.sh/docs/install/),
if you haven't done so already,
and try out the [Analysis example](https://github.com/keptn/lifecycle-toolkit/tree/main/examples/support/analysis).
