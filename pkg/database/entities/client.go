package entities

import (
	"smatflow/platform-installer/pkg/database"

	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Name      string
	Email     int
	CountryID uint
	Country   Country
}

type ClientRepository struct{}

func (ClientRepository) GetByID(ID uint) *Client {
	object := &Client{}
	database.Conn.Last(object, ID)

	return object
}

func (ClientRepository) Get(country string, clientEmail string) *Client {
	object := &Client{}
	// database.Conn.Last(object, ID)

	return object
}

func (ClientRepository) Create(client Client) *Client {
	object := &Client{}
	// database.Conn.Last(object, ID)

	return object
}

func init() {
	database.Conn.AutoMigrate(&Client{})
}
