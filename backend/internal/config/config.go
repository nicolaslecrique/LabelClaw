package config

import "os"

type Config struct {
	Addr              string
	ConfigurationPath string
}

func Load() Config {
	return Config{
		Addr:              getenv("LABELCLAW_ADDR", "127.0.0.1:8080"),
		ConfigurationPath: getenv("LABELCLAW_CONFIG_PATH", "data/active-config.json"),
	}
}

func getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
