package structs

type Platform struct {
	Name     string       `json:"name" binding:"required"`
	Metadata *interface{} `json:"metadata" binding:"json"`
}
