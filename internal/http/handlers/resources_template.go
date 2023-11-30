package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/icloudeng/platform-installer/internal/database/entities"
	"github.com/icloudeng/platform-installer/internal/resources/db"
	"github.com/icloudeng/platform-installer/internal/resources/utilities"
)

// Create resources from platform template
func (resourcesHandler) CreateResourcesFromTemplate(ctx *gin.Context) {
	var tmpBody resourcesBody
	ctx.ShouldBindBodyWith(&tmpBody, binding.JSON)

	// Platform name must be fulfilled
	if tmpBody.Platform == nil || len(tmpBody.Platform.Name) == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Platform name cannot be empty",
		})
		return
	}

	// Platform name must be fulfilled
	if len(tmpBody.Environment) == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Environment cannot be empty",
		})
		return
	}

	platform_name := tmpBody.Platform.Name
	template := entities.ResourcesTemplateRepository{}.GetByPlatform(platform_name)

	if template.ID == 0 {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "No platform template found: " + platform_name,
		})
		return
	}

	client := &clientBody{}
	if tmpBody.Client != nil {
		client = tmpBody.Client
	}

	body := resourcesBody{
		Client:      client,
		Environment: tmpBody.Environment,
		Domain:      template.Domain.Data(),
		Vm:          template.Vm.Data(),
		Platform:    template.Platform.Data(),
	}

	// Merge body request into template (body)
	ctx.ShouldBindBodyWith(&body, binding.JSON)

	// Dynamic reference
	reference := utilities.Helpers.ConcatenateAndCleanParams(
		body.Platform.Name,
		body.Domain.Subdomain,
		body.Environment,
		body.Domain.Zone,
	)

	body.Ref = reference
	body.Vm.Name = reference

	// Validate request
	if err := ctx.ShouldBindBodyWith(&body, binding.JSON); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Concatenate subdomain with env
	if !utilities.Helpers.IsProdEnv(body.Environment) {
		body.Domain.Subdomain = utilities.Helpers.ConcatenateSubdomain(body.Domain.Subdomain, body.Environment)
	}

	// Create or get client from DB
	clientdb := db.ClientCountry{}.CreateClientCountry(db.ClientCountry{
		CountryName: body.Client.CountryName,
		CountryCode: body.Client.CountryCode,
		ClientName:  body.Client.ClientName,
		ClientEmail: body.Client.ClientEmail,
	})

	// Generate VM ID based on client details
	if body.Vm.Vmid == 0 && (tmpBody.Vm == nil || tmpBody.Vm.Vmid == 0) {
		vmid, err := utilities.Helpers.GenerateVMId(
			body.Platform.Name,
			body.Environment,
			int(clientdb.CountryID),
			clientdb.Code,
		)

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		body.Vm.Vmid = vmid
	}

	// Update reference
	if body.Vm.Vmid > 0 {
		reference := utilities.Helpers.ConcatenateAndCleanParams(
			body.Ref,
			strconv.Itoa(body.Vm.Vmid),
		)

		body.Ref = reference
		body.Vm.Name = reference
	}

	job := createResourceJob(ctx, body)

	// if job == nil {
	// 	return
	// }

	ctx.JSON(http.StatusOK, gin.H{"data": body, "job": job})
}
