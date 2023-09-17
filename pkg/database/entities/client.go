package entities

import (
	"smatflow/platform-installer/pkg/database"

	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Country     string `gorm:"uniqueIndex"`
	CountryCode int
	ClientEmail string `gorm:"uniqueIndex"`
	ClientCode  int
}

type ClientRepository struct{}

func (ClientRepository) GetByID(ID uint) *Client {
	object := &Client{}

	database.Conn.Last(object, ID)

	if object.ID == 0 {
		return nil
	}

	return object
}

func (ClientRepository) Get(country string, clientEmail string) *Client {
	client := &Client{
		Country:     country,
		ClientEmail: clientEmail,
	}

	database.Conn.Last(client)

	return client
}

func (ClientRepository) Create(client Client) *Client {
	client_email := Client{
		Country:     client.Country,
		ClientEmail: client.ClientEmail,
	}
	database.Conn.Last(&client_email)

	if client_email.ID != 0 {
		return &client_email
	}

	client_country := Client{Country: client.Country}
	database.Conn.Last(&client_country)

	if client_country.ID != 0 {
		// Count client by country
		var clients int64
		database.Conn.Where(&Client{
			Country: client.Country,
		}).Count(&clients)

		client_country.ClientCode = int(clients + 1)
		client_country.ClientEmail = client.ClientEmail

		database.Conn.Save(&client_country)
		return &client_country
	}

	// Count countries
	var countries int64
	database.Conn.Model(&Client{}).Select("country, count(ID) as total_country").Group("name").Count(&countries)

	client.CountryCode = int(countries + 1)
	client.ClientCode = 1

	database.Conn.Save(&client)

	return &client
}

func init() {
	database.Conn.AutoMigrate(&Client{})
}
