package main

import (
	"fmt"
	"log"
	"os"
	"sigs.k8s.io/yaml"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	klcv1beta1 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const keptnAnnotation = "app.kubernetes.io/managed-by"
const keptn = "keptn"

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go app_cr_input_file.yaml output_file.yaml")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	inputContent, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	var keptnApp klcv1alpha3.KeptnApp
	if err := yaml.Unmarshal(inputContent, &keptnApp); err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}

	var keptnAppV1beta1 klcv1beta1.KeptnApp
	if err := yaml.Unmarshal(inputContent, &keptnAppV1beta1); err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}

	addKeptnAnnotation(&keptnAppV1beta1.ObjectMeta)
	keptnAppV1beta1.TypeMeta.APIVersion = "lifecycle.keptn.sh/v1beta1"

	// Transform KeptnAppContext
	keptnAppContext := klcv1beta1.KeptnAppContext{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KeptnAppContext",
			APIVersion: "lifecycle.keptn.sh/v1beta1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      keptnApp.Name,
			Namespace: keptnApp.Namespace,
		},
		Spec: klcv1beta1.KeptnAppContextSpec{
			DeploymentTaskSpec: klcv1beta1.DeploymentTaskSpec{

				PreDeploymentTasks:        keptnApp.Spec.PreDeploymentTasks,
				PreDeploymentEvaluations:  keptnApp.Spec.PreDeploymentEvaluations,
				PostDeploymentTasks:       keptnApp.Spec.PostDeploymentTasks,
				PostDeploymentEvaluations: keptnApp.Spec.PostDeploymentEvaluations,
			},
		},
	}

	// Convert to YAML and write to output file
	keptnAppV1beta1YAML, err := yaml.Marshal(keptnAppV1beta1)
	if err != nil {
		log.Fatalf("Error marshalling to YAML: %v", err)
	}

	keptnAppContextYAML, err := yaml.Marshal(keptnAppContext)
	if err != nil {
		log.Fatalf("Error marshalling to YAML: %v", err)
	}

	// Combine and write to output file
	outputContent := fmt.Sprintf("%s\n---\n%s", string(keptnAppV1beta1YAML), string(keptnAppContextYAML))
	if err := os.WriteFile(outputFile, []byte(outputContent), 0644); err != nil {
		log.Fatalf("Error writing to output file: %v", err)
	}

	fmt.Println("Transformation completed. Output written to", outputFile)
}

func addKeptnAnnotation(resource *metav1.ObjectMeta) {
	annotations := resource.GetAnnotations()
	if annotations == nil {
		resource.Annotations = make(map[string]string, 1)
	}
	resource.Annotations[keptnAnnotation] = keptn
}
