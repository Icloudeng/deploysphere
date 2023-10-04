package entities

import (
	"smatflow/platform-installer/internal/database"

	"gorm.io/gorm"
)

type Country struct {
	gorm.Model
	Name string
	Code string
}

type CountryRepository struct{}

func (CountryRepository) GetByID(ID uint) *Country {
	object := &Country{}

	database.Conn.Last(object, ID)

	return object
}

func (CountryRepository) GetByName(name string) *Country {
	object := &Country{}

	database.Conn.Where(&Country{Name: name}).Last(object)

	return object
}

func (CountryRepository) Create(data *Country) {
	database.Conn.Create(data)
}

func init() {
	database.Conn.AutoMigrate(&Country{})
}
