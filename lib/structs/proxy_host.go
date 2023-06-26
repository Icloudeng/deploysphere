package structs

type ProxyHost struct {
	Platform string `json:"platform" binding:"required"`
	Domain   string `json:"domain" binding:"required"`
	Ip       string `json:"ip" binding:"required"`
}

type ProxyHostDomain struct {
	Domain string `json:"domain" binding:"required"`
}
