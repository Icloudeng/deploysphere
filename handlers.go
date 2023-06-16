package main

import (
	"context"
	"fmt"
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

	go func() {
		if err := Queue.QueueTask(func(ctx context.Context) error {
			fmt.Println("Process QueueTask - Create domain Record and apply")
			// Working on ovh resource
			resourceData := ResourceJSONData{}
			resourceData.ParseOVHresourcesJSON()

			// Add domain to the resource
			resourceData.GetResource().AddDomainZoneRerord(json.Ref, json.Domain)

			// Execute Terraform
			defer Tf.apply()
			defer fmt.Println("Process QueueTask - Terraform apply")

			// Write resource data
			defer resourceData.WriteOVHresources()

			return nil
		}); err != nil {
			panic(err)
		}
	}()

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

	go func() {
		if err := Queue.QueueTask(func(ctx context.Context) error {
			fmt.Println("Process QueueTask - Delete domain Record and apply")
			// Working on ovh resource
			resourceData := ResourceJSONData{}
			resourceData.ParseOVHresourcesJSON()

			// Add domain to the resource
			resourceData.GetResource().DeleteDomainZoneRerord(data.Ref)

			// Execute Terraform
			defer Tf.apply()
			defer fmt.Println("Process QueueTask - Terraform apply")

			// Write resource data
			defer resourceData.WriteOVHresources()

			return nil
		}); err != nil {
			panic(err)
		}
	}()

	c.AsciiJSON(http.StatusOK, data)
}
