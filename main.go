package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/consul/api"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

const schemaPath string = "config/component/testing/service/config-validator/schema"

func main() {
	flag.Parse()

	// Initialize flags, env variables, defaults
	config := GetConfig()

	// Get schema file from Consul
	schema, err := getSchema(config)
	check(err)

	// Get file data for validation
	fileData, err := getFileData(config)
	check(err)

	// Run the validation
	err = validateSchema(fileData, schema)
	check(err)
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
	result, err := ioutil.ReadFile(config.FilePath)
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

func validateSchema(fileData interface{}, schema string) error {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	configLoader := gojsonschema.NewGoLoader(fileData)

	result, err := gojsonschema.Validate(schemaLoader, configLoader)
	check(err)

	if !result.Valid() {
		return fmt.Errorf("Failed to validate document: %s", result.Errors())
	}

	fmt.Print("Template is valid\n")
	return nil

}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
