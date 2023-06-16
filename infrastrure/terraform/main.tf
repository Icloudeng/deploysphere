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

  pm_api_url = var.pm_api_url
}
