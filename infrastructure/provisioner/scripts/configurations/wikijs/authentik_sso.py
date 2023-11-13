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

    # wikijs
    wikijs_reference = args.reference
    wikijs_state = get_resources_state(wikijs_reference)["data"]
    wikijs_domain = wikijs_state["State"]["ovh_domain_zone_record"]["values"]
    wikijs_domain = concatenate_domain(
        sub_domain=wikijs_domain["subdomain"],
        root_domain=wikijs_domain["zone"],
    )
    wikijs_credentials = wikijs_state["Credentials"][0]

    body = {
        "ref": wikijs_reference,
        "platform": {
            "name": args.platform,
            "metadata": {
                "wikijs_url": "https://%s" % wikijs_domain,
                "wikijs_admin_username": wikijs_credentials["username"],
                "wikijs_admin_password": wikijs_credentials["password"],
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
