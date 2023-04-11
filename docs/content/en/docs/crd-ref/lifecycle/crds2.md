# API Reference

Packages:

- [lifecycle.keptn.sh/v1alpha1](#lifecyclekeptnshv1alpha1)
- [lifecycle.keptn.sh/v1alpha2](#lifecyclekeptnshv1alpha2)
- [lifecycle.keptn.sh/v1alpha3](#lifecyclekeptnshv1alpha3)

# lifecycle.keptn.sh/v1alpha1

Resource Types:

- [KeptnEvaluation](#keptnevaluation)




## KeptnEvaluation
<sup><sup>[↩ Parent](#lifecyclekeptnshv1alpha1 )</sup></sup>






KeptnEvaluation is the Schema for the keptnevaluations API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>lifecycle.keptn.sh/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>KeptnEvaluation</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#keptnevaluationspec">spec</a></b></td>
        <td>object</td>
        <td>
          KeptnEvaluationSpec defines the desired state of KeptnEvaluation<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#keptnevaluationstatus">status</a></b></td>
        <td>object</td>
        <td>
          KeptnEvaluationStatus defines the observed state of KeptnEvaluation<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### KeptnEvaluation.spec
<sup><sup>[↩ Parent](#keptnevaluation)</sup></sup>



KeptnEvaluationSpec defines the desired state of KeptnEvaluation

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>evaluationDefinition</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>workloadVersion</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>appName</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>appVersion</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>checkType</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failAction</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>retries</b></td>
        <td>integer</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: 10<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>retryInterval</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: 5s<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>workload</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### KeptnEvaluation.status
<sup><sup>[↩ Parent](#keptnevaluation)</sup></sup>



KeptnEvaluationStatus defines the observed state of KeptnEvaluation

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#keptnevaluationstatusevaluationstatuskey">evaluationStatus</a></b></td>
        <td>map[string]object</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>overallStatus</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: Pending<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>retryCount</b></td>
        <td>integer</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: 0<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>endTime</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>startTime</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### KeptnEvaluation.status.evaluationStatus[key]
<sup><sup>[↩ Parent](#keptnevaluationstatus)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

# lifecycle.keptn.sh/v1alpha2

Resource Types:

- [KeptnEvaluation](#keptnevaluation)




## KeptnEvaluation
<sup><sup>[↩ Parent](#lifecyclekeptnshv1alpha2 )</sup></sup>






KeptnEvaluation is the Schema for the keptnevaluations API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>lifecycle.keptn.sh/v1alpha2</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>KeptnEvaluation</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#keptnevaluationspec-1">spec</a></b></td>
        <td>object</td>
        <td>
          KeptnEvaluationSpec defines the desired state of KeptnEvaluation<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#keptnevaluationstatus-1">status</a></b></td>
        <td>object</td>
        <td>
          KeptnEvaluationStatus defines the observed state of KeptnEvaluation<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### KeptnEvaluation.spec
<sup><sup>[↩ Parent](#keptnevaluation-1)</sup></sup>



KeptnEvaluationSpec defines the desired state of KeptnEvaluation

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>evaluationDefinition</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>workloadVersion</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>appName</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>appVersion</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>checkType</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failAction</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>retries</b></td>
        <td>integer</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: 10<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>retryInterval</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: 5s<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>workload</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### KeptnEvaluation.status
<sup><sup>[↩ Parent](#keptnevaluation-1)</sup></sup>



KeptnEvaluationStatus defines the observed state of KeptnEvaluation

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#keptnevaluationstatusevaluationstatuskey-1">evaluationStatus</a></b></td>
        <td>map[string]object</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>overallStatus</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: Pending<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>retryCount</b></td>
        <td>integer</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: 0<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>endTime</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>startTime</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### KeptnEvaluation.status.evaluationStatus[key]
<sup><sup>[↩ Parent](#keptnevaluationstatus-1)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>

# lifecycle.keptn.sh/v1alpha3

Resource Types:

- [KeptnEvaluation](#keptnevaluation)




## KeptnEvaluation
<sup><sup>[↩ Parent](#lifecyclekeptnshv1alpha3 )</sup></sup>






KeptnEvaluation is the Schema for the keptnevaluations API

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
      <td><b>apiVersion</b></td>
      <td>string</td>
      <td>lifecycle.keptn.sh/v1alpha3</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>KeptnEvaluation</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#keptnevaluationspec-1">spec</a></b></td>
        <td>object</td>
        <td>
          KeptnEvaluationSpec defines the desired state of KeptnEvaluation<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b><a href="#keptnevaluationstatus-1">status</a></b></td>
        <td>object</td>
        <td>
          KeptnEvaluationStatus defines the observed state of KeptnEvaluation<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### KeptnEvaluation.spec
<sup><sup>[↩ Parent](#keptnevaluation-1)</sup></sup>



KeptnEvaluationSpec defines the desired state of KeptnEvaluation

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>evaluationDefinition</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>workloadVersion</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>appName</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>appVersion</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>checkType</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>failAction</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>retries</b></td>
        <td>integer</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: 10<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>retryInterval</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: 5s<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>workload</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### KeptnEvaluation.status
<sup><sup>[↩ Parent](#keptnevaluation-1)</sup></sup>



KeptnEvaluationStatus defines the observed state of KeptnEvaluation

<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b><a href="#keptnevaluationstatusevaluationstatuskey-1">evaluationStatus</a></b></td>
        <td>map[string]object</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>overallStatus</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: Pending<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>retryCount</b></td>
        <td>integer</td>
        <td>
          <br/>
          <br/>
            <i>Default</i>: 0<br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>endTime</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>startTime</b></td>
        <td>string</td>
        <td>
          <br/>
          <br/>
            <i>Format</i>: date-time<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### KeptnEvaluation.status.evaluationStatus[key]
<sup><sup>[↩ Parent](#keptnevaluationstatus-1)</sup></sup>





<table>
    <thead>
        <tr>
            <th>Name</th>
            <th>Type</th>
            <th>Description</th>
            <th>Required</th>
        </tr>
    </thead>
    <tbody><tr>
        <td><b>status</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>value</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>message</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>