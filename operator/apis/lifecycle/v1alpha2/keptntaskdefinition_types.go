/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// KeptnTaskDefinitionSpec defines the desired state of KeptnTaskDefinition
type KeptnTaskDefinitionSpec struct {
	Function FunctionSpec `json:"function,omitempty"`
}

type FunctionSpec struct {
	FunctionReference  FunctionReference  `json:"functionRef,omitempty"`
	Inline             Inline             `json:"inline,omitempty"`
	HttpReference      HttpReference      `json:"httpRef,omitempty"`
	ConfigMapReference ConfigMapReference `json:"configMapRef,omitempty"`
	Parameters         TaskParameters     `json:"parameters,omitempty"`
	SecureParameters   SecureParameters   `json:"secureParameters,omitempty"`
}

type ConfigMapReference struct {
	Name string `json:"name,omitempty"`
}

type FunctionReference struct {
	Name string `json:"name,omitempty"`
}

type Inline struct {
	Code string `json:"code,omitempty"`
}

type HttpReference struct {
	Url string `json:"url,omitempty"`
}

type ContainerSpec struct {
}

// KeptnTaskDefinitionStatus defines the observed state of KeptnTaskDefinition
type KeptnTaskDefinitionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Function FunctionStatus `json:"function,omitempty"`
}

type FunctionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	ConfigMap string `json:"configMap,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// 
// A `KeptnTaskDefinition` defines the tasks
// that the Keptn Lifecycle Toolkit can run
// as part of pre- and post-deployment phases of a deployment.

// ## YAML synopsis

// ```yaml
// apiVersion: lifecycle.keptn.sh/v1alpha2
// kind: KeptnTaskDefinition
// metadata:
//   name: <task-name>
// spec:
//   function:
//     inline: | httpRef | functionRef
//     [parameters:
//       map:
//         textMessage: "This is my configuration"]
//     [secureParameters:
//       secret: slack-token]
// ```
//
// * **apiVersion** - API version being used.
//   The default for KLT Release 0.5.0 is `lifecycle.keptn.sh/v1alpha2`.
//   * This must match <whatever>
//   * Other information
//
// * **kind** CRD type.  This is `KeptnTaskDefinition`
//
// * **name** Unique name of this task.
//   This must be an alphanumeric string and, by convention, is all lowercase.
//   It can use the special characters `_`, `-` ... <what others>.
//   It should not include spaces.
//
// * **function** - Code to be executed.
//   This can be expressed as one of the following:
//
//   * **inline** - Include the actual executable code to execute.
//     This can be written as a full-fledged Deno script.
//     For example:
//     ```function:
//     inline:
//       code: |
//         console.log("Deployment Task has been executed");
//     ```
//     The task definition is a [Deno](https://deno.land/) script.
//     Please, refer to the [function runtime](https://github.com/keptn/lifecycle-toolkit/tree/main/functions-runtime)
//     for more information about the runtime.
//     In the future, we also intend to support other runtimes,
//     especially running a container image directly.
//   * **httpRef** - Specify a Deno script to be executed at runtime
//     from the remote webserver that is specified.
//     For example:
// 
//     ```yaml
//     name: hello-keptn-http
//       spec:
//         function:
//           httpRef:
//             url: <url>
//     ```
//   * **functionRef** -- Execute another `KeptnTaskDefinition` that has been defined.
//     Populate this field with the value of the `name` field
//     for the `KeptnTaskDefinition` to be called.
//     This is commonly used to call a general function
//     that is used in multiple place with different parameters.
//     An example is:
//      ```yaml
//      spec:
//        function:
//          functionRef:
//            name: slack-notification
//      ```
// 
//     This can also be used to group a set of tasks into a single `KeptnTaskDefinitions`,
//     such as defining a `KeptnTaskDefinition` for testing
//     and have it call a `KeptnTaskDefinition` for each type of test to be run.
// 
//     <Explain what happens if one task fails.
//    Will the subsequent tasks execute or does the pipeline stop
//    or pass control to whatever would execute next.
//    Can I control that behavior?
//
//* **parameters** - An optional field to supply input parameters to a function.
//  The Lifecycle Toolkit passes the values defined inside the `map` field
//  as a JSON object.
//  At the moment, multi-level maps are not supported.
//  For example:
//   ```spec:
//       parameters:
//         map:
//           textMessage: "This is my configuration"
//   ```
//
//   The JSON object can be read
//   through the `DATA` environment variable using `Deno.env.get("DATA");`.
//
//* **secureParameters** -- An optional field used to pass a Kubernetes secret.
//  The `secret` value is the Kubernetes secrete name
//  that is mounted into the runtime
//  and made available to the function
//  using the `SECURE_DATA` environment variable.
//
//  ```yaml
//      secureParameters:
//        secret: slack-token
//   ```
// ## Field spec
//
type KeptnTaskDefinition struct {
	metav1.TypeMeta   `json:",inline"`
metav1.ObjectMeta `json:"metadata,omitempty"`

Spec   KeptnTaskDefinitionSpec   `json:"spec,omitempty"`
Status KeptnTaskDefinitionStatus `json:"status,omitempty"`

+kubebuilder:object:root=true

KeptnTaskDefinitionList contains a list of KeptnTaskDefinition
type KeptnTaskDefinitionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KeptnTaskDefinition `json:"items"`
}

func init() {
	SchemeBuilder.Register(&KeptnTaskDefinition{}, &KeptnTaskDefinitionList{})
}

// ## Usage
// 
// A Task executs the TaskDefinition of a workload.
// The execution is done by spawning a Kubernetes Job to handle a single Task.
// In its state, it tracks the current status of the Kubernetes Job it created.
// 
// <!-- How is this CRD "activated".  For example, which event uses this CRD -->
// <!-- Can I execute tasks in parallel? -->
// <!-- Instructions and guidelines for when and how to customize a CRD -->
// 
// ## Examples
// 
// This section can do any of the following:
// 
// * Include annotated examples
// * Link to formal `examples`; include an annotation about what they illustrate
// 
// ## Files
// 
// * link to source code file where this is defined.
// 
// ## Differences between versions
// 
// ## See also
// 
// [function runtime](https://github.com/keptn/lifecycle-toolkit/tree/main/functions-runtime)
//
// * Link to "use-case" guide pages that do something interesting with this CRD
// * Link to reference pages for any related CRDs

