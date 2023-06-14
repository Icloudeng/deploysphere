terraform {
  required_providers {
    proxmox = {
      source  = "telmate/proxmox"
      version = "v2.9.14"
    }
  }
}

provider "proxmox" {
  pm_api_url    = var.pm_api_url
  pm_log_enable = true
  pm_log_file   = "terraform-plugin-proxmox.log"
  pm_debug      = true
  pm_log_levels = {
    _default    = "debug"
    _capturelog = ""
  }
}
