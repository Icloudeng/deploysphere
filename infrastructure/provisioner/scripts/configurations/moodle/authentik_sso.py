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

    # moodle
    moodle_reference = args.reference
    moodle_state = get_resources_state(moodle_reference)["data"]
    moodle_domain = moodle_state["State"]["ovh_domain_zone_record"]["values"]
    moodle_domain = concatenate_domain(
        sub_domain=moodle_domain["subdomain"],
        root_domain=moodle_domain["zone"],
    )
    moodle_credentials = moodle_state["Credentials"][0]

    body = {
        "ref": moodle_reference,
        "platform": {
            "name": args.platform,
            "metadata": {
                "moodle_url": "https://%s" % moodle_domain,
                "moodle_admin_username": moodle_credentials["username"],
                "moodle_admin_password": moodle_credentials["password"],
                "configuration_type": args.type,
                "configuration": {
                    "authentik_url": "https://%s" % authentik_domain,
                    "authentik_admin": authentik_credentials["username"],
                    "authentik_password": authentik_credentials["password"],
                },
            },
        },
    }
    post_provisioning_configuration(body)


if __name__ == "__main__":
    args = get_command_args()
    main(args)
