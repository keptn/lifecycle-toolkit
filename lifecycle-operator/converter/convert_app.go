package main

import (
	"fmt"
	"log"
	"os"

	klcv1alpha3 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1alpha3"
	klcv1beta1 "github.com/keptn/lifecycle-toolkit/lifecycle-operator/apis/lifecycle/v1beta1"
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
		return fmt.Errorf("error reading input file: %v", err)
	}

	keptnAppV1beta1, keptnAppContext, err := parseAndTransform(inputContent)
	if err != nil {
		return err
	}

	outputContent := combineYAML(keptnAppV1beta1, keptnAppContext)
	if err := os.WriteFile(outputFile, []byte(outputContent), 0644); err != nil {
		return fmt.Errorf("error writing to output file: %v", err)
	}
	return nil
}

func parseAndTransform(inputContent []byte) ([]byte, []byte, error) {
	var keptnApp klcv1alpha3.KeptnApp
	if err := yaml.Unmarshal(inputContent, &keptnApp); err != nil {
		return nil, nil, fmt.Errorf("error unmarshalling KeptnApp: %v", err)
	}

	var keptnAppV1beta1 klcv1beta1.KeptnApp
	if err := yaml.Unmarshal(inputContent, &keptnAppV1beta1); err != nil {
		return nil, nil, fmt.Errorf("error unmarshalling KeptnAppV1beta1: %v", err)
	}

	addKeptnAnnotation(&keptnAppV1beta1.ObjectMeta)
	keptnAppV1beta1.TypeMeta.APIVersion = "lifecycle.keptn.sh/v1beta1"

	keptnAppContext := transformKeptnAppContext(keptnApp)

	keptnAppV1beta1YAML, err := yaml.Marshal(keptnAppV1beta1)
	if err != nil {
		return nil, nil, fmt.Errorf("error marshalling KeptnAppV1beta1 to YAML: %v", err)
	}

	keptnAppContextYAML, err := yaml.Marshal(keptnAppContext)
	if err != nil {
		return nil, nil, fmt.Errorf("error marshalling KeptnAppContext to YAML: %v", err)
	}

	return keptnAppV1beta1YAML, keptnAppContextYAML, nil
}

func transformKeptnAppContext(keptnApp klcv1alpha3.KeptnApp) klcv1beta1.KeptnAppContext {
	return klcv1beta1.KeptnAppContext{
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
}

func combineYAML(keptnAppV1beta1YAML, keptnAppContextYAML []byte) string {
	return fmt.Sprintf("%s\n---\n%s", string(keptnAppV1beta1YAML), string(keptnAppContextYAML))
}

func addKeptnAnnotation(resource *metav1.ObjectMeta) {
	annotations := resource.GetAnnotations()
	if annotations == nil {
		resource.Annotations = make(map[string]string, 1)
	}
	resource.Annotations[keptnAnnotation] = keptn
}
