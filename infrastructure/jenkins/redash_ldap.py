#!/usr/bin/env python3

import base64
import json
import re
import os
import tempfile
import time

import requests


print(
    """
    =============================================
    # redash freeipa SSO
    =============================================
    """
)


JOB_FETCH_INTERVAL = 5
LAST_LOGS = ""

BUILD_TAG = os.environ.get("BUILD_TAG", "")

INSTALLER_URL = os.environ.get("INSTALLER_URL", "https://installer.homelab.xyz")
LDAP_USER_NAME = os.environ.get("LDAP_USER_NAME", "jenkins")
LDAP_PASSWORD = os.environ.get("LDAP_PASSWORD", "")

if not LDAP_PASSWORD or not LDAP_USER_NAME:
    raise Exception("LDAP_USER_NAME or LDAP_PASSWORD cannot be empty")

credentials = f"{LDAP_USER_NAME}:{LDAP_PASSWORD}"
encoded_credentials = base64.b64encode(credentials.encode("utf-8")).decode("utf-8")


headers = {
    "Authorization": f"Basic {encoded_credentials}",
    "Content-Type": "application/json",
}


def format_to_filename_standard(input_string):
    # Remove leading and trailing whitespace
    input_string = input_string.strip()

    # Replace spaces and other non-alphanumeric characters with hyphens
    formatted_string = re.sub(r"[^\w\s-]", "", input_string)
    formatted_string = re.sub(r"[\s]+", "-", formatted_string)

    return formatted_string


def extract_variables(input_string: str):
    # Initialize an empty dictionary to store the metadata variables and their values
    variable_dict = {}

    # Define a regex pattern to match lines containing "METADATA_" followed by a variable name and its value
    pattern = r"^([A-Za-z_][A-Za-z0-9_]*)=(.*?)$"

    # Use re.finditer to find all matches in the input string
    matches = re.finditer(pattern, input_string, re.MULTILINE)

    # Iterate over the matches and populate the dictionary
    for match in matches:
        variable_name = match.group(1).lower()
        value = match.group(2)
        variable_dict[variable_name] = value.strip()

    return variable_dict


def read_and_extract_downstream_variables():
    # Get the temporary directory path
    temp_dir = tempfile.gettempdir()
    # Specify the filename
    file_name = os.path.join(temp_dir, format_to_filename_standard(BUILD_TAG))

    with open(file_name, "r") as file:
        # Read and print the content of the file
        file_content = file.read()

    return extract_variables(file_content)


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
    print(f"Post Provisioning Configuration: " + json.dumps(response.json(), indent=4))
    response.raise_for_status()

    return response.json()


def wait_job_status(jobid: str) -> bool:
    global LAST_LOGS

    response = requests.get(f"{INSTALLER_URL}/jobs/{jobid}", headers=headers)
    response.raise_for_status()

    job = response.json()
    logs: str = job["Logs"].replace("\\n", "\n")

    status = job["Status"]  # (idle | completed | failed | running)
    done = {"completed": True, "failed": True}

    DIFF_LOGS = logs.replace(LAST_LOGS, "")
    LAST_LOGS = logs

    # Print logs
    print(DIFF_LOGS)

    if status == "failed":
        raise Exception("Job failed")

    return bool(done.get(status, False))


# =============================================================================
# redash Functions
# =============================================================================
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


def main():
    variables = read_and_extract_downstream_variables()

    # redash
    redash_reference = variables["redash_reference"]

    # FreeIPA
    freeipa_state = get_resources_state(variables["freeipa_reference"])["data"]
    ipa_domain = freeipa_state["Job"]["PostBody"]["platform"]["metadata"]["ipa_domain"]
    ipa_domain_dc = domain_to_ldap_dc(ipa_domain)
    freeipa_credentials = freeipa_state["Credentials"][0]
    freeipa_ipv4_address = freeipa_state["State"]["proxmox_vm_qemu"]["values"][
        "default_ipv4_address"
    ]

    body = {
        "ref": redash_reference,
        "platform": {
            "name": "redash",
            "metadata": {
                "configuration_type": "ldap",
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
    jobid = response["job"]["ID"]
    print(f"JobID: {jobid} ... ")

    # Wait for job status
    while True:
        if wait_job_status(jobid=jobid):
            break
        time.sleep(JOB_FETCH_INTERVAL)


if __name__ == "__main__":
    time.sleep(5)
    main()
