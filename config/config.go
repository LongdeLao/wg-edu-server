// Package config provides configuration management for the application
package config

// Config holds all application configurations
type Config struct {
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	JWTSecret  string
	ServerPort string
}

// NewConfig returns a new Config with the default values
//
// Returns:
//   - Config: Configuration object with default values
func NewConfig() Config {
	return Config{
		DBHost:     "47.79.6.81",
		DBPort:     "5432",
		DBName:     "WG_EDU",
		DBUser:     "longdelao",
		DBPassword: "2008",
		JWTSecret:  "wg-edu-secret-key",
		ServerPort: ":8080",
	}
}
