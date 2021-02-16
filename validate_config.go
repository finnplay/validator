package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

const schemaPath string = "config/component/automation/validator/schema"

// ValidateSchema is
func ValidateSchema(config Config) {
	fileData, err := getFileData(config)
	check(err)

	schemaPath, _ := filepath.Abs("schema/" + config.Schema)
	schemaPath = filepath.ToSlash(schemaPath)

	result, isValid := runSchemaValidation(fileData, schemaPath)
	check(err)

	if isValid {
		fmt.Printf("Config %s is valid\n", config.Path)
		os.Exit(0)
	} else {
		fmt.Printf("Config %s is invalid:\n", config.Path)

		for _, error := range result.Errors() {
			fmt.Printf("%s\n", error)
		}
		os.Exit(1)
	}
}

func getFileData(config Config) (interface{}, error) {
	result, err := ioutil.ReadFile(config.Path)
	check(err)

	var document interface{}
	if err := yaml.Unmarshal([]byte(result), &document); err != nil {
		panic(err)
	}

	var fileData interface{} = convertConfig(document)

	return fileData, nil
}

func convertConfig(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convertConfig(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convertConfig(v)
		}
	}
	return i
}

func runSchemaValidation(fileData interface{}, schemaPath string) (*gojsonschema.Result, bool) {
	schemaLoader := gojsonschema.NewReferenceLoader("file://" + schemaPath)
	configLoader := gojsonschema.NewGoLoader(fileData)

	result, err := gojsonschema.Validate(schemaLoader, configLoader)
	check(err)

	if result.Valid() {
		return result, true
	}

	return result, false
}
