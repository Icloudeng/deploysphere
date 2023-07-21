package validators

import (
	"encoding/json"
	"fmt"
	"net/http"
	"smatflow/platform-installer/pkg/files"
	"smatflow/platform-installer/pkg/structs"

	"github.com/gin-gonic/gin"
)

func ValidatePlatformMetadata(c *gin.Context, platform structs.Platform) bool {
	if len(platform.Name) > 0 {
		// Check if plaform has provisionner script
		if !files.ExistsProvisionerPlaformReadDir(platform.Name) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Cannot found the correspoding platform"})
			return false
		}

		// Check if platform meta fields exist
		requiredFields := files.ReadPlatformMetadataFields()
		meta := structs.PlatformMetadataFields{}
		json.Unmarshal(requiredFields, &meta)

		metadata := *platform.Metadata

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

func ValidateLdapMetadata(c *gin.Context, platform structs.Platform) bool {
	if len(platform.Name) > 0 {
		// Check if plaform has provisionner script
		if !files.ExistsProvisionerPlaformReadDir(platform.Name) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Cannot found the correspoding platform"})
			return false
		}

		// Check if platform meta fields exist
		requiredFields := files.ReadLdapMetadataFields()
		meta := structs.LdapMetadataFields{}
		json.Unmarshal(requiredFields, &meta)

		metadata := *platform.Metadata

		if values, found := meta[platform.Name]; found {
			// Validate fields
			for _, val := range values.Fields {
				if _, exists := metadata[val]; !exists {
					c.AbortWithStatusJSON(
						http.StatusBadRequest,
						gin.H{"error": fmt.Sprintf("Platform (%s), Metadata field (%s) required", platform.Name, val)},
					)
					return false
				}
			}

			// Validate LDAP Fields
			ldap_type, exists := metadata["ldap"]

			if exists {
				ldap_values, ok := ldap_type.(map[string]interface{})

				if !ok {
					c.AbortWithStatusJSON(
						http.StatusBadRequest,
						gin.H{"error": "LDAP field be an object type"},
					)
					return false
				}

				for _, val := range values.Ldap {
					if _, exists := ldap_values[val]; !exists {
						c.AbortWithStatusJSON(
							http.StatusBadRequest,
							gin.H{"error": fmt.Sprintf("Platform (%s), Ldap Metadata field (%s) required", platform.Name, val)},
						)
						return false
					}
				}
			} else {
				c.AbortWithStatusJSON(
					http.StatusBadRequest,
					gin.H{"error": "Unable to find LDAP fields in metadata object"},
				)
				return false
			}
		}
	}

	return true
}
