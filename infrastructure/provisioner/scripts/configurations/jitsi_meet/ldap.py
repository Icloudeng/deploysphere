import sys
import os

# Get the parent directory
parent_dir = os.getcwd()
# # Add the parent directory to sys.path
sys.path.append(parent_dir)

from lib.utilities.auto_configuration import (
    get_resources_state,
    post_provisioning_configuration,
    domain_to_ldap_dc,
    get_command_args,
    log,
)


def main(args):
    # jitsi
    jitsi_reference = args.reference

    # FreeIPA
    freeipa_state = get_resources_state(args.config_reference)["data"]
    ipa_domain = freeipa_state["Job"]["PostBody"]["platform"]["metadata"]["ipa_domain"]
    ipa_domain_dc = domain_to_ldap_dc(ipa_domain)
    freeipa_credentials = freeipa_state["Credentials"][0]
    freeipa_ipv4_address = freeipa_state["State"]["proxmox_vm_qemu"]["values"][
        "default_ipv4_address"
    ]

    body = {
        "ref": jitsi_reference,
        "platform": {
            "name": args.platform,
            "metadata": {
                "configuration_type": args.type,
                "configuration": {
                    "ldap_server_url": "ldap://%s" % freeipa_ipv4_address,
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
