package config

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// AppConfig holds core application configurations
type AppConfig struct {
	Name  string `mapstructure:"name" json:"name"`
	Env   string `mapstructure:"env" json:"env"`
	Port  string `mapstructure:"port" json:"port"`
	Debug bool   `mapstructure:"debug" json:"debug"`
}

// DatabaseConfig holds database connection and isolation settings
type DatabaseConfig struct {
	Driver    string `mapstructure:"driver" json:"driver"`
	Database  string `mapstructure:"database" json:"database"`
	Isolation string `mapstructure:"isolation" json:"isolation"`
}

// Config represents the complete application configuration
type Config struct {
	App      AppConfig      `mapstructure:"app" json:"app"`
	Database DatabaseConfig `mapstructure:"database" json:"database"`
}

// Load reads configuration from .env file, config/app.yaml, and environment variables
func Load() (*Config, error) {
	_ = godotenv.Load()

	viper.AutomaticEnv()

	viper.SetDefault("app.name", "tgo-app")
	viper.SetDefault("app.env", "development")
	viper.SetDefault("app.port", "8080")
	viper.SetDefault("app.debug", true)
	viper.SetDefault("database.driver", "sqlite")
	viper.SetDefault("database.database", "data/app.db")
	viper.SetDefault("database.isolation", "single")

	// Support env overrides like APP_PORT, APP_NAME, etc.
	if port := viper.GetString("APP_PORT"); port != "" {
		viper.Set("app.port", port)
	}
	if name := viper.GetString("APP_NAME"); name != "" {
		viper.Set("app.name", name)
	}
	if env := viper.GetString("APP_ENV"); env != "" {
		viper.Set("app.env", env)
	}
	if dbDriver := viper.GetString("DB_DRIVER"); dbDriver != "" {
		viper.Set("database.driver", dbDriver)
	}
	if dbName := viper.GetString("DB_DATABASE"); dbName != "" {
		viper.Set("database.database", dbName)
	}

	viper.SetConfigName("app")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(".")

	_ = viper.ReadInConfig()

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Ensure fallback if empty
	if cfg.App.Port == "" {
		cfg.App.Port = "8080"
	}
	if cfg.App.Name == "" {
		cfg.App.Name = "tgo-app"
	}
	if cfg.App.Env == "" {
		cfg.App.Env = "development"
	}

	return &cfg, nil
}

// MustLoad calls Load and panics if there is an error (Fail Fast & Loud)
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load configuration: %v", err))
	}
	return cfg
}
