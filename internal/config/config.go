package config

import (
	"fmt"
	"os"
	"time"

	wbf "github.com/wb-go/wbf/config"
)

type Config struct {
	Logger  Logger  `mapstructure:"logger"`
	Server  Server  `mapstructure:"server"`
	Storage Storage `mapstructure:"database"`
	Cache   Cache   `mapstructure:"cache"`
}

type Logger struct {
	Debug  bool   `mapstructure:"debug_mode"`
	LogDir string `mapstructure:"log_directory"`
}

type Server struct {
	Port            string        `mapstructure:"port"`
	ReadTimeout     time.Duration `mapstructure:"read_timeout"`
	WriteTimeout    time.Duration `mapstructure:"write_timeout"`
	MaxHeaderBytes  int           `mapstructure:"max_header_bytes"`
	ShutdownTimeout time.Duration `mapstructure:"shutdown_timeout"`
}

type Storage struct {
	Host               string        `mapstructure:"host"`
	Port               string        `mapstructure:"port"`
	Username           string        `mapstructure:"username"`
	Password           string        `mapstructure:"password"`
	DBName             string        `mapstructure:"dbname"`
	SSLMode            string        `mapstructure:"sslmode"`
	MaxOpenConns       int           `mapstructure:"max_open_conns"`
	MaxIdleConns       int           `mapstructure:"max_idle_conns"`
	ConnMaxLifetime    time.Duration `mapstructure:"conn_max_lifetime"`
	QueryRetryStrategy RetryStrategy `mapstructure:"query_retry_strategy"`
}

type Cache struct {
	Host           string        `mapstructure:"host"`
	Port           string        `mapstructure:"port"`
	Password       string        `mapstructure:"password"`
	MaxMemory      string        `mapstructure:"max_memory"`
	Policy         string        `mapstructure:"policy"`
	RetryStrategy  RetryStrategy `mapstructure:"retry_strategy"`
	ExpirationTime time.Duration `mapstructure:"expiration_time"`
}

type RetryStrategy struct {
	Attempts int           `mapstructure:"attempts"`
	Delay    time.Duration `mapstructure:"delay"`
	Backoff  float64       `mapstructure:"backoff"`
}

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

func loadEnvs(conf *Config) {

	conf.Storage.Username = os.Getenv("DB_USER")
	conf.Storage.Password = os.Getenv("DB_PASSWORD")

	conf.Cache.Password = os.Getenv("REDIS_PASSWORD")

}
