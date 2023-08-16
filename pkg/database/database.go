package database

import (
	"fmt"
	"log"
	"smatflow/platform-installer/pkg/env"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MapInterface map[string]interface{}

var dbConn *gorm.DB

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

	if env.EnvConfig.DB_TYPE == "postgres" {
		_db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			PrepareStmt: true,
		})

		if dbConn = _db; err != nil {
			log.Fatal(err.Error())
		}
	} else {
		_db, err := gorm.Open(sqlite.Open(".data.sqlite"), &gorm.Config{
			PrepareStmt: true,
		})

		if dbConn = _db; err != nil {
			log.Fatal(err.Error())
		}
	}
}
