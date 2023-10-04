# Platform Installer

## Getting started

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
```

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

## Test and Deploy

...

## Usage

...

## License

...
