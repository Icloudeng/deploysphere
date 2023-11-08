package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/icloudeng/platform-installer/internal/database/entities"
	"github.com/icloudeng/platform-installer/internal/resources/provisioning"
	"github.com/icloudeng/platform-installer/internal/resources/terraform"
	"github.com/icloudeng/platform-installer/internal/structs"
)

type (
	autoProvisioning struct {
		Provision struct {
			Type string `json:"type" binding:"required"`

			Url *struct {
				PlatformUrl       string `json:"platform_url" binding:"required_without_all=Reference,fqdn"`
				PlatformConfigUrl string `json:"platform_config_url" binding:"required_without_all=Reference,fqdn"`
			} `json:"url" binding:"required_without_all=Reference"`

			Reference *struct {
				PlatformRef       string `json:"platform_ref" binding:"required_without_all=Url,resourceref"`
				PlatformConfigRef string `json:"platform_config_ref" binding:"required_without_all=Url,resourceref"`
			} `json:"reference" binding:"required_without_all=Url"`
		} `json:"provision" binding:"required"`
	}
)

func platformDomainFromUrl(platformUrl string, platformConfigUrl string) (string, string, error) {
	platformRef, _, err := terraform.Resources.GetOvhDomainZoneFromUrl(platformUrl)
	if err != nil {
		return "", "", errors.Join(err, errors.New("url: "+platformUrl))
	}

	platformConfigRef, _, err := terraform.Resources.GetOvhDomainZoneFromUrl(platformConfigUrl)
	if err != nil {
		return "", "", errors.Join(err, errors.New("url: "+platformUrl))
	}

	return platformRef, platformConfigRef, nil
}

func (provisioningHandler) CreateAutoConfigurationProvisioning(ctx *gin.Context) {
	response_body := &autoProvisioning{}

	if err := ctx.ShouldBindJSON(response_body); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"type":  "BindJSON",
		})
		return
	}

	body := response_body.Provision

	var platformRef string
	var platformConfigRef string

	// Reference bing from terraform resources
	if body.Url != nil {
		platform, platformConfig, err := platformDomainFromUrl(body.Url.PlatformUrl, body.Url.PlatformConfigUrl)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
				"type":  "platformDomainFromUrl: PlatformUrl, PlatformConfigUrl",
			})
			return
		}

		platformRef = platform
		platformConfigRef = platformConfig
	}

	// Reference bing from request body
	if body.Url == nil && body.Reference != nil {
		platformRef = body.Reference.PlatformRef
		platformConfigRef = body.Reference.PlatformConfigRef
	}

	platformName, err := terraform.Resources.GetPlatformNameByReference(platformRef)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"type":  "GetPlatformNameByReference: platformName",
		})
		return
	}

	output, err := provisioning.CreateAutoConfigurationProvisioning(structs.AutoConfiguration{
		Type:              body.Type,
		Platform:          platformName,
		PlatformRef:       platformRef,
		PlatformConfigRef: platformConfigRef,
	})

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"type":  "CreateAutoConfigurationProvisioning",
		})
		return
	}

	jobid_str := provisioning.ExtractDataFromConfigurationOutputCommand(output)

	jobid, err := strconv.ParseUint(jobid_str, 10, 64)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"type":  "ParseUint jobid_str",
		})
		return
	}

	job := entities.JobRepository{}.Get(uint(jobid))

	ctx.JSON(http.StatusOK, gin.H{"data": body, "job": job})
}
