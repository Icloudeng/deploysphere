package structs

type DomainZoneRecord struct {
	Zone      string `json:"zone" binding:"required"`
	Subdomain string `json:"subdomain" binding:"required,alpha,lowercase"`
	Fieldtype string `json:"fieldtype" binding:"required,oneof=A AAAA CNAME DNAME NS MX SPF DKIM DMARC TXT SRV CAA NAPTR LOC SSHFP TLSA"`
	Ttl       int    `json:"ttl" binding:"required,number,gte=60"`
	Target    string `json:"target" binding:"required"`
}

type PmVmQemuNetwork struct{}

type ProxmoxVmQemu struct {
	Name       string `json:"name" binding:"required"`
	TargetNode string `json:"target_node" binding:"required"`

	Vmid        int    `json:"vmid" binding:"number"`
	Description string `json:"desc"`
	Clone       string `json:"clone"`
	FullClone   bool   `json:"full_clone"`
	OsType      string `json:"os_type"`
	OnBoot      bool   `json:"onboot"`
	Agent       int    `json:"agent" binding:"number"`
	Memory      int    `json:"memory" binding:"number"`

	Network *PmVmQemuNetwork `json:"network" binding:"json"`
}
