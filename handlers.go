package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler top lever
type Handler struct{}

var Handlers = Handler{}

// Structs
type DomainZoneRecord struct {
	Zone      string `json:"zone" binding:"required"`
	Subdomain string `json:"subdomain" binding:"required"`
	Fieldtype string `json:"fieldtype" binding:"required,oneof=A AAAA CNAME DNAME NS MX SPF DKIM DMARC TXT SRV CAA NAPTR LOC SSHFP TLSA"`
	Ttl       int    `json:"ttl" binding:"required,number,gte=60"`
	Target    string `json:"target" binding:"required"`
}

type Provision struct {
	Ref    string            `json:"ref" binding:"required"`
	Domain *DomainZoneRecord `json:"domain" binding:"required,json"`
}

// Handler methods
func (s *Handler) provision(c *gin.Context) {
	var json Provision

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// plan, _ := c.Params.Get("plan")

	c.AsciiJSON(http.StatusOK, &json)
}
