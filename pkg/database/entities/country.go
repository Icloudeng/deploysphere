package entities

import (
	"smatflow/platform-installer/pkg/database"

	"gorm.io/gorm"
)

type Country struct {
	gorm.Model
	Name    string
	Code    string
	Clients []Client
}

type CountryRepository struct{}

func (CountryRepository) GetByID(ID uint) *Country {
	object := &Country{}

	return object
}

func (CountryRepository) Get(name string) *Country {
	object := &Country{}

	return object
}

func (CountryRepository) Create(country *Country) *Country {
	object := &Country{}

	return object
}

func init() {
	database.Conn.AutoMigrate(&Country{})
}
