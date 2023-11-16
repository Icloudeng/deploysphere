package structs

type DomainZoneRecord struct {
	Zone      string `json:"zone" binding:"required"`
	Subdomain string `json:"subdomain" binding:"omitempty,ascii,lowercase"`
	Fieldtype string `json:"fieldtype" binding:"required,oneof=A AAAA CNAME DNAME NS MX SPF DKIM DMARC TXT SRV CAA NAPTR LOC SSHFP TLSA"`
	Ttl       int    `json:"ttl" binding:"required,number,gte=60"`
	Target    string `json:"target" binding:"required"`
}
