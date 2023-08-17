set -u
: "${GITLAB_ACCESS_TOKEN}"

set -u
: "${GITLAB_USERNAME}"

set -u
: "${TERRAFORM_STATE_FILENAME}"

set -u
: "${GITLAB_PROJECT_PATH}" # https://gitlab.example.net/api/v4/projects/20

bin/terraform -chdir=infrastructure/terraform init \
    -backend-config="address=$GITLAB_PROJECT_PATH/terraform/state/$TERRAFORM_STATE_FILENAME" \
    -backend-config="lock_address=$GITLAB_PROJECT_PATH/terraform/state/$TERRAFORM_STATE_FILENAME/lock" \
    -backend-config="unlock_address=$GITLAB_PROJECT_PATH/terraform/state/$TERRAFORM_STATE_FILENAME/lock" \
    -backend-config="username=$GITLAB_USERNAME" \
    -backend-config="password=$GITLAB_ACCESS_TOKEN" \
    -backend-config="lock_method=POST" \
    -backend-config="unlock_method=DELETE" \
    -backend-config="retry_wait_min=5"
