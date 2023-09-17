package handlers

import (
	"net/http"
	"smatflow/platform-installer/pkg/database/entities"

	"github.com/gin-gonic/gin"
)

type (
	clientBody struct {
		Country     string `json:"country" binding:"required,alpha"`
		ClientEmail string `json:"client_email" binding:"required,email"`
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

	repository := entities.ClientRepository{}

	client := repository.Create(entities.Client{
		Country:     json.Country,
		ClientEmail: json.ClientEmail,
	})

	c.JSON(http.StatusOK, client)
}

func (clientHandler) GetClient(c *gin.Context) {
	var json clientBody

	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	repository := entities.ClientRepository{}

	client := repository.Get(json.Country, json.ClientEmail)

	c.JSON(http.StatusOK, client)
}
