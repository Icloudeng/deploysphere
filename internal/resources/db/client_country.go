package db

import "github.com/icloudeng/platform-installer/internal/database/entities"

type (
	ClientCountry struct {
		CountryName string
		CountryCode string
		ClientName  string
		ClientEmail string
	}
)

func (ClientCountry) CreateClientCountry(body ClientCountry) *entities.Client {
	// Country
	var countryRepo entities.CountryRepository
	country := countryRepo.GetByName(body.CountryName)

	if country.ID == 0 {
		country = &entities.Country{
			Name: body.CountryName,
			Code: body.CountryCode,
		}
		countryRepo.Create(country)
	}

	// Client
	var clientRepo entities.ClientRepository
	client := clientRepo.GetByCountry(body.ClientEmail, country.ID)

	if client.ID != 0 {
		return client
	}

	// clientRepo
	client_count := clientRepo.CountByCountry(country.ID)
	client = &entities.Client{
		Name:      body.ClientName,
		Email:     body.ClientEmail,
		Code:      int(client_count + 1),
		CountryID: country.ID,
		Country:   *country,
	}

	clientRepo.Create(client)

	return client
}
