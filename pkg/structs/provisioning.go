package structs

type Provisioning struct {
	MachineUser string    `json:"machine_user" binding:"required"`
	MachineIp   string    `json:"machine_ip" binding:"required,ipv4"`
	Platform    *Platform `json:"platform" binding:"required"`
}
