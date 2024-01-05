import sys
import os

# Get the parent directory
parent_dir = os.getcwd()
# # Add the parent directory to sys.path
sys.path.append(parent_dir)

from lib.utilities.auto_configuration import (
    post_provisioning_configuration,
    get_command_args,
    log,
    freeipa_resources_state,
    get_resources_state,
    concatenate_domain,
)


def main(args):
    # Limesurvey
    limesurvey_reference = args.reference
    limesurvey_state = get_resources_state(limesurvey_reference)["data"]
    limesurvey_domain = limesurvey_state["State"]["ovh_domain_zone_record"]["values"]
    limesurvey_domain = concatenate_domain(
        sub_domain=limesurvey_domain["subdomain"],
        root_domain=limesurvey_domain["zone"],
    )
    limesurvey_credentials = limesurvey_state["Credentials"][0]

    # FreeIPA
    ipa_domain_dc, freeipa_credentials, freeipa_ipv4_address = freeipa_resources_state(args.config_reference)

    body = {
        "ref": limesurvey_reference,
        "platform": {
            "name": args.platform,
            "metadata": {
                "limesurvey_url": "https://%s" % limesurvey_domain,
                "limesurvey_admin_username": limesurvey_credentials["username"],
                "limesurvey_admin_password": limesurvey_credentials["password"],
                "configuration_type": args.type,
                "configuration": {
                    "ldap_server_host": "ldap://%s" % freeipa_ipv4_address,
                    "ldap_server_port": "389",
                    "ldap_search_base": "cn=users,cn=accounts,%s" % ipa_domain_dc,
                },
            },
        },
    }

    response = post_provisioning_configuration(body)
    log(response["job"]["ID"])


if __name__ == "__main__":
    args = get_command_args()
    main(args)
