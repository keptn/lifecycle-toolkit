# API Reference

Packages:

- [lifecycle.keptn.sh/v1alpha1](#lifecyclekeptnshv1alpha1)
- [lifecycle.keptn.sh/v1alpha2](#lifecyclekeptnshv1alpha2)
- [lifecycle.keptn.sh/v1alpha3](#lifecyclekeptnshv1alpha3)

# lifecycle.keptn.sh/v1alpha1

Resource Types:

- [KeptnEvaluationDefinition](#keptnevaluationdefinition)




## KeptnEvaluationDefinition
<sup><sup>[↩ Parent](#lifecyclekeptnshv1alpha1 )</sup></sup>






KeptnEvaluationDefinition is the Schema for the keptnevaluationdefinitions API

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
      <td>KeptnEvaluationDefinition</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#keptnevaluationdefinitionspec">spec</a></b></td>
        <td>object</td>
        <td>
          KeptnEvaluationDefinitionSpec defines the desired state of KeptnEvaluationDefinition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>object</td>
        <td>
          KeptnEvaluationDefinitionStatus defines the observed state of KeptnEvaluationDefinition<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### KeptnEvaluationDefinition.spec
<sup><sup>[↩ Parent](#keptnevaluationdefinition)</sup></sup>



KeptnEvaluationDefinitionSpec defines the desired state of KeptnEvaluationDefinition

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
        <td><b><a href="#keptnevaluationdefinitionspecobjectivesindex">objectives</a></b></td>
        <td>[]object</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>source</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### KeptnEvaluationDefinition.spec.objectives[index]
<sup><sup>[↩ Parent](#keptnevaluationdefinitionspec)</sup></sup>





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
        <td><b>evaluationTarget</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>query</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>

# lifecycle.keptn.sh/v1alpha2

Resource Types:

- [KeptnEvaluationDefinition](#keptnevaluationdefinition)




## KeptnEvaluationDefinition
<sup><sup>[↩ Parent](#lifecyclekeptnshv1alpha2 )</sup></sup>






KeptnEvaluationDefinition is the Schema for the keptnevaluationdefinitions API

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
      <td>KeptnEvaluationDefinition</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#keptnevaluationdefinitionspec-1">spec</a></b></td>
        <td>object</td>
        <td>
          KeptnEvaluationDefinitionSpec defines the desired state of KeptnEvaluationDefinition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>object</td>
        <td>
          KeptnEvaluationDefinitionStatus defines the observed state of KeptnEvaluationDefinition<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### KeptnEvaluationDefinition.spec
<sup><sup>[↩ Parent](#keptnevaluationdefinition-1)</sup></sup>



KeptnEvaluationDefinitionSpec defines the desired state of KeptnEvaluationDefinition

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
        <td><b><a href="#keptnevaluationdefinitionspecobjectivesindex-1">objectives</a></b></td>
        <td>[]object</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>source</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### KeptnEvaluationDefinition.spec.objectives[index]
<sup><sup>[↩ Parent](#keptnevaluationdefinitionspec-1)</sup></sup>





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
        <td><b>evaluationTarget</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>query</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>

# lifecycle.keptn.sh/v1alpha3

Resource Types:

- [KeptnEvaluationDefinition](#keptnevaluationdefinition)




## KeptnEvaluationDefinition
<sup><sup>[↩ Parent](#lifecyclekeptnshv1alpha3 )</sup></sup>






KeptnEvaluationDefinition is the Schema for the keptnevaluationdefinitions API

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
      <td>KeptnEvaluationDefinition</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#keptnevaluationdefinitionspec-1">spec</a></b></td>
        <td>object</td>
        <td>
          KeptnEvaluationDefinitionSpec defines the desired state of KeptnEvaluationDefinition<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>object</td>
        <td>
          KeptnEvaluationDefinitionStatus defines the observed state of KeptnEvaluationDefinition<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### KeptnEvaluationDefinition.spec
<sup><sup>[↩ Parent](#keptnevaluationdefinition-1)</sup></sup>



KeptnEvaluationDefinitionSpec defines the desired state of KeptnEvaluationDefinition

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
        <td><b><a href="#keptnevaluationdefinitionspecobjectivesindex-1">objectives</a></b></td>
        <td>[]object</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### KeptnEvaluationDefinition.spec.objectives[index]
<sup><sup>[↩ Parent](#keptnevaluationdefinitionspec-1)</sup></sup>





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
        <td><b>evaluationTarget</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b><a href="#keptnevaluationdefinitionspecobjectivesindexkeptnmetricref">keptnMetricRef</a></b></td>
        <td>object</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr></tbody>
</table>


### KeptnEvaluationDefinition.spec.objectives[index].keptnMetricRef
<sup><sup>[↩ Parent](#keptnevaluationdefinitionspecobjectivesindex-1)</sup></sup>





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
        <td><b>name</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>true</td>
      </tr><tr>
        <td><b>namespace</b></td>
        <td>string</td>
        <td>
          <br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>