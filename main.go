package main

import "fmt"

const version = "1.2"

func main() {
	// Initialize flags, env variables, defaults
	config := GetConfig()

	switch {
	case config.PrintVersion:
		fmt.Printf("Version: %s", version)
	case config.ValidateConfig:
		ValidateSchema(config)
	case config.ValidateStructure:
		ValidateStructure(config)
	}

}
