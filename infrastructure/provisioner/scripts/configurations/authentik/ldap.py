import sys
import os

# Get the parent directory
parent_dir = os.getcwd()
# # Add the parent directory to sys.path
sys.path.append(parent_dir)

from lib.utilities.auto_configuration import (
    get_resources_state,
    concatenate_domain,
    freeipa_resources_state,
    post_provisioning_configuration,
    get_command_args,
    log,
)


def main(args):
    # Authentik
    authentik_reference = args.reference
    authentik_state = get_resources_state(authentik_reference)["data"]
    authentik_domain = authentik_state["State"]["ovh_domain_zone_record"]["values"]
    authentik_domain = concatenate_domain(
        sub_domain=authentik_domain["subdomain"],
        root_domain=authentik_domain["zone"],
    )
    authentik_credentials = authentik_state["Credentials"][0]

    # FreeIPA
    ipa_domain_dc, freeipa_credentials, freeipa_ipv4_address = freeipa_resources_state(
        args.config_reference
    )

    body = {
        "ref": authentik_reference,
        "platform": {
            "name": args.platform,
            "metadata": {
                "authentik_url": "https://%s" % authentik_domain,
                "authentik_admin": authentik_credentials["username"],
                "authentik_password": authentik_credentials["password"],
                "configuration_type": args.type,
                "configuration": {
                    "ldap_bind_dn": "uid=admin,cn=users,cn=accounts,%s" % ipa_domain_dc,
                    "ldap_bind_password": freeipa_credentials["password"],
                    "ldap_search_base": "cn=accounts,%s" % ipa_domain_dc,
                    "ldap_server_url": "ldap://%s" % freeipa_ipv4_address,
                },
            },
        },
    }

    response = post_provisioning_configuration(body)
    log(response["job"]["ID"])


if __name__ == "__main__":
    args = get_command_args()
    main(args)
