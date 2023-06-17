package structs

type DomainZoneRecord struct {
	Zone      string `json:"zone" binding:"required"`
	Subdomain string `json:"subdomain" binding:"required,alpha,lowercase"`
	Fieldtype string `json:"fieldtype" binding:"required,oneof=A AAAA CNAME DNAME NS MX SPF DKIM DMARC TXT SRV CAA NAPTR LOC SSHFP TLSA"`
	Ttl       int    `json:"ttl" binding:"required,number,gte=60"`
	Target    string `json:"target" binding:"required"`
}

type PmVmQemuNetwork struct {
	Bridge string `json:"bridge" binding:"required"`
	Model  string `json:"model" binding:"required,oneof=e1000 e1000-82540em e1000-82544gc e1000-82545em i82551 i82557b i82559er ne2k_isa ne2k_pci pcnet rtl8139 virtio vmxnet3"`
}

type ProxmoxVmQemu struct {
	Name        string `json:"name" binding:"required"`
	TargetNode  string `json:"target_node" binding:"required"`
	Vmid        int    `json:"vmid" binding:"number"`
	Description string `json:"desc"`
	Clone       string `json:"clone" binding:"required"`

	FullClone bool `json:"full_clone" binding:"boolean"`

	OsType  string            `json:"os_type" binding:"oneof=ubuntu centos cloud-init"`
	OnBoot  bool              `json:"onboot"  binding:"boolean"`
	Agent   int               `json:"agent"   binding:"number,oneof=1 0"`
	Memory  int               `json:"memory"  binding:"number,min=1024"` //RAM
	Scsihw  string            `json:"scsihw"  binding:"oneof=lsi lsi53c810 megasas pvscsi virtio-scsi-pci virtio-scsi-single"`
	Network []PmVmQemuNetwork `json:"network"`

	Cores   int    `json:"cores" binding:"number,min=1,max=10"`   //CPU
	Sockets int    `json:"sockets" binding:"number,min=1,max=10"` //CPU
	Cpu     string `json:"cpu"`                                   //CPU
}

func NewProxmoxVmQemu() *ProxmoxVmQemu {

	pm := ProxmoxVmQemu{
		Vmid:      0,
		FullClone: true,
		OsType:    "cloud-init",
		OnBoot:    true,
		Agent:     1,
		Memory:    1024,
		Scsihw:    "virtio-scsi-pci",
		Cores:     1,
		Sockets:   1,
		Cpu:       "host",
	}

	pm.Network = append(pm.Network, PmVmQemuNetwork{
		Bridge: "vmbr0",
		Model:  "virtio",
	})

	return &pm
}
