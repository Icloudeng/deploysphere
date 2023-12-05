package env

import (
	"log"
	"strings"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type config struct {
	// LDAP Auth
	LDAP_AUTH                 bool   `env:"LDAP_AUTH" envDefault:"false"`
	LDAP_SERVER_URL           string `env:"LDAP_SERVER_URL"`
	LDAP_BIND_TEMPLATE        string `env:"LDAP_BIND_TEMPLATE"`
	LDAP_AUTHORIZED_USERNAMES string `env:"LDAP_AUTHORIZED_USERNAMES"`

	// DB
	DB_TYPE        string `env:"DB_TYPE"`
	DB_PG_HOST     string `env:"DB_PG_HOST"`
	DB_PG_PORT     string `env:"DB_PG_PORT"`
	DB_PG_DBNAME   string `env:"DB_PG_DBNAME"`
	DB_PG_USER     string `env:"DB_PG_USER"`
	DB_PG_PASSWORD string `env:"DB_PG_PASSWORD"`
	DB_PG_SSLMODE  string `env:"DB_PG_SSLMODE"`
	DB_PG_TIMEZONE string `env:"DB_PG_TIMEZONE"`

	// Front
	FRONT_PROXY bool   `env:"FRONT_PROXY,required"`
	FRONT_URL   string `env:"FRONT_URL,required"`

	// Redis
	REDIS_URL string `env:"REDIS_URL,required"`

	// Proxmox
	PROXMOX_API_URL  string `env:"PROXMOX_API_URL,required"`
	PROXMOX_USERNAME string `env:"PROXMOX_USERNAME,required"`
	PROXMOX_PASSWORD string `env:"PROXMOX_PASSWORD,required"`

	// Sentry
	SENTRY_DSN string `env:"SENTRY_DSN"`
}

var Config config

func (c config) getLdapAuthorizedUsernames() []string {
	var usernames []string

	for _, v := range strings.Split(c.LDAP_AUTHORIZED_USERNAMES, " ") {
		usernames = append(usernames, strings.TrimSpace(v))
	}

	return usernames
}

func (c config) ExistsLdapAuthorizedUsername(username string) bool {
	usernames := c.getLdapAuthorizedUsernames()

	for _, v := range usernames {
		if v == username {
			return true
		}
	}

	return false
}

func init() {
	// Loading the environment variables from '.env' file.
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("unable to load .env file: %e", err)
	}

	err = env.Parse(&Config) // 👈 Parse environment variables into `Config`
	if err != nil {
		log.Fatalf("unable to parse ennvironment variables: %e", err)
	}
}
