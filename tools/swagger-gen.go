// +build ignore

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// This tool can be used to convert between JSON and YAML OpenAPI specs
// Usage: go run tools/swagger-gen.go

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run swagger-gen.go <command>")
		fmt.Println("Commands:")
		fmt.Println("  validate  - Validate the OpenAPI spec")
		fmt.Println("  json      - Convert YAML spec to JSON")
		fmt.Println("  yaml      - Convert JSON spec to YAML")
		return
	}

	command := os.Args[1]

	switch command {
	case "validate":
		validateSpec()
	case "json":
		convertToJSON()
	case "yaml":
		convertToYAML()
	default:
		fmt.Printf("Unknown command: %s\n", command)
	}
}

func validateSpec() {
	fmt.Println("Validating OpenAPI specification...")
	
	// Read the YAML file
	data, err := os.ReadFile("docs/swagger.yaml")
	if err != nil {
		log.Fatalf("Error reading swagger.yaml: %v", err)
	}

	// Parse YAML
	var spec map[string]interface{}
	err = yaml.Unmarshal(data, &spec)
	if err != nil {
		log.Fatalf("Error parsing YAML: %v", err)
	}

	// Basic validation
	requiredFields := []string{"openapi", "info", "paths"}
	for _, field := range requiredFields {
		if _, exists := spec[field]; !exists {
			log.Fatalf("Missing required field: %s", field)
		}
	}

	fmt.Println("✅ OpenAPI specification is valid!")
}

func convertToJSON() {
	fmt.Println("Converting YAML to JSON...")
	
	// Read YAML
	data, err := os.ReadFile("docs/swagger.yaml")
	if err != nil {
		log.Fatalf("Error reading swagger.yaml: %v", err)
	}

	// Parse YAML
	var spec map[string]interface{}
	err = yaml.Unmarshal(data, &spec)
	if err != nil {
		log.Fatalf("Error parsing YAML: %v", err)
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		log.Fatalf("Error converting to JSON: %v", err)
	}

	// Write JSON file
	err = os.WriteFile("docs/swagger.json", jsonData, 0644)
	if err != nil {
		log.Fatalf("Error writing JSON file: %v", err)
	}

	fmt.Println("✅ Converted to docs/swagger.json")
}

func convertToYAML() {
	fmt.Println("Converting JSON to YAML...")
	
	// Read JSON
	data, err := os.ReadFile("docs/swagger.json")
	if err != nil {
		log.Fatalf("Error reading swagger.json: %v", err)
	}

	// Parse JSON
	var spec map[string]interface{}
	err = json.Unmarshal(data, &spec)
	if err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Convert to YAML
	yamlData, err := yaml.Marshal(spec)
	if err != nil {
		log.Fatalf("Error converting to YAML: %v", err)
	}

	// Write YAML file
	err = os.WriteFile("docs/swagger.yaml", yamlData, 0644)
	if err != nil {
		log.Fatalf("Error writing YAML file: %v", err)
	}

	fmt.Println("✅ Converted to docs/swagger.yaml")
}