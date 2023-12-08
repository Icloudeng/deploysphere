import sys
import os

# Get the parent directory
parent_dir = os.getcwd()
# # Add the parent directory to sys.path
sys.path.append(parent_dir)

from lib.utilities.auto_configuration import (
    freeipa_resources_state,
    post_provisioning_configuration,
    get_command_args,
    log,
)


def main(args):
    # redash
    redash_reference = args.reference

    # FreeIPA
    ipa_domain_dc, freeipa_credentials, freeipa_ipv4_address = freeipa_resources_state(
        args.config_reference
    )

    body = {
        "ref": redash_reference,
        "platform": {
            "name": args.platform,
            "metadata": {
                "configuration_type": args.type,
                "configuration": {
                    "ldap_server_url": "ldap://%s:389" % freeipa_ipv4_address,
                    "ldap_search_base": "cn=accounts,%s" % ipa_domain_dc,
                    "ldap_bind_dn": "uid=admin,cn=users,cn=accounts,%s" % ipa_domain_dc,
                    "ldap_bind_password": freeipa_credentials["password"],
                },
            },
        },
    }

    response = post_provisioning_configuration(body)
    log(response["job"]["ID"])


if __name__ == "__main__":
    args = get_command_args()
    main(args)
