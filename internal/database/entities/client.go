package entities

import (
	"smatflow/platform-installer/internal/database"

	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Name      string
	Email     string
	Code      int
	CountryID uint
	Country   Country
}

type ClientRepository struct{}

func (ClientRepository) GetByID(ID uint) *Client {
	object := &Client{}
	database.Conn.Last(object, ID)

	return object
}

func (ClientRepository) GetByEmail(email string) *Client {
	object := &Client{}

	database.Conn.Where(&Client{Email: email}).Last(object)

	return object
}

func (ClientRepository) CountByCountry(countryID uint) int64 {
	var count int64
	database.Conn.Model(&Client{}).Select("count(ID)").Where(&Client{CountryID: countryID}).Scan(&count)

	return count
}

func (ClientRepository) GetByCountry(email string, countryID uint) *Client {
	var object Client
	database.Conn.Joins("Country").Where(&Client{CountryID: countryID, Email: email}).Last(&object)

	return &object
}

func (ClientRepository) Create(client *Client) {
	database.Conn.Create(client)
}

func init() {
	database.Conn.AutoMigrate(&Client{})
}
