package structs

type Provisioning struct {
	Ref         string    `json:"ref" binding:"required_without_all=MachineUser MachineIp"`
	MachineUser string    `json:"machine_user" binding:"required_without=Ref"`
	MachineIp   string    `json:"machine_ip" binding:"required_without=Ref"`
	Platform    *Platform `json:"platform" binding:"required"`
}

type AutoConfiguration struct {
	Type              string
	PlatformRef       string
	PlatformConfigRef string
	Platform          string
}
