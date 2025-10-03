package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Security SecurityConfig `mapstructure:"security"`
	Features FeatureFlags   `mapstructure:"features"`
}

type AppConfig struct {
	Name        string `mapstructure:"name"`
	Version     string `mapstructure:"version"`
	Environment string `mapstructure:"environment"`
	Mode        string `mapstructure:"mode"` // debug, release, test
}

type ServerConfig struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	IdleTimeout  int    `mapstructure:"idle_timeout"`
}

type DatabaseConfig struct {
	Driver          string `mapstructure:"driver"`
	DSN             string `mapstructure:"dsn"`
	MaxOpenConns    int    `mapstructure:"max_open_conns"`
	MaxIdleConns    int    `mapstructure:"max_idle_conns"`
	ConnMaxLifetime int    `mapstructure:"conn_max_lifetime"`
	LogLevel        string `mapstructure:"log_level"`
}

type LoggingConfig struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"` // json, text
	Output     string `mapstructure:"output"` // stdout, file
	FilePath   string `mapstructure:"file_path"`
	MaxSize    int    `mapstructure:"max_size"`    // megabytes
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"` // days
}

type SecurityConfig struct {
	JWTSecret            string   `mapstructure:"jwt_secret"`
	JWTExpiry            int      `mapstructure:"jwt_expiry"` // minutes
	BCryptCost           int      `mapstructure:"bcrypt_cost"`
	RateLimit            int      `mapstructure:"rate_limit"` // requests per minute
	AllowedOrigins       []string `mapstructure:"allowed_origins"`
	TrustedProxies       []string `mapstructure:"trusted_proxies"`
	EnableCSRF           bool     `mapstructure:"enable_csrf"`
	EnableSecurityHeaders bool     `mapstructure:"enable_security_headers"`
}

type FeatureFlags struct {
	EnableSwagger   bool `mapstructure:"enable_swagger"`
	EnableMetrics   bool `mapstructure:"enable_metrics"`
	EnableProfiling bool `mapstructure:"enable_profiling"`
	MaintenanceMode bool `mapstructure:"maintenance_mode"`
	DebugMode       bool `mapstructure:"debug_mode"`
}

// Load loads configuration from file and environment variables
func Load() (*Config, error) {
	v := viper.New()

	// Set defaults
	setDefaults(v)

	// Set config name and paths
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")
	v.AddConfigPath("/etc/banking-system/")

	// Read config file
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Override with environment-specific config
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	v.SetConfigName(fmt.Sprintf("config.%s", env))
	if err := v.MergeInConfig(); err == nil {
		fmt.Printf("Loaded environment config: config.%s.yaml\n", env)
	}

	// Bind environment variables
	v.SetEnvPrefix("BANK")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()

	// Unmarshal config
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	// Validate config
	if err := validate(&config); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &config, nil
}

func setDefaults(v *viper.Viper) {
	// App defaults
	v.SetDefault("app.name", "Banking System")
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.environment", "development")
	v.SetDefault("app.mode", "debug")

	// Server defaults
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.read_timeout", 15)
	v.SetDefault("server.write_timeout", 15)
	v.SetDefault("server.idle_timeout", 60)

	// Database defaults
	v.SetDefault("database.driver", "sqlite")
	v.SetDefault("database.dsn", "banking.db?_journal_mode=WAL")
	v.SetDefault("database.max_open_conns", 25)
	v.SetDefault("database.max_idle_conns", 5)
	v.SetDefault("database.conn_max_lifetime", 300)
	v.SetDefault("database.log_level", "info")

	// Logging defaults
	v.SetDefault("logging.level", "info")
	v.SetDefault("logging.format", "json")
	v.SetDefault("logging.output", "stdout")
	v.SetDefault("logging.max_size", 100)
	v.SetDefault("logging.max_backups", 3)
	v.SetDefault("logging.max_age", 7)

	// Security defaults
	v.SetDefault("security.jwt_expiry", 60)
	v.SetDefault("security.bcrypt_cost", 10)
	v.SetDefault("security.rate_limit", 100)
	v.SetDefault("security.enable_security_headers", true)

	// Feature flags defaults
	v.SetDefault("features.enable_metrics", true)
	v.SetDefault("features.enable_profiling", false)
	v.SetDefault("features.maintenance_mode", false)
	v.SetDefault("features.debug_mode", true)
}

func validate(cfg *Config) error {
	if cfg.Server.Port < 1 || cfg.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", cfg.Server.Port)
	}

	if cfg.Security.BCryptCost < 4 || cfg.Security.BCryptCost > 31 {
		return fmt.Errorf("invalid bcrypt cost: %d", cfg.Security.BCryptCost)
	}

	if cfg.Database.MaxOpenConns < 1 {
		return fmt.Errorf("invalid max open connections: %d", cfg.Database.MaxOpenConns)
	}

	return nil
}