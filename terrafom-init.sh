set -u
: "${GITLAB_ACCESS_TOKEN}"

statefile="platform-installer-dev"

bin/terraform -chdir=infrastrure/terraform init \
    -backend-config="address=https://hub.smatflow.net/api/v4/projects/20/terraform/state/$statefile" \
    -backend-config="lock_address=https://hub.smatflow.net/api/v4/projects/20/terraform/state/$statefile/lock" \
    -backend-config="unlock_address=https://hub.smatflow.net/api/v4/projects/20/terraform/state/$statefile/lock" \
    -backend-config="username=paradoxe.ngwasi" \
    -backend-config="password=$GITLAB_ACCESS_TOKEN" \
    -backend-config="lock_method=POST" \
    -backend-config="unlock_method=DELETE" \
    -backend-config="retry_wait_min=5"