variable "ovh_endpoint" {
  description = "OVH endpoint."
  type        = string
  default     = "ovh-eu"
}


variable "ovh_application_key" {
  description = "OVH Application Key."
  type        = string
  sensitive   = true
}


variable "ovh_application_secret" {
  description = "OVH Application Secret."
  type        = string
  sensitive   = true
}

variable "ovh_consumer_key" {
  description = "OVH Consumer Key."
  type        = string
  sensitive   = true
}


# Proxmox
variable "proxmox_api_url" {
  description = "Proxmox api url."
  type        = string
}

variable "proxmox_api_token_id" {
  description = "Proxmox api token id."
  type        = string
}

variable "proxmox_api_token_secret" {
  description = "Proxmox api token secret."
  type        = string
}
