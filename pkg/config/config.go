package config

// Config - stores config.
type Config struct {
	Welcome string       `json:"welcome"`
	Server  ServerConfig `json:"server"`
}

// ServerConfig - stores server settings
type ServerConfig struct {
	Port string `json:"port"`
}
