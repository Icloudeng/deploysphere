package database

import (
	"fmt"
	"log"

	"github.com/icloudeng/platform-installer/internal/env"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type MapInterface map[string]interface{}

var Conn *gorm.DB

func init() {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		env.Config.DB_PG_HOST,
		env.Config.DB_PG_USER,
		env.Config.DB_PG_PASSWORD,
		env.Config.DB_PG_DBNAME,
		env.Config.DB_PG_PORT,
		env.Config.DB_PG_SSLMODE,
		env.Config.DB_PG_TIMEZONE,
	)

	if env.Config.DB_TYPE == "postgres" {
		_db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
			PrepareStmt: true,
		})

		if Conn = _db; err != nil {
			log.Fatalln(err.Error(), "DSN: ", dsn)
		}
	} else {
		_db, err := gorm.Open(sqlite.Open(".data.sqlite"), &gorm.Config{
			PrepareStmt: true,
		})

		if Conn = _db; err != nil {
			log.Fatalln(err.Error())
		}
	}
}
