package structs

type Environments struct {
	Environments map[string]int `json:"environments"`
	Platforms    map[string]int `json:"platforms"`
}
