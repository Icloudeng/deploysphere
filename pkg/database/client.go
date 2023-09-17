package database

import "gorm.io/gorm"

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

	dbConn.Last(object, ID)

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

	dbConn.Last(client)

	return client
}

func (ClientRepository) Create(client Client) *Client {
	client_email := Client{
		Country:     client.Country,
		ClientEmail: client.ClientEmail,
	}
	dbConn.Last(&client_email)

	if client_email.ID != 0 {
		return &client_email
	}

	client_country := Client{Country: client.Country}
	dbConn.Last(&client_country)

	if client_country.ID != 0 {
		// Count client by country
		var clients int64
		dbConn.Where(&Client{
			Country: client.Country,
		}).Count(&clients)

		client_country.ClientCode = int(clients + 1)
		client_country.ClientEmail = client.ClientEmail

		dbConn.Save(&client_country)
		return &client_country
	}

	// Count countries
	var countries int64
	dbConn.Model(&Client{}).Select("country, count(ID) as total_country").Group("name").Count(&countries)

	client.CountryCode = int(countries + 1)
	client.ClientCode = 1

	dbConn.Save(&client)

	return &client
}

func init() {
	dbConn.AutoMigrate(&Client{})
}
