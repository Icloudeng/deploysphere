package validators

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smatflow/platform-installer/pkg/filesystem"
	"smatflow/platform-installer/pkg/resources/terraform"
	"smatflow/platform-installer/pkg/structs"

	"github.com/gin-gonic/gin"
)

func ValidatePlatformMetadata(c *gin.Context, platform structs.Platform) bool {
	if len(platform.Name) > 0 {
		// Check if plaform has provisionner script
		if !filesystem.ExistsProvisionerPlaformReadDir(platform.Name) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Cannot found the correspoding platform"})
			return false
		}

		// Check if platform meta fields exist
		required_fields := filesystem.ReadPlatformMetadataFields()
		meta := structs.PlatformMetadataFields{}
		json.Unmarshal(required_fields, &meta)

		metadata := platform.Metadata

		if values, found := meta[platform.Name]; found {
			for _, val := range values {
				if _, exists := metadata[val]; !exists {
					c.AbortWithStatusJSON(
						http.StatusBadRequest,
						gin.H{"error": fmt.Sprintf("platform (%s), Metadata field (%s) required", platform.Name, val)},
					)
					return false
				}
			}
		}
	}

	return true
}

func ValidateConfigurationMetadata(c *gin.Context, platform structs.Platform) bool {
	if len(platform.Name) > 0 {
		// Check if plaform has provisionner script
		if !filesystem.ExistsProvisionerPlaformReadDir(platform.Name) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Cannot found the correspoding platform"})
			return false
		}

		// Check if platform meta fields exist
		required_fields := filesystem.ReadConfigurationMetadataFields()
		metadata_struct := structs.ConfigurationMetadataFields{}
		json.Unmarshal(required_fields, &metadata_struct)

		if platform_config, found := metadata_struct[platform.Name]; found {
			request_metadata := platform.Metadata

			// Check Configuration type provided
			configuration_type, exists := request_metadata["configuration_type"]
			config_type, err := configuration_type.(string)

			if !exists || err {
				c.AbortWithStatusJSON(
					http.StatusBadRequest,
					gin.H{"error": "Configuration Type must be provided"},
				)
				return false
			}

			object, exists := platform_config[config_type]
			if !exists {
				c.AbortWithStatusJSON(
					http.StatusBadRequest,
					gin.H{"error": "Unable to find the corresponding configuration type"},
				)
				return false
			}

			// Validate platform fields
			for _, val := range object.Fields {
				if _, exists := request_metadata[val]; !exists {
					c.AbortWithStatusJSON(
						http.StatusBadRequest,
						gin.H{"error": fmt.Sprintf("Platform (%s), Metadata field (%s) required", platform.Name, val)},
					)
					return false
				}
			}

			// Validate Configuration Fields
			configuration, exists := request_metadata["configuration"]
			if exists {
				config_values, ok := configuration.(map[string]interface{})
				if !ok {
					c.AbortWithStatusJSON(
						http.StatusBadRequest,
						gin.H{"error": "Configuration field be an object type"},
					)
					return false
				}

				for _, val := range object.Configuration {
					if _, exists := config_values[val]; !exists {
						c.AbortWithStatusJSON(
							http.StatusBadRequest,
							gin.H{"error": fmt.Sprintf("Platform (%s), Configuration Metadata field (%s) required", platform.Name, val)},
						)
						return false
					}
				}
			} else {
				c.AbortWithStatusJSON(
					http.StatusBadRequest,
					gin.H{"error": "Unable to find Configuration field in metadata object"},
				)
				return false
			}
		}
	}

	return true
}

func ValidatePlatformProvisionAndBindResourceState(body *structs.Provisioning) bool {
	if len(body.Ref) > 0 {
		resourceState := terraform.ResourceState{Module: "proxmox"}
		vm_resource := resourceState.GetResourceState(body.Ref)

		if vm_resource != nil {
			values := vm_resource.AttributeValues

			muser, muser_ok := values["ciuser"]
			mip, mip_ok := values["default_ipv4_address"]

			if mip_ok && muser_ok {
				body.MachineUser = muser.(string)
				body.MachineIp = mip.(string)
			}
		}
	}

	if len(body.MachineUser) == 0 || len(body.MachineIp) == 0 {
		return false
	}

	return true
}
