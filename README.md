# Platform Installer

## Getting started

### Usage

Infrastructure and auto platform provisioning, POST `/resources`

```js
{
  // Terraform Resource reference (Can be anything following this format: /^[a-z]+([0-9a-z]+(?:-[0-9a-z]+)?)*$/ )
  "ref": "wordpress-resource-infra",
  //  Domain
  "domain": {
    "zone": "domain.xyz", // OVH Domain Zone
    "subdomain": "wordpress",
    "fieldtype": "A", // A, CNAME etc.
    "ttl": 3600,
    "target": "x.x.x.x" // Target IP
  },
  "vm": {
    "name": "wordpress-resource-vm",
    "target_node": "promox-srv2", // Proxmox Node
    "clone": "ubuntu-22.04-cloudinit-template", // Proxmox Cloud image template (Must exists)
    "vmid": 0,
    "memory": 4096,
    "network": [
      {
        "bridge": "vmbr1",
        "tag": 10
      }
    ]
  },
  "platform": {
    "name": "wordpress"
  }
}
```

### 1. Prerequisite

1. You must have a Linux environment, preferably `Ubuntu`
2. Make sure the following programs are installed on your computer:

   - Golang, [Docs\*](https://go.dev/dl/)
   - NodeJS [Docs\*](https://nodejs.org/en/download/current) (with PNPM [Docs\*](https://pnpm.io/installation))
   - Python (Version 3.8 and up), [Docs\*](https://www.python.org/downloads/)
   - Bash shell, [Docs\*](https://www.gnu.org/software/bash/)

3. You must have also Redis installed, [Docs\*]()

### 2. Clone the Project

```
git remote add origin https://github.com/Icloudeng/platform-installer.git
cd platform-installer/
```

### 3. Configurations

#### 1. Create `.env` environment file from template

```
cp .env.example .env
```

#### 2. Write env environment variables

The following env variables are required

```bash
ADMIN_SYSTEM_EMAIL=

# Nginx Proxy Manager Interface
NGINX_PM_URL=
NGINX_PM_EMAIL=
NGINX_PM_PASSWORD=

# Redis
REDIS_URL=

# Proxmox
PROXMOX_API_URL=
PROXMOX_USERNAME=
PROXMOX_PASSWORD=
```

Check all variables here: `.env.example`

#### 3. Write Required Terraform Variables

- Create variables file

```
touch infrastructure/terraform/variables.tfvars
```

- Please check all required variables in `infrastructure/terraform/variables.tf`

- Fill the variables inside `infrastructure/terraform/variables.tfvars` using the following format:

```bash
ovh_endpoint           = "**"
ovh_application_key    = "**"
ovh_application_secret = "**"
ovh_consumer_key       = "**"

# Proxmox
proxmox_api_url          = "**"
proxmox_api_token_id     = "**"
proxmox_api_token_secret = "**"
```

## Supported Platforms

Check all supported platforms here: `infrastructure/provisioner/scripts/platforms`

## Test and Deploy

...

## License

...
