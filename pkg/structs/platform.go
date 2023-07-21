package structs

type Platform struct {
	Name     string                  `json:"name"`
	Metadata *map[string]interface{} `json:"metadata"`
}

type PlatformMetadataFields map[string][]string

type LdapMetadataFields map[string]struct {
	Fields []string `json:"fields"`
	Ldap   []string `json:"ldap"`
}
