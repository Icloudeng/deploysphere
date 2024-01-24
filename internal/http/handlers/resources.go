package handlers

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/icloudeng/platform-installer/internal/database/entities"
	"github.com/icloudeng/platform-installer/internal/filesystem"
	"github.com/icloudeng/platform-installer/internal/pubsub"
	"github.com/icloudeng/platform-installer/internal/resources/jobs"
	"github.com/icloudeng/platform-installer/internal/resources/terraform"
	"github.com/icloudeng/platform-installer/internal/structs"
)

// Store resources and apply
type (
	resourcesBody struct {
		Ref           string                    `json:"ref" binding:"required,resourceref"`
		Domain        *structs.DomainZoneRecord `json:"domain" binding:"required,json"`
		Subdomains    []string                  `json:"subdomains" binding:"dive,alpha"`
		Vm            *structs.ProxmoxVmQemu    `json:"vm" binding:"required,json"`
		Platform      *structs.Platform         `json:"platform" binding:"omitempty,json"`
		Client        *clientBody               `json:"client" binding:"omitempty,json"`
		Environment   string                    `json:"environment" binding:"omitempty,alpha"`
		MxDomain      *string                   `json:"mx_domain" binding:"omitempty,fqdn|eq=auto"`
		MxDomainValue *structs.DomainZoneRecord `json:"mx_domain_value" binding:"omitempty,json"`
	}

	resourcesRefUri struct {
		Ref string `uri:"ref" binding:"required,resourceref"`
	}

	resourcesHandler struct{}
)

var Resources resourcesHandler

const ResourceExistsError = `The resource reference already exists. If you plan to create a new resource, 
please use a different resource reference name or use PUT method to update resource.
`

func (resourcesHandler) CreateResources(c *gin.Context) {
	json := resourcesBody{
		Vm:       structs.NewProxmoxVmQemu(""),
		Platform: &structs.Platform{Metadata: map[string]interface{}{}},
	}

	if err := c.ShouldBindJSON(&json); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job := createResourceJob(c, &json)

	if job == nil {
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": json, "job": job})
}

// DELETE resources
func (resourcesHandler) DeleteResources(c *gin.Context) {
	var uri resourcesRefUri

	if err := c.ShouldBindUri(&uri); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err.Error()})
		return
	}

	task := jobs.ResourcesJob{
		Ref:           uri.Ref,
		PostBody:      uri,
		Description:   "Resources deletion",
		ResourceState: false, // Disable on resource deletion
		Handler:       c.Request.URL.String(),
		Method:        c.Request.Method,
		Task: func(ctx context.Context, job *entities.Job) error {
			// Remove resources
			deleteResourcesDomain(uri.Ref)
			terraform.Resources.DeleteProxmoxVmQemuResource(uri.Ref)

			// Terraform Apply changes
			err := terraform.Exec.Apply(true)

			if err == nil {
				pubsub.BusEvent.Publish(pubsub.RESOURCES_NOTIFIER_EVENT, structs.Notifier{
					Status:  "info",
					Details: fmt.Sprintf("Job ID: %d\nRef: %s", job.ID, uri.Ref),
					Logs:    "Resources deleted",
				})
			}

			return err
		},
	}

	job := jobs.ResourcesJobTask(task)

	c.JSON(http.StatusOK, gin.H{"data": uri, "job": job})
}

// Get resources state from terraform
func (resourcesHandler) GetTerraformResourcesState(c *gin.Context) {
	state := terraform.Exec.Show()

	if state == nil {
		c.JSON(http.StatusOK, struct{}{})
		return
	}

	c.JSON(http.StatusOK, state)
}

func (resourcesHandler) GetResources(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"Proxmox": terraform.Resources.GetProxmoxResource(),
		"Ovh":     terraform.Resources.GetOvhResource(),
	})
}

func (resourcesHandler) GetResourcesByReference(c *gin.Context) {
	var uri resourcesRefUri

	if err := c.ShouldBindUri(&uri); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"Vm":     terraform.Resources.GetProxmoxVmQemuResource(uri.Ref),
		"Domain": terraform.Resources.GetOvhDomainZoneResource(uri.Ref),
		"Ref":    uri.Ref,
	})
}

func (resourcesHandler) GetPlatforms(c *gin.Context) {
	c.JSON(http.StatusOK, filesystem.ReadProvisionerPlatforms())
}
