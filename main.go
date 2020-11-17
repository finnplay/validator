package main

import (
	"fmt"
	"io/ioutil"

	"github.com/finnplay/validator/config"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v2"
)

const schemaPath string = "config/component/testing/service/config-validator/schema"

func main() {
	// Initialize flags, env variables, defaults
	cfg := config.Initialize()

	// Initialize Consul config
	cfgConsul := config.Consul(cfg)

	fmt.Printf("%+v\n", cfgConsul)

	// Get schema file from Consul
	// schema := getSchema(*schemaType)

}

/*
func getSchema(schemaType string) string {

	keyPath := schemaPath + "/" + schemaType

	// Initialize default config
	apiConfig := api.DefaultConfig()

	client, err := api.NewClient(apiConfig)
	check(err)

	// Get a handle to the KV API
	kv := client.KV()

	// Lookup the pair
	pair, _, err := kv.Get(keyPath, nil)
	check(err)

	return string(pair.Value)
}
*/

func prepareConfig(path string) interface{} {
	result, err := ioutil.ReadFile(path)
	check(err)

	var document interface{}
	if err := yaml.Unmarshal([]byte(result), &document); err != nil {
		panic(err)
	}

	var config interface{} = convertConfig(document)

	return config
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

func validateSchema(config interface{}, schema string) {
	schemaLoader := gojsonschema.NewStringLoader(schema)
	configLoader := gojsonschema.NewGoLoader(config)

	result, err := gojsonschema.Validate(schemaLoader, configLoader)

	check(err)

	if result.Valid() {
		fmt.Printf("The document is valid\n")
	} else {
		fmt.Printf("The document is not valid. see errors :\n")
		for _, desc := range result.Errors() {
			fmt.Printf("- %s\n", desc)
		}
	}
}

func check(e error) {
	if e != nil {
		panic(e.Error())
	}
}
