package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const inputFileName = "example_keptnapp.yaml"
const outputFileName = "example_output.yaml"

func TestMigration(t *testing.T) {
	// Set up a temporary directory for test files
	// Run the main function with test arguments
	inputFile := inputFileName
	outputFile := outputFileName
	os.Args = []string{"main", inputFile, outputFile}
	main()

	// Read the expected output file content
	expectedOutput, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Error reading expected output file: %v", err)
	}

	// Read the actual output file content
	actualOutput, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Error reading actual output file: %v", err)
	}

	// Assert that the actual output file content matches the expected output
	if string(actualOutput) != string(expectedOutput) {
		t.Errorf("Unexpected output content. Expected:\n%s\n\nActual:\n%s", string(expectedOutput), string(actualOutput))
	}

}

func TestAddKeptnAnnotation(t *testing.T) {
	// Test case 1: Annotations map is nil
	t.Run("AnnotationsMapIsNil", func(t *testing.T) {
		resource := &metav1.ObjectMeta{}
		addKeptnAnnotation(resource)

		// Check if the annotation was added
		if value, exists := resource.Annotations[keptnAnnotation]; !exists || value != keptn {
			t.Errorf("Annotation not added correctly. Expected: %s=%s, Actual: %s=%s", keptnAnnotation, keptn, keptnAnnotation, value)
		}
	})

	// Test case 2: Annotations map is not nil
	t.Run("AnnotationsMapIsNotNil", func(t *testing.T) {
		// Existing annotations
		existingAnnotations := map[string]string{
			"existing-key": "existing-value",
		}

		resource := &metav1.ObjectMeta{
			Annotations: existingAnnotations,
		}

		addKeptnAnnotation(resource)

		// Check if the annotation was added
		if value, exists := resource.Annotations[keptnAnnotation]; !exists || value != keptn {
			t.Errorf("Annotation not added correctly. Expected: %s=%s, Actual: %s=%s", keptnAnnotation, keptn, keptnAnnotation, value)
		}

		// Check if existing annotations are preserved
		for key, value := range existingAnnotations {
			if resource.Annotations[key] != value {
				t.Errorf("Existing annotation %s=%s is not preserved", key, value)
			}
		}
	})
}

func TestParseAndTransformBadYAML(t *testing.T) {
	invalidYAML := []byte("This is not valid YAML")
	// Attempt to parse and transform the invalid app YAML
	_, _, err := parseAndTransform(invalidYAML)
	if err == nil {
		t.Error("Expected an error but got nil")
	}
}

func TestTransformKeptnApp_UnmarshalFailure(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a temporary file with invalid YAML content
	invalidYAML := []byte("This is not valid YAML")
	err := os.WriteFile(tmpDir+"/validfile.yaml", invalidYAML, 0644)
	require.NoError(t, err)

	// Attempt to transform the invalid YAML file
	err = transformKeptnApp(tmpDir+"/validfile.yaml", tmpDir+"/output.yaml")

	require.Error(t, err)
}

func TestTransformKeptnApp_ReadFileFailure(t *testing.T) {
	//unexisting file
	tmpDir := t.TempDir()
	err := transformKeptnApp(tmpDir+"/validfile.yaml", tmpDir+"/output.yaml")
	require.Error(t, err)
	t.Log(err)
}
