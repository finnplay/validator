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

var path = flag.String("path", "", "Path to use for validation")
var validateConfig = flag.Bool("config", false, "Validate config file on path against JSON schema")
var validateStructure = flag.Bool("structure", false, "Validate file and directory structure on path")

func init() {
	flag.Parse()
}

// Config is
type Config struct {
	Path              string
	Schema            string
	ValidateConfig    bool
	ValidateStructure bool
	ConsulConfig      *api.Config
	ConsulPrefix      string
}

// GetConfig is
func GetConfig() Config {
	var config Config

	err := checkFlags()
	check(err)

	absPath, err := filepath.Abs(*path)
	check(err)

	if *validateConfig {
		schemaName, err := getSchemaName(*path)
		check(err)

		viper.SetDefault("consul_address", "127.0.0.1")
		viper.SetDefault("consul_port", "8500")
		viper.SetDefault("consul_scheme", "http")
		viper.SetDefault("consul_datacenter", "dc1")
		viper.SetDefault("consul_namepsace", "default")
		viper.SetDefault("consul_kv_prefix", "monitoring-poc")

		viper.AutomaticEnv()

		consulConfig := api.DefaultConfig()

		consulConfig.Address = viper.GetString("consul_address") + ":" + viper.GetString("consul_port")
		consulConfig.Scheme = viper.GetString("consul_scheme")
		consulConfig.Token = viper.GetString("consul_token")
		consulConfig.Datacenter = viper.GetString("consul_datacenter")
		consulConfig.Namespace = viper.GetString("consul_namepsace")

		config = Config{
			Path:           filepath.ToSlash(absPath),
			Schema:         schemaName,
			ValidateConfig: *validateConfig,
			ConsulConfig:   consulConfig,
			ConsulPrefix:   viper.GetString("consul_kv_prefix"),
		}
	}

	if *validateStructure {
		config = Config{
			Path:              filepath.ToSlash(absPath),
			ValidateStructure: *validateStructure,
		}
	}

	return config
}

func checkFlags() error {
	if *path == "" {
		return fmt.Errorf("Path to use for validation was not provided")
	}

	if _, err := os.Stat(*path); os.IsNotExist(err) {
		return fmt.Errorf("File to validate does not exist: %q", *path)
	}

	if !*validateConfig && !*validateStructure {
		return fmt.Errorf("Flag on what to validate was not provided")
	}

	if *validateConfig && *validateStructure {
		return fmt.Errorf("Please provide either -schema or -structure flags, not both")
	}

	return nil
}

func getSchemaName(path string) (string, error) {
	extension := filepath.Ext(path)

	if extension != ".yml" && extension != ".yaml" {
		err := fmt.Errorf("Wrong config file extension type %q, expecting .yml or .yaml", extension)
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
