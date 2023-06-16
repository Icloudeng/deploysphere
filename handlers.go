package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler top lever
type Handler struct{}

var Handlers = Handler{}

/**
Handler methods
**/

// Store provistion and apply
type Provision struct {
	Ref    string            `json:"ref" binding:"required,alpha,lowercase"`
	Domain *DomainZoneRecord `json:"domain" binding:"required,json"`
}

func (s *Handler) provision(c *gin.Context) {
	var json Provision

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Working on ovh resource
	resourceData := ResourceJSONData{}
	resourceData.ParseOVHresourcesJSON()
	// Write resource data
	defer resourceData.WriteOVHresources()

	// Add domain to the resource
	resourceData.GetResource().AddDomainZoneRerord(json.Ref, json.Domain)

	c.AsciiJSON(http.StatusOK, json)
}

// Delete provistion and apply
type ProvisionRef struct {
	Ref string `uri:"ref" binding:"required,alpha,lowercase"`
}

func (s *Handler) deleteProvision(c *gin.Context) {
	var data ProvisionRef
	if err := c.ShouldBindUri(&data); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	// Working on ovh resource
	resourceData := ResourceJSONData{}
	resourceData.ParseOVHresourcesJSON()
	// Write resource data
	defer resourceData.WriteOVHresources()

	// Add domain to the resource
	resourceData.GetResource().DeleteDomainZoneRerord(data.Ref)

	c.AsciiJSON(http.StatusOK, data)
}
