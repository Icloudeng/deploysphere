import re
import argparse
import base64

import requests
from .dotenv import config

INSTALLER_URL = config.get("SERVER_URL")
LDAP_USER_NAME = config.get("PROVISIONING_LDAP_USERNAME")
LDAP_PASSWORD = config.get("PROVISIONING_LDAP_PASSWORD")


credentials = f"{LDAP_USER_NAME}:{LDAP_PASSWORD}"
encoded_credentials = base64.b64encode(credentials.encode("utf-8")).decode("utf-8")


headers = {
    "Authorization": f"Basic {encoded_credentials}",
    "Content-Type": "application/json",
}


def concatenate_domain(sub_domain=None, root_domain=None):
    # Create a list of variables
    variables = [sub_domain, root_domain]

    # Filter out any variables that are None or empty
    filtered_variables = [v for v in variables if v]

    # Join the filtered variables with "-" separator
    domain_result = ".".join(filtered_variables)

    return domain_result


# Get resource state
def get_resources_state(ref: str):
    response = requests.get(
        f"{INSTALLER_URL}/resources-state/{ref.strip()}",
        headers=headers,
    )
    if response.status_code >= 400:
        print(f"Type: get_resources_state, Error: {response.text}")

    response.raise_for_status()

    return response.json()


def post_provisioning_configuration(body):
    response = requests.post(
        f"{INSTALLER_URL}/provisioning/configuration", headers=headers, json=body
    )

    if response.status_code >= 400:
        print(f"Type: post_provisioning_configuration, Error: {response.text}")

    response.raise_for_status()

    return response.json()


def domain_to_ldap_dc(domain: str):
    # Remove leading and trailing whitespaces and convert to lowercase
    domain = domain.strip().lower()

    # Split the domain into components
    domain_components = domain.split(".")

    # Prefix each component with "dc="
    ldap_dc_components = ["dc=" + component for component in domain_components]

    # Join the components with commas
    ldap_dc = ",".join(ldap_dc_components)

    return ldap_dc


def remove_first_segment(domain: str):
    parts = domain.split(".")
    if len(parts) > 1:
        return ".".join(parts[1:])
    else:
        return ""


def freeipa_resources_state(reference: str):
    freeipa_state = get_resources_state(reference)["data"]
    metadata = freeipa_state["Job"]["PostBody"]["platform"]["metadata"]
    ipa_domain = metadata.get("ipa_domain")

    if not ipa_domain:
        ipa_domain = extract_subdomain(metadata.get("domain"))

    ipa_domain_dc = domain_to_ldap_dc(ipa_domain)
    freeipa_credentials = freeipa_state["Credentials"][0]
    freeipa_ipv4_address = freeipa_state["State"]["proxmox_vm_qemu"]["values"][
        "default_ipv4_address"
    ]

    return (ipa_domain_dc, freeipa_credentials, freeipa_ipv4_address)


def extract_root_domain(domain):
    regex = r"(?:[a-zA-Z0-9-]+\.)?([a-zA-Z0-9-]+\.[a-zA-Z0-9-]+)$"
    match = re.search(regex, domain)
    return match.group(1) if match else None


def extract_subdomain(full_domain: str):
    parts = full_domain.split(".")
    if len(parts) > 2:
        return ".".join(parts[-(len(parts) - 1) :])
    return full_domain


def log(text: str):
    print(f"%%{text}%%")


def get_command_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("--reference", required=True)
    parser.add_argument("--config-reference", required=True)
    parser.add_argument("--type", required=True)
    parser.add_argument("--platform", required=True)

    return parser.parse_args()
