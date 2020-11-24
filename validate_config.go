package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hashicorp/consul/api"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

const schemaPath string = "config/component/testing/service/config-validator/schema"

// ValidateSchema is
func ValidateSchema(config Config) {
	schema, err := getSchema(config)
	check(err)

	fileData, err := getFileData(config)
	check(err)

	result, isValid := runSchemaValidation(fileData, schema)
	check(err)

	if isValid {
		fmt.Printf("Config %s is valid\n", config.Path)
		os.Exit(0)
	} else {
		fmt.Printf("Config %s is invalid: %s\n", config.Path, result.Errors())
		os.Exit(1)
	}
}

func getSchema(config Config) (string, error) {
	keyPath := config.ConsulPrefix + "/" + schemaPath + "/" + config.Schema

	client, err := api.NewClient(config.ConsulConfig)
	check(err)

	// Get a handle to the KV API
	kv := client.KV()

	// Lookup the pair
	pair, _, err := kv.Get(keyPath, nil)
	check(err)

	return string(pair.Value), nil
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

func runSchemaValidation(fileData interface{}, schema string) (*gojsonschema.Result, bool) {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	configLoader := gojsonschema.NewGoLoader(fileData)

	result, err := gojsonschema.Validate(schemaLoader, configLoader)
	check(err)

	if result.Valid() {
		return result, true
	}

	return result, false
}
