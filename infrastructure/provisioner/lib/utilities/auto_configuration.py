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
    response.raise_for_status()

    return response.json()


def post_provisioning_configuration(body):
    response = requests.post(
        f"{INSTALLER_URL}/provisioning/configuration", headers=headers, json=body
    )
    response.raise_for_status()

    return response.json()


def domain_to_ldap_dc(domain):
    # Remove leading and trailing whitespaces and convert to lowercase
    domain = domain.strip().lower()

    # Split the domain into components
    domain_components = domain.split(".")

    # Prefix each component with "dc="
    ldap_dc_components = ["dc=" + component for component in domain_components]

    # Join the components with commas
    ldap_dc = ",".join(ldap_dc_components)

    return ldap_dc


def log(text: str):
    print(f"%%{text}%%")


def get_command_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("--reference", required=True)
    parser.add_argument("--config-reference", required=True)
    parser.add_argument("--type", required=True)
    parser.add_argument("--platform", required=True)

    return parser.parse_args()
