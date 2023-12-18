import sys
import os

# Get the parent directory
parent_dir = os.getcwd()
# # Add the parent directory to sys.path
sys.path.append(parent_dir)

from lib.utilities.auto_configuration import (
    get_resources_state,
    concatenate_domain,
    post_provisioning_configuration,
    get_command_args,
    log,
)


def main(args):
    # Authentik
    authentik_state = get_resources_state(args.config_reference)["data"]
    authentik_domain = authentik_state["State"]["ovh_domain_zone_record"]["values"]
    authentik_domain = concatenate_domain(
        sub_domain=authentik_domain["subdomain"],
        root_domain=authentik_domain["zone"],
    )
    authentik_credentials = authentik_state["Credentials"][0]

    # discourse
    discourse_reference = args.reference
    discourse_state = get_resources_state(discourse_reference)["data"]
    discourse_domain = discourse_state["State"]["ovh_domain_zone_record"]["values"]
    discourse_domain = concatenate_domain(
        sub_domain=discourse_domain["subdomain"],
        root_domain=discourse_domain["zone"],
    )
    discourse_credentials = discourse_state["Credentials"][0]

    body = {
        "ref": discourse_reference,
        "platform": {
            "name": args.platform,
            "metadata": {
                "discourse_url": "https://%s" % discourse_domain,
                "discourse_admin_username": discourse_credentials["username"],
                "discourse_admin_password": discourse_credentials["password"],
                "configuration_type": args.type,
                "configuration": {
                    "authentik_url": "https://%s" % authentik_domain,
                    "authentik_admin": authentik_credentials["username"],
                    "authentik_password": authentik_credentials["password"],
                },
            },
        },
    }

    response = post_provisioning_configuration(body)
    log(response["job"]["ID"])


if __name__ == "__main__":
    args = get_command_args()
    main(args)
