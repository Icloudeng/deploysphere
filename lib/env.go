package lib

import (
	"fmt"
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

type Config struct {
	LdapServerUrl    string `env:"LDAP_SERVER_URL,required"`
	LdapBindTemplate string `env:"LDAP_BIND_TEMPLATE,required"`
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

	fmt.Println("Config:")
	fmt.Printf("Ldap Server Url: %s\n", EnvConfig.LdapServerUrl)
	fmt.Printf("Ldap Bind Template: %s\n", EnvConfig.LdapBindTemplate)
}
