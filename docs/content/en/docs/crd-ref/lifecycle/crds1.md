# API Reference

Packages:

- [options.keptn.sh/v1alpha1](#optionskeptnshv1alpha1)

# options.keptn.sh/v1alpha1

Resource Types:

- [KeptnConfig](#keptnconfig)




## KeptnConfig
<sup><sup>[↩ Parent](#optionskeptnshv1alpha1 )</sup></sup>






KeptnConfig is the Schema for the keptnconfigs API

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
      <td>options.keptn.sh/v1alpha1</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b>kind</b></td>
      <td>string</td>
      <td>KeptnConfig</td>
      <td>true</td>
      </tr>
      <tr>
      <td><b><a href="https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.20/#objectmeta-v1-meta">metadata</a></b></td>
      <td>object</td>
      <td>Refer to the Kubernetes API documentation for the fields of the `metadata` field.</td>
      <td>true</td>
      </tr><tr>
        <td><b><a href="#keptnconfigspec">spec</a></b></td>
        <td>object</td>
        <td>
          KeptnConfigSpec defines the desired state of KeptnConfig<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>status</b></td>
        <td>object</td>
        <td>
          KeptnConfigStatus defines the observed state of KeptnConfig<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>


### KeptnConfig.spec
<sup><sup>[↩ Parent](#keptnconfig)</sup></sup>



KeptnConfigSpec defines the desired state of KeptnConfig

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
        <td><b>OTelCollectorUrl</b></td>
        <td>string</td>
        <td>
          OTelCollectorUrl can be used to set the Open Telemetry collector that the operator should use<br/>
        </td>
        <td>false</td>
      </tr><tr>
        <td><b>keptnAppCreationRequestTimeout</b></td>
        <td>integer</td>
        <td>
          KeptnAppCreationRequestTimeout is used to set the interval in which automatic app discovery searches for workload to put into the same auto-generated KeptnApp<br/>
          <br/>
            <i>Default</i>: 30<br/>
        </td>
        <td>false</td>
      </tr></tbody>
</table>