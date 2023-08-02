package structs

type Provisioning struct {
	Ref         string    `json:"ref" binding:"required_without_all=machine_user machine_ip"`
	MachineUser string    `json:"machine_user" binding:"required_without=ref"`
	MachineIp   string    `json:"machine_ip" binding:"required_without=ref,ipv4"`
	Platform    *Platform `json:"platform" binding:"required"`
}
