package config

import (
	"fmt"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type ServerMode string

const (
	ServerModeAPI    ServerMode = "api"
	ServerModeWorker ServerMode = "worker"
	ServerModeBoth   ServerMode = "both"
)

type PlaidEnv string

const (
	PlaidEnvSandbox     PlaidEnv = "sandbox"
	PlaidEnvDevelopment PlaidEnv = "development"
	PlaidEnvProduction  PlaidEnv = "production"
)

type LogLevel string

const (
	LogLevelDebug LogLevel = "debug"
	LogLevelInfo  LogLevel = "info"
	LogLevelWarn  LogLevel = "warn"
	LogLevelError LogLevel = "error"
)

type Config struct {
	ServerPort int        `envconfig:"SERVER_PORT" default:"8080"`
	ServerMode ServerMode `envconfig:"SERVER_MODE" default:"both"`

	PostgresDSN string `envconfig:"POSTGRES_DSN" required:"true"`

	RedisAddr     string `envconfig:"REDIS_ADDR" required:"true"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" default:""`
	RedisDB       int    `envconfig:"REDIS_DB" default:"0"`

	KafkaBrokers []string `envconfig:"KAFKA_BROKERS" required:"true"`
	KafkaTopic   string   `envconfig:"KAFKA_TOPIC" default:"sync-relay.transactions.v1"`

	PlaidClientID      string   `envconfig:"PLAID_CLIENT_ID" required:"true"`
	PlaidSecret        string   `envconfig:"PLAID_SECRET" required:"true"`
	PlaidEnv           PlaidEnv `envconfig:"PLAID_ENV" default:"sandbox"`
	PlaidWebhookSecret string   `envconfig:"PLAID_WEBHOOK_SECRET" required:"true"`

	EncryptionKey string `envconfig:"ENCRYPTION_KEY" required:"true"`

	WorkerConcurrency int `envconfig:"WORKER_CONCURRENCY" default:"10"`
	LockTTLSeconds    int `envconfig:"LOCK_TTL_SECONDS" default:"60"`

	LogLevel LogLevel `envconfig:"LOG_LEVEL" default:"info"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

func (c *Config) Validate() error {
	if c.ServerMode != ServerModeAPI && c.ServerMode != ServerModeWorker && c.ServerMode != ServerModeBoth {
		return fmt.Errorf("invalid SERVER_MODE: %s (must be api, worker, or both)", c.ServerMode)
	}

	if c.PlaidEnv != PlaidEnvSandbox && c.PlaidEnv != PlaidEnvDevelopment && c.PlaidEnv != PlaidEnvProduction {
		return fmt.Errorf("invalid PLAID_ENV: %s", c.PlaidEnv)
	}

	if len(c.EncryptionKey) != 64 {
		return fmt.Errorf("ENCRYPTION_KEY must be 64 hex characters (32 bytes)")
	}

	return nil
}

func LoadOrPanic() *Config {
	cfg, err := Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Configuration error: %v\n", err)
		os.Exit(1)
	}
	return cfg
}
