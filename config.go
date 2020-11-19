package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

const schemaExtension string = ".json"

var path = flag.String("path", "", "File path to validate")

func init() {
	flag.Parse()
}

// Config is
type Config struct {
	FilePath     string
	Schema       string
	ConsulConfig *api.Config
	ConsulPrefix string
}

// GetConfig is
func GetConfig() Config {
	err := checkFlags()
	check(err)

	absPath, err := filepath.Abs(*path)
	check(err)

	schemaName, err := getSchemaName(*path)
	check(err)

	viper.SetDefault("consul_address", "127.0.0.1")
	viper.SetDefault("consul_port", "8500")
	viper.SetDefault("consul_scheme", "http")
	viper.SetDefault("consul_datacenter", "dc1")
	viper.SetDefault("consul_namepsace", "default")
	viper.SetDefault("consul_kv_prefix", "config")

	viper.AutomaticEnv()

	consulConfig := api.DefaultConfig()

	consulConfig.Address = viper.GetString("consul_address") + ":" + viper.GetString("consul_port")
	consulConfig.Scheme = viper.GetString("consul_scheme")
	consulConfig.Token = viper.GetString("consul_token")
	consulConfig.Datacenter = viper.GetString("consul_datacenter")
	consulConfig.Namespace = viper.GetString("consul_namepsace")

	config := Config{
		FilePath:     filepath.ToSlash(absPath),
		Schema:       schemaName,
		ConsulConfig: consulConfig,
		ConsulPrefix: viper.GetString("consul_kv_prefix"),
	}

	return config
}

func checkFlags() error {
	if *path == "" {
		return fmt.Errorf("Path for file to validate was not provided")
	}

	if _, err := os.Stat(*path); os.IsNotExist(err) {
		return fmt.Errorf("File to validate does not exist: %q", *path)
	}

	return nil
}

func getSchemaName(path string) (string, error) {
	extension := filepath.Ext(path)

	if extension != ".yml" && extension != ".yaml" {
		err := fmt.Errorf("Wrong extension type %q, expecting .yml or .yaml", extension)
		return "", err
	}

	absPath, err := filepath.Abs(path)
	check(err)

	dir := filepath.Dir(absPath)
	standardDir := filepath.ToSlash(dir)
	dirItems := strings.Split(standardDir, "/")
	schema := dirItems[len(dirItems)-1]

	err = schemaNameIsValid(schema)
	check(err)

	return schema + schemaExtension, nil
}

func schemaNameIsValid(name string) error {
	switch name {
	case "customer", "environment", "component", "site":
		return nil
	}

	return fmt.Errorf("Schema name %q that was extracted from file path is not valid", name)
}
