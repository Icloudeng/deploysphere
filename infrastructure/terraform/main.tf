terraform {
  backend "http" {
  }
}


module "ovh" {
  source = "./modules/ovh"

  endpoint           = var.ovh_endpoint
  application_key    = var.ovh_application_key
  application_secret = var.ovh_application_secret
  consumer_key       = var.ovh_consumer_key
}


module "proxmox" {
  source = "./modules/proxmox"

  pm_api_url          = var.proxmox_api_url
  pm_api_token_id     = var.proxmox_api_token_id
  pm_api_token_secret = var.proxmox_api_token_secret
}
