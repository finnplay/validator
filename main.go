package main

func main() {
	// Initialize flags, env variables, defaults
	config := GetConfig()

	switch {
	case config.ValidateConfig:
		ValidateSchema(config)
	case config.ValidateStructure:
		ValidateStructure(config)
	}
}
