package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

const schemaExtension string = ".json"

var path = flag.String("path", "", "Path to use for validation")
var validateConfig = flag.Bool("config", false, "Validate config file on path against JSON schema")
var validateStructure = flag.Bool("structure", false, "Validate file and directory structure on path")
var printVersion = flag.Bool("version", false, "Print version information")

func init() {
	flag.Parse()
}

// Config is
type Config struct {
	Path              string
	Schema            string
	ValidateConfig    bool
	ValidateStructure bool
	PrintVersion      bool
	ConsulConfig      *api.Config
	ConsulPrefix      string
}

// GetConfig is
func GetConfig() Config {
	var config Config

	if !*printVersion {
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
	}

	if *printVersion {
		config = Config{
			PrintVersion: true,
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
	schema := ""

	extension := filepath.Ext(path)

	if extension != ".yml" && extension != ".yaml" {
		err := fmt.Errorf("Wrong config file extension type %q, expecting .yml or .yaml", extension)
		return "", err
	}

	absPath, err := filepath.Abs(path)
	check(err)
	path = filepath.ToSlash(absPath)

	// component 	= component/component_name.yml
	// global 		= component/component_name/global.yml
	// subcomponent = component/component_name/subcomponent_name.yml
	// customer 	= customer/customer_name.yml
	// environment	= cutomer/customer_name/environment/environment_name.yml
	// site 		= cutomer/customer_name/environment/environment_name/site/site_name.yml
	// default

	globalConfig := regexp.MustCompile(".*global\\.yml")
	componentConfig := regexp.MustCompile(".*component/[a-zA-Z0-9-_.]+\\.yml")
	customerConfig := regexp.MustCompile(".*customer/[a-zA-Z0-9-_.]+\\.yml")
	environmentConfig := regexp.MustCompile(".*environment/[a-zA-Z0-9-_.]+\\.yml")
	siteConfig := regexp.MustCompile(".*site/[a-zA-Z0-9-_.]+\\.yml")

	switch {
	case globalConfig.Match([]byte(path)):
		schema = "global"
	case componentConfig.Match([]byte(path)):
		schema = "component"
	case customerConfig.Match([]byte(path)):
		schema = "customer"
	case environmentConfig.Match([]byte(path)):
		schema = "environment"
	case siteConfig.Match([]byte(path)):
		schema = "site"
	default:
		schema = "default"
	}

	return schema + schemaExtension, nil
}
