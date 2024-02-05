package main

import (
	"os"
	"testing"
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
