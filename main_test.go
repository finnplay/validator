package main

import "testing"

func TestGetSchemaNameRelativePath(t *testing.T) {
	testPath := "./environment/prod.yml"
	name, _ := getSchemaName(testPath)

	if name != "environment.json" {
		t.Errorf("Parsed schema name was incorrect, got: %q, expected: environment.json", name)
	}
}

func TestGetSchemaNameAbsolutePath(t *testing.T) {
	testPath := "/tmp/config/customer/tabella/environment/prod.yml"
	name, _ := getSchemaName(testPath)

	if name != "environment.json" {
		t.Errorf("Parsed schema name was incorrect, got: %q, expected: environment.json", name)
	}
}
