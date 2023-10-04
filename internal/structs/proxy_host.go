package structs

type ProxyHost struct {
	Domain   string `json:"domain" binding:"required,fqdn"`
	Platform string `json:"platform" binding:"required"`
	Hostname string `json:"hostname" binding:"required,hostname|ipv4|fqdn"`
}

type ProxyHostDomain struct {
	Domain string `json:"domain" binding:"required,fqdn"`
}
