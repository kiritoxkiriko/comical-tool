package config

import "os"

type Config struct {
	BaseURL    string
	AdminToken string
	Output     string
}

func Default() Config {
	baseURL := os.Getenv("COMICAL_API_BASE_URL")
	if baseURL == "" {
		baseURL = "http://localhost:8080"
	}
	return Config{
		BaseURL:    baseURL,
		AdminToken: os.Getenv("COMICAL_ADMIN_TOKEN"),
		Output:     "json",
	}
}
