package structs

type Platform struct {
	Name     string                 `json:"name"`
	Metadata map[string]interface{} `json:"metadata"`
}

type PlatformMetadataFields map[string][]string

type ConfigurationMetadataFields map[string]struct {
	Fields        []string `json:"fields"`
	Configuration []string `json:"configuration"`
}
