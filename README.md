# Platform Installer

## Getting started

To make it easy for you to get started with GitLab, here's a list of recommended next steps.

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
git remote add origin https://hub.smatflow.net/smatflow-projects/platform-installer.git

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

## Integrate with your tools

- [ ] [Set up project integrations](https://hub.smatflow.net/smatflow-projects/platform-installer/-/settings/integrations)

## Collaborate with your team

- [ ] [Invite team members and collaborators](https://docs.gitlab.com/ee/user/project/members/)
- [ ] [Create a new merge request](https://docs.gitlab.com/ee/user/project/merge_requests/creating_merge_requests.html)
- [ ] [Automatically close issues from merge requests](https://docs.gitlab.com/ee/user/project/issues/managing_issues.html#closing-issues-automatically)
- [ ] [Enable merge request approvals](https://docs.gitlab.com/ee/user/project/merge_requests/approvals/)
- [ ] [Automatically merge when pipeline succeeds](https://docs.gitlab.com/ee/user/project/merge_requests/merge_when_pipeline_succeeds.html)

## Test and Deploy

Use the built-in continuous integration in GitLab.

- [ ] [Get started with GitLab CI/CD](https://docs.gitlab.com/ee/ci/quick_start/index.html)
- [ ] [Analyze your code for known vulnerabilities with Static Application Security Testing(SAST)](https://docs.gitlab.com/ee/user/application_security/sast/)
- [ ] [Deploy to Kubernetes, Amazon EC2, or Amazon ECS using Auto Deploy](https://docs.gitlab.com/ee/topics/autodevops/requirements.html)
- [ ] [Use pull-based deployments for improved Kubernetes management](https://docs.gitlab.com/ee/user/clusters/agent/)
- [ ] [Set up protected environments](https://docs.gitlab.com/ee/ci/environments/protected_environments.html)

---

# Editing this README

When you're ready to make this README your own, just edit this file and use the handy template below (or feel free to structure it however you want - this is just a starting point!). Thank you to [makeareadme.com](https://www.makeareadme.com/) for this template.

## Suggestions for a good README

Every project is different, so consider which of these sections apply to yours. The sections used in the template are suggestions for most open source projects. Also keep in mind that while a README can be too long and detailed, too long is better than too short. If you think your README is too long, consider utilizing another form of documentation rather than cutting out information.

## Usage

Use examples liberally, and show the expected output if you can. It's helpful to have inline the smallest example of usage that you can demonstrate, while providing links to more sophisticated examples if they are too long to reasonably include in the README.

## License

For open source projects, say how it is licensed.
