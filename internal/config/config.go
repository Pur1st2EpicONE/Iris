// Package config provides the configuration structures and loading logic for the Iris application.
// It supports loading configuration from YAML files and environment variables.
package config

import (
	"fmt"
	"os"
	"time"

	wbf "github.com/wb-go/wbf/config"
)

// Config is the top-level application configuration, containing logger, server, storage and cache settings.
type Config struct {
	Logger  Logger  `mapstructure:"logger"`   // logger configuration
	Server  Server  `mapstructure:"server"`   // server configuration
	Storage Storage `mapstructure:"database"` // database/storage configuration
	Cache   Cache   `mapstructure:"cache"`    // cache configuration
}

// Logger defines logging configuration.
type Logger struct {
	Debug  bool   `mapstructure:"debug_mode"`    // enable debug mode
	LogDir string `mapstructure:"log_directory"` // directory for log files
}

// Server contains HTTP server configuration.
type Server struct {
	Port            string        `mapstructure:"port"`             // server port
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`     // HTTP read timeout
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`    // HTTP write timeout
	MaxHeaderBytes  int           `mapstructure:"max_header_bytes"` // max header size
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"` // graceful shutdown timeout
}

// Storage defines database connection and query retry configuration.
type Storage struct {
	Host               string        `mapstructure:"host"`                 // DB host
	Port               string        `mapstructure:"port"`                 // DB port
	Username           string        `mapstructure:"username"`             // DB username
	Password           string        `mapstructure:"password"`             // DB password
	DBName             string        `mapstructure:"dbname"`               // database name
	SSLMode            string        `mapstructure:"sslmode"`              // SSL mode
	MaxOpenConns       int           `mapstructure:"max_open_conns"`       // maximum open connections
	MaxIdleConns       int           `mapstructure:"max_idle_conns"`       // maximum idle connections
	ConnMaxLifetime    time.Duration `mapstructure:"conn_max_lifetime"`    // max lifetime per connection
	QueryRetryStrategy RetryStrategy `mapstructure:"query_retry_strategy"` // retry strategy for queries
}

// Cache defines Redis cache connection and retry configuration.
type Cache struct {
	Host           string        `mapstructure:"host"`            // cache host
	Port           string        `mapstructure:"port"`            // cache port
	Password       string        `mapstructure:"password"`        // cache password
	MaxMemory      string        `mapstructure:"max_memory"`      // max memory for Redis
	Policy         string        `mapstructure:"policy"`          // eviction policy
	RetryStrategy  RetryStrategy `mapstructure:"retry_strategy"`  // retry strategy for cache operations
	ExpirationTime time.Duration `mapstructure:"expiration_time"` // key expiration duration
}

// RetryStrategy defines retry behavior for operations.
type RetryStrategy struct {
	Attempts int           `mapstructure:"attempts"` // number of retry attempts
	Delay    time.Duration `mapstructure:"delay"`    // delay between retries
	Backoff  float64       `mapstructure:"backoff"`  // backoff multiplier
}

// Load reads configuration from files and environment variables.
func Load() (Config, error) {

	cfg := wbf.New()

	if err := cfg.LoadConfigFiles("./config.yaml"); err != nil {
		return Config{}, err
	}

	if err := cfg.LoadEnvFiles(".env"); err != nil && !cfg.GetBool("docker") {
		return Config{}, err
	}

	var conf Config

	if err := cfg.Unmarshal(&conf); err != nil {
		return Config{}, fmt.Errorf("unmarshal config: %w", err)
	}

	loadEnvs(&conf)

	return conf, nil

}

// loadEnvs overrides sensitive fields from environment variables.
func loadEnvs(conf *Config) {

	conf.Storage.Username = os.Getenv("DB_USER")
	conf.Storage.Password = os.Getenv("DB_PASSWORD")

	conf.Cache.Password = os.Getenv("REDIS_PASSWORD")

}
