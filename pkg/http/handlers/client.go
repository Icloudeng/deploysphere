package handlers

import (
	"net/http"
	"smatflow/platform-installer/pkg/resources/db"

	"github.com/gin-gonic/gin"
)

type (
	clientBody struct {
		CountryName string `json:"country_name" binding:"required,alpha,lowercase"`
		CountryCode string `json:"country_code" binding:"omitempty,iso3166_1_alpha2,lowercase"`
		ClientEmail string `json:"client_email" binding:"required,email,lowercase"`
		ClientName  string `json:"client_name" binding:"omitempty,ascii"`
	}

	clientHandler struct{}
)

var Client clientHandler

func (clientHandler) CreateClient(c *gin.Context) {
	var json clientBody

	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var resource db.ClientCountry

	client := resource.CreateClientCountry(db.ClientCountry{
		CountryName: json.CountryName,
		CountryCode: json.CountryCode,
		ClientName:  json.ClientName,
		ClientEmail: json.ClientEmail,
	})

	c.JSON(http.StatusOK, client)
}
