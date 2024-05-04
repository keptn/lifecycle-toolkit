package main

import (
	"fmt"
	"log"
	"os"

	apilifecycle "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1"
	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
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

	if err := transformKeptnApp(inputFile, outputFile); err != nil {
		log.Fatalf("Error transforming KeptnApp: %v", err)
	}

	fmt.Println("Transformation completed. Output written to", outputFile)
}

func transformKeptnApp(inputFile, outputFile string) error {
	inputContent, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	keptnAppV1, keptnAppContext, err := parseAndTransform(inputContent)
	if err != nil {
		return err
	}

	outputContent := combineYAML(keptnAppV1, keptnAppContext)
	if err := os.WriteFile(outputFile, []byte(outputContent), 0644); err != nil {
		return fmt.Errorf("error writing to output file: %w", err)
	}
	return nil
}

func parseAndTransform(inputContent []byte) ([]byte, []byte, error) {
	var keptnApp klcv1alpha3.KeptnApp
	if err := yaml.Unmarshal(inputContent, &keptnApp); err != nil {
		return nil, nil, fmt.Errorf("error unmarshalling KeptnApp: %w", err)
	}

	var keptnAppV1 apilifecycle.KeptnApp
	if err := yaml.Unmarshal(inputContent, &keptnAppV1); err != nil {
		return nil, nil, fmt.Errorf("error unmarshalling KeptnAppV1: %w", err)
	}

	addKeptnAnnotation(&keptnAppV1.ObjectMeta)
	keptnAppV1.TypeMeta.APIVersion = "lifecycle.keptn.sh/v1"

	keptnAppContext := transformKeptnAppContext(keptnApp)

	keptnAppV1YAML, err := yaml.Marshal(keptnAppV1)
	if err != nil {
		return nil, nil, fmt.Errorf("error marshalling KeptnAppV1 to YAML: %w", err)
	}

	keptnAppContextYAML, err := yaml.Marshal(keptnAppContext)
	if err != nil {
		return nil, nil, fmt.Errorf("error marshalling KeptnAppContext to YAML: %w", err)
	}

	return keptnAppV1YAML, keptnAppContextYAML, nil
}

func transformKeptnAppContext(keptnApp klcv1alpha3.KeptnApp) apilifecycle.KeptnAppContext {
	return apilifecycle.KeptnAppContext{
		TypeMeta: metav1.TypeMeta{
			Kind:       "KeptnAppContext",
			APIVersion: "lifecycle.keptn.sh/v1",
		},

		ObjectMeta: metav1.ObjectMeta{
			Name:      keptnApp.Name,
			Namespace: keptnApp.Namespace,
		},
		Spec: apilifecycle.KeptnAppContextSpec{
			DeploymentTaskSpec: apilifecycle.DeploymentTaskSpec{
				PreDeploymentTasks:        keptnApp.Spec.PreDeploymentTasks,
				PreDeploymentEvaluations:  keptnApp.Spec.PreDeploymentEvaluations,
				PostDeploymentTasks:       keptnApp.Spec.PostDeploymentTasks,
				PostDeploymentEvaluations: keptnApp.Spec.PostDeploymentEvaluations,
			},
		},
	}
}

func combineYAML(keptnAppV1YAML, keptnAppContextYAML []byte) string {
	return fmt.Sprintf("%s\n---\n%s", string(keptnAppV1YAML), string(keptnAppContextYAML))
}

func addKeptnAnnotation(resource *metav1.ObjectMeta) {
	annotations := resource.GetAnnotations()
	if annotations == nil {
		resource.Annotations = make(map[string]string, 1)
	}
	resource.Annotations[keptnAnnotation] = keptn
}
