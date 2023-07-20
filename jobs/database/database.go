package database

import (
	"fmt"
	"log"
	"smatflow/platform-installer/pkg/env"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func init() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		env.EnvConfig.DB_PG_HOST,
		env.EnvConfig.DB_PG_USER,
		env.EnvConfig.DB_PG_PASSWORD,
		env.EnvConfig.DB_PG_DBNAME,
		env.EnvConfig.DB_PG_PORT,
		env.EnvConfig.DB_PG_SSLMODE,
		env.EnvConfig.DB_PG_TIMEZONE,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if Db = db; err != nil {
		log.Print(err.Error())
	}
}
