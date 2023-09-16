package database

import "gorm.io/gorm"

type Client struct {
	gorm.Model
	Country     string
	CountryCode int
	ClientEmail string
	ClientCode  int
}

type ClientRepository struct{}

func (ClientRepository) Get(ID uint) *Client {
	object := &Client{}

	dbConn.Last(object, ID)

	if object.ID == 0 {
		return nil
	}

	return object
}

func (ClientRepository) Create(client *Client) {
	dbConn.Create(client)
}

func init() {
	dbConn.AutoMigrate(&Client{})
}
