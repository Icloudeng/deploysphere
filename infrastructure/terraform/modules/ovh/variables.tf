variable "endpoint" {
  description = "OVH endpoint."
  type        = string
  default     = "ovh-eu"
}


variable "application_key" {
  description = "OVH Application Key."
  type        = string
  sensitive   = true
}


variable "application_secret" {
  description = "OVH Application Secret."
  type        = string
  sensitive   = true
}

variable "consumer_key" {
  description = "OVH Consumer Key."
  type        = string
  sensitive   = true
}
