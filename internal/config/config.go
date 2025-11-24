package config

import "os"

type Config struct {
	ServerAddress string
	MockDataPath  string
}

func Load() *Config {
	cfg := &Config{
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		MockDataPath:  getEnv("MOCK_DATA_PATH", "./mock/kaspi_data.json"),
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
