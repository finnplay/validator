package config

import (
	"flag"

	"github.com/hashicorp/consul/api"
	"github.com/spf13/viper"
)

// Config is
type Config struct {
	Files            string
	VaultAddress     string
	VaultPort        string
	VaultToken       string
	VaultScheme      string
	VaultPrefix      string
	ConsulAddress    string
	ConsulPort       string
	ConsulToken      string
	ConsulScheme     string
	ConsulPrefix     string
	ConsulDatacenter string
	ConsulNamespace  string
}

// Initialize is
func Initialize() Config {
	files := flag.String("schema", "defaults.json", "Comma-separated list of file paths to validate")

	flag.Parse()

	viper.SetDefault("vault_address", "127.0.0.1")
	viper.SetDefault("vault_port", "8200")
	viper.SetDefault("vault_scheme", "http")
	viper.SetDefault("vault_kv_prefix", "secret")
	viper.SetDefault("consul_address", "127.0.0.1")
	viper.SetDefault("consul_port", "8500")
	viper.SetDefault("consul_scheme", "http")
	viper.SetDefault("consul_datacenter", "dc1")
	viper.SetDefault("consul_namepsace", "default")
	viper.SetDefault("consul_kv_prefix", "config")

	viper.AutomaticEnv()

	viper.Get("vault_address")
	viper.Get("vault_port")
	viper.Get("vault_token")
	viper.Get("vault_scheme")
	viper.Get("vault_kv_prefix")
	viper.Get("consul_address")
	viper.Get("consul_port")
	viper.Get("consul_token")
	viper.Get("consul_scheme")
	viper.Get("consul_datacenter")
	viper.Get("consul_namepsace")
	viper.Get("consul_kv_prefix")

	return Config{
		Files:            *files,
		VaultAddress:     viper.GetString("vault_address"),
		VaultPort:        viper.GetString("vault_port"),
		VaultToken:       viper.GetString("vault_token"),
		VaultScheme:      viper.GetString("vault_scheme"),
		VaultPrefix:      viper.GetString("vault_kv_prefix"),
		ConsulAddress:    viper.GetString("consul_address"),
		ConsulPort:       viper.GetString("consul_port"),
		ConsulToken:      viper.GetString("consul_token"),
		ConsulScheme:     viper.GetString("consul_scheme"),
		ConsulDatacenter: viper.GetString("consul_datacenter"),
		ConsulNamespace:  viper.GetString("consul_namepsace"),
		ConsulPrefix:     viper.GetString("consul_kv_prefix"),
	}
}

// Consul is
func Consul(cfg Config) api.Config {

	consulConfig := api.DefaultConfig()

	address := cfg.ConsulAddress + ":" + cfg.ConsulPort

	consulConfig.Address = address
	consulConfig.Scheme = cfg.ConsulScheme
	consulConfig.Token = cfg.ConsulToken
	consulConfig.Datacenter = cfg.ConsulDatacenter
	consulConfig.Namespace = cfg.ConsulNamespace

	return *consulConfig
}
