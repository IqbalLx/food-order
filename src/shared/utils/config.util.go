package utils

import (
	"fmt"
)

type postgresConfig struct {
	ConnStringAvailable bool
	ConnString string
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func NewPostgresConfig(env Env) *postgresConfig{
	connString := env.ReadWithDefaultVal("POSTGRES_CONN_STRING", "")
	if (connString != "") {
		return &postgresConfig{
			ConnStringAvailable: true,
			ConnString: connString,
		}
	}

	return &postgresConfig{
		ConnStringAvailable: false,
		Host: env.Read("POSTGRES_HOST"),
		Port: env.ReadWithDefaultVal("POSTGRES_PORT", "5432"),
		User: env.Read("POSTGRES_USER"),
		Password: env.Read("POSTGRES_PASSWORD"),
		Database: env.ReadWithDefaultVal("POSTGRES_DB_NAME", "food-order-db"),
		SSLMode: env.ReadWithDefaultVal("POSTGRES_SSL_MODE", "disable"),
	}
}

func (cfg *postgresConfig) DSNString() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode)
}

type AppConfig struct {
	Database postgresConfig
	Name string
	Environment string
	Address string
}

func NewAppConfig(env Env, postgresConfig *postgresConfig) *AppConfig {
	return &AppConfig{
		Database: *postgresConfig,
		Name: env.ReadWithDefaultVal("APP_NAME", "food-order"),
		Environment: env.ReadWithDefaultVal("APP_ENV", "dev"),
		Address: env.ReadWithDefaultVal("APP_ADDRESS", ":3000"),
	}
}