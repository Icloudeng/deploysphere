package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icloudeng/platform-installer/internal/database/entities"
	"github.com/icloudeng/platform-installer/internal/http/validators"
	"github.com/icloudeng/platform-installer/internal/structs"
	"gorm.io/datatypes"
)

// Store resources and apply
type (
	templateBody struct {
		Domain     *structs.DomainZoneRecord `json:"domain" binding:"required,json"`
		Subdomains []string                  `json:"subdomains" binding:"dive,alpha"`
		Vm         *structs.ProxmoxVmQemu    `json:"vm" binding:"required,json"`
		Platform   *structs.Platform         `json:"platform" binding:"required,json"`
	}

	templateURI struct {
		Platform string `uri:"platform" binding:"required"`
	}

	templateHandler struct{}
)

var PlatformsTemplates templateHandler

func (templateHandler) CreateResourcesTemplate(c *gin.Context) {
	body := templateBody{
		Vm: structs.NewProxmoxVmQemu(""),

		Platform: &structs.Platform{Metadata: map[string]interface{}{
			"domain": "",
		}},
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Chech if platform the password corresponse to an existing platform folder
	if !validators.ValidatePlatformMetadata(c, *body.Platform) {
		return
	}

	delete(body.Platform.Metadata, "domain")

	repository := entities.ResourcesTemplateRepository{}

	entity := &entities.ResourcesTemplate{
		PlatformName: body.Platform.Name,
		Subdomains:   datatypes.NewJSONType(body.Subdomains),
		Domain:       datatypes.NewJSONType(body.Domain),
		Vm:           datatypes.NewJSONType(body.Vm),
		Platform:     datatypes.NewJSONType(body.Platform),
	}

	repository.UpdateOrCreate(entity)

	c.JSON(http.StatusOK, gin.H{"data": entity})
}

func (templateHandler) GetResourcesTemplate(c *gin.Context) {
	var uri templateURI

	if err := c.ShouldBindUri(&uri); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": err.Error()})
		return
	}

	repository := entities.ResourcesTemplateRepository{}
	template := repository.GetByPlatform(uri.Platform)

	if template.ID == 0 {
		c.JSON(http.StatusOK, gin.H{"data": nil})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": template})
}
