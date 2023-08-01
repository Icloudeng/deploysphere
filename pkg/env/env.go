package env

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	LDAP_AUTH          bool   `env:"LDAP_AUTH" envDefault:"false"`
	LDAP_SERVER_URL    string `env:"LDAP_SERVER_URL"`
	LDAP_BIND_TEMPLATE string `env:"LDAP_BIND_TEMPLATE"`
	// DB
	DB_TYPE        string `env:"DB_TYPE"`
	DB_PG_HOST     string `env:"DB_PG_HOST"`
	DB_PG_PORT     string `env:"DB_PG_PORT"`
	DB_PG_DBNAME   string `env:"DB_PG_DBNAME"`
	DB_PG_USER     string `env:"DB_PG_USER"`
	DB_PG_PASSWORD string `env:"DB_PG_PASSWORD"`
	DB_PG_SSLMODE  string `env:"DB_PG_SSLMODE"`
	DB_PG_TIMEZONE string `env:"DB_PG_TIMEZONE"`
}

var EnvConfig Config

func init() {
	// Loading the environment variables from '.env' file.
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("unable to load .env file: %e", err)
	}

	err = env.Parse(&EnvConfig) // 👈 Parse environment variables into `Config`
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}
}
