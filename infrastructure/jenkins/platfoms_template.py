#!/usr/bin/env python3

import re
import base64
import time
import requests
import os
import json
import tempfile


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


platforms = {
    "freeipa": "10",
    "authentik": "12",
    "keycloak": "13",
    "zimbra": "20",
    "nextcloud": "21",
    "wordpress": "22",
    "discourse": "23",
    "joomla": "24",
    "xwiki": "25",
    "moodle": "30",
    "jupyterhub": "35",
    "erpnext": "40",
    "suitecrm": "41",
    "mautic": "42",
    "jenkins": "70",
    "nifi": "71",
    "prometheus": "75",
    "elastic": "76",
    "postgresql": "80",
    "mariadb": "81",
    "redash": "95",
    "grafana": "96",
    "wikijs": "33",
}


platforms_subdomain = {
    "freeipa": "auth",
    "authentik": "sso",
    "keycloak": "sso",
    "zimbra": "mail",
    "nextcloud": "drive",
    "wordpress": "www",
    "discourse": "forum",
    "joomla": "cms",
    "xwiki": "wiki",
    "wikijs": "wiki",
    "moodle": "lms",
    "jupyterhub": "lab",
    "erpnext": "erp",
    "suitecrm": "crm",
    "mautic": "marketing",
    "jenkins": "cicd",
    "nifi": "etl",
    "prometheus": "monitor",
    "elastic": "indexer",
    "postgresql": "postgres",
    "mariadb": "mysql",
    "redash": "analytics",
    "grafana": "metrics",
}

cloud_init_templates = {
    "freeipa": "ubuntu-20.04-cloudinit-template",
    "authentik": "ubuntu-22.04-cloudinit-template",
    "keycloak": "ubuntu-20.04-cloudinit-template",
    "zimbra": "ubuntu-20.04-cloudinit-template",
    "nextcloud": "ubuntu-20.04-cloudinit-template",
    "wordpress": "ubuntu-22.04-cloudinit-template",
    "discourse": "ubuntu-20.04-cloudinit-template",
    "joomla": "ubuntu-22.04-cloudinit-template",
    "moodle": "ubuntu-22.04-cloudinit-template",
    "jupyterhub": "ubuntu-22.04-cloudinit-template",
    "erpnext": "ubuntu-22.04-cloudinit-template",
    "suitecrm": "ubuntu-22.04-cloudinit-template",
    "mautic": "ubuntu-22.04-cloudinit-template",
    "jenkins": "ubuntu-20.04-cloudinit-template",
    "mysql": "ubuntu-22.04-cloudinit-template",
    "redash": "ubuntu-22.04-cloudinit-template",
}


environments = {
    "prod": "1",
    "stage": "2",
    "test": "3",
    "dev": "4",
    "lab": "5",
    "tpl": "6",
    "incub": "7",
    "ngo": "8",
    "infra": "9",
}


def concatenate_resources(
    platform_name=None, sub_domain=None, env_domain=None, root_domain=None
):
    # Create a list of variables
    variables = [platform_name, sub_domain, env_domain, root_domain]

    # If root_domain has a ".", replace it with "-"
    if root_domain:
        root_domain = root_domain.replace(".", "-")
        variables[
            -1
        ] = root_domain  # Update the last element with the modified root_domain

    # Filter out any variables that are None or empty
    filtered_variables = [v for v in variables if v]

    # Join the filtered variables with "-" separator
    resource_ref = "-".join(filtered_variables)

    return resource_ref


def concatenate_domain(sub_domain=None, env_domain=None, root_domain=None):
    # Create a list of variables
    variables = [sub_domain, env_domain, root_domain]

    # Filter out any variables that are None or empty
    filtered_variables = [v for v in variables if v]

    # Join the filtered variables with "-" separator
    domain_result = ".".join(filtered_variables)

    return domain_result


def concatenate_ldap_schema(env_domain=None, root_domain=None):
    # Create a list of variables
    variables = [env_domain, root_domain]

    # Filter out any variables that are None or empty
    filtered_variables = [v for v in variables if v]

    # Join the filtered variables with "-" separator
    ldap_schema_result = ".".join(filtered_variables)

    return ldap_schema_result


def concatenate_subdomain(sub_domain=None, env_domain_1=None, env_domain_2=None):
    # Create a list of variables
    variables = [sub_domain, env_domain_1, env_domain_2]

    # Filter out any variables that are None or empty
    filtered_variables = [v for v in variables if v]

    # Join the filtered variables with "-" separator
    sub_domain_result = ".".join(filtered_variables)

    return sub_domain_result


def check_empty_fields(data, path=[]):
    for key, value in data.items():
        current_path = path + [key]

        if isinstance(value, dict):
            check_empty_fields(value, current_path)
        elif value is None or (isinstance(value, str) and value.strip() == ""):
            raise ValueError(f"Empty field at path: {'.'.join(current_path)}")


def format_to_filename_standard(input_string):
    # Remove leading and trailing whitespace
    input_string = input_string.strip()

    # Replace spaces and other non-alphanumeric characters with hyphens
    formatted_string = re.sub(r"[^\w\s-]", "", input_string)
    formatted_string = re.sub(r"[\s]+", "-", formatted_string)

    return formatted_string


def add_default_email(data, default_email):
    if isinstance(data, dict):
        for key, value in data.items():
            if "email" in key.lower() and isinstance(value, str):
                # Check if the value looks like an email address
                if len(value.strip()) == 0:
                    # Field appears to be an email; add the default email if it's empty
                    data[key] = default_email
            elif isinstance(value, (dict, list)):
                # Recursively check nested dictionaries and lists
                add_default_email(value, default_email)
    elif isinstance(data, list):
        # Recursively check items in a list
        for item in data:
            add_default_email(item, default_email)


def pad_project_code(PROJECT_CODE, PROJECT_COUNTRY_NUMBER):
    if not PROJECT_CODE.isdigit():
        raise ValueError(f"'{PROJECT_CODE}' is not a valid digit sequence.")

    num_digits = len(str(PROJECT_COUNTRY_NUMBER))

    # Calculate the upper limit based on num_digits
    upper_limit = int("9" * (6 - num_digits))

    code_as_int = int(PROJECT_CODE)

    if not (0 <= code_as_int <= upper_limit):
        raise ValueError(f"'{PROJECT_CODE}' is not in the range [0, {upper_limit}].")

    return PROJECT_CODE.zfill(6 - num_digits)


def get_metadata_variables():
    metadata_variables = {}
    for var_name, var_value in os.environ.items():
        if var_name.startswith("METADATA_"):
            # Remove the "METADATA_" prefix and convert variable name to lowercase
            key = var_name[len("METADATA_") :].lower()
            metadata_variables[key] = var_value
    return metadata_variables


def extract_metadata_variables(input_string):
    # Initialize an empty dictionary to store the metadata variables and their values
    metadata_dict = {}

    # Define a regex pattern to match lines containing "METADATA_" followed by a variable name and its value
    pattern = r"^METADATA_([A-Za-z_][A-Za-z0-9_]*)=(.*?)$"

    # Use re.finditer to find all matches in the input string
    matches = re.finditer(pattern, input_string, re.MULTILINE)

    # Iterate over the matches and populate the dictionary
    for match in matches:
        variable_name = match.group(1).lower()
        value = match.group(2)
        metadata_dict[variable_name] = value.strip()

    return metadata_dict


def is_empty_key(dictionary: dict, key: str):
    if key in dictionary and not dictionary.get(key):
        return True
    else:
        return False


def default_metadata_values_platform(metadata: dict, platfom: str):
    if platfom == "freeipa" and is_empty_key(metadata, "ipa_domain"):
        metadata["ipa_domain"] = DOMAIN_ROOT


def is_memory_binary_format(number):
    # Check if the number is a power of 2
    return number & (number - 1) == 0 and number != 0


def convert_dict_values(data):
    if isinstance(data, dict):
        for key, value in data.items():
            data[key] = convert_dict_values(value)
    elif isinstance(data, str):
        if data.lower() == "true":
            return True
        elif data.lower() == "false":
            return False
        elif data.replace(".", "", 1).isdigit():
            # Check if the string is a valid number (float or int)
            if "." in data:
                return float(data)
            else:
                return int(data)
    return data


def prune_empty(d):
    """Recursively prune dictionary items with None values, empty strings, and empty dictionaries."""
    if not isinstance(d, dict):
        return d
    pruned = {}
    for k, v in d.items():
        if isinstance(v, dict):
            v = prune_empty(v)
            if v:  # Check if the pruned dictionary is non-empty
                pruned[k] = v
        elif v not in (None, "", []):
            pruned[k] = v
    return pruned


def write_jobid_to_upstream_file(
    platform: str, upstream: str, jobid: int, reference: str
):
    if not upstream or not jobid:
        return

    try:
        # Get the temporary directory path
        temp_dir = tempfile.gettempdir()
        # Specify the filename
        file_name = os.path.join(temp_dir, format_to_filename_standard(upstream))
        # Open the file in append mode ('a')
        with open(file_name, "a") as file:
            # Append content to the file
            file.write(f"{platform}_jobid={jobid}\n")
            file.write(f"{platform}_reference={reference}\n")

    except Exception as e:
        pass


# ========== Requests =================


def create_resource():
    response = requests.put(
        f"{INSTALLER_URL}/resources{RESOURCES_TYPE}", headers=headers, json=body
    )
    print("Create Resource Response: " + json.dumps(response.json(), indent=4))
    response.raise_for_status()

    return response.json()


def create_or_get_client(
    client_email: str = "",
    country_name: str = "",
    country_code: str = "",
    client_name: str = "",
):
    response = requests.post(
        f"{INSTALLER_URL}/clients",
        headers=headers,
        json={
            "country_name": country_name.lower(),
            "country_code": country_code.lower(),
            "client_email": client_email.lower(),
            "client_name": client_name.lower(),
        },
    )

    print("Create Create Or Get: " + json.dumps(response.json(), indent=4))

    response.raise_for_status()

    return response.json()


# Resource Type
RESOURCES_TYPE = os.environ.get("RESOURCES_TYPE", "all").strip()
if RESOURCES_TYPE == "all":
    RESOURCES_TYPE = ""
else:
    RESOURCES_TYPE = f"/{RESOURCES_TYPE}"

# PLATFORM
PLATFORM_NAME = os.environ.get("PLATFORM_NAME", "").strip()
PLATFORM_CODE = platforms[PLATFORM_NAME]

# DOMAIN  VARIABLES
DOMAIN_ROOT = os.environ.get("DOMAIN_ROOT", "").strip()
DOMAIN_SUB = os.environ.get(
    "DOMAIN_SUB", platforms_subdomain.get(PLATFORM_NAME, "")
).strip()
DOMAIN_FIELDTYPE = os.environ.get("DOMAIN_FIELDTYPE", "CNAME").strip()
DOMAIN_TTL = os.environ.get("DOMAIN_TTL", "3600").strip()
DOMAIN_TARGET = os.environ.get("DOMAIN_TARGET", "gateway.homelab.net.")
DOMAIN_ENV = os.environ.get("DOMAIN_ENV", "").strip()

if not DOMAIN_SUB or DOMAIN_SUB == "":
    raise Exception("DOMAIN_SUB Cannot be empty")

# ENV
ENV_CODE = environments[DOMAIN_ENV]

# Proxmox Variable
PROXMOX_VM_ID = os.environ.get("PROXMOX_VM_ID", None)
PROXMOX_VM_MEMORY = os.environ.get("PROXMOX_VM_MEMORY", "4096").strip()
PROXMOX_CLOUD_INIT_TEMPLATE = os.environ.get(
    "PROXMOX_CLOUD_INIT_TEMPLATE", cloud_init_templates.get(PLATFORM_NAME, "")
).strip()
PROXMOX_TARGET_NODE = os.environ.get("PROXMOX_TARGET_NODE", "auto").strip()
PROXMOX_VM_CORES = os.environ.get("PROXMOX_VM_CORES", "2").strip()
PROXMOX_NETWORK_BRIDGE = os.environ.get("PROXMOX_NETWORK_BRIDGE", "vmbr3").strip()
PROXMOX_NETWORK_TAG = os.environ.get("PROXMOX_NETWORK_TAG", "10").strip()

PROXMOX_VM_NAME = concatenate_resources(
    PLATFORM_NAME, DOMAIN_SUB, DOMAIN_ENV, DOMAIN_ROOT
)

# Client
PROJECT_EMAIL = os.environ.get("PROJECT_EMAIL", "admin@homelab.com").strip()
PROJECT_COUNTRY = os.environ.get("PROJECT_COUNTRY", "french").strip()
PROJECT_COUNTRY_CODE = os.environ.get("PROJECT_COUNTRY_CODE", "").strip()
PROJECT_NAME = os.environ.get("PROJECT_NAME", "").strip()

# Metadata
METADATA = os.environ.get("METADATA", "")

# UPSTREAM
UPSTREAM_BUILD_TAG = os.environ.get("UPSTREAM_BUILD_TAG")

# ENV TYPE
DOMAIN_SUB_RESULT = DOMAIN_SUB
if DOMAIN_ENV != "prod":
    DOMAIN_SUB_RESULT = concatenate_domain(DOMAIN_SUB, DOMAIN_ENV)


client = create_or_get_client(
    client_email=PROJECT_EMAIL,
    client_name=PROJECT_NAME,
    country_name=PROJECT_COUNTRY,
    country_code=PROJECT_COUNTRY_CODE,
)

if not PROXMOX_VM_ID:
    PROXMOX_VM_ID = f"{client['CountryID']}{client['Code']}{ENV_CODE}{PLATFORM_CODE}"


if not is_memory_binary_format(int(PROXMOX_VM_MEMORY)) or int(PROXMOX_VM_MEMORY) < 2080:
    raise Exception("Wrong ram memory binary format")

if int(PROXMOX_VM_CORES) < 0 or int(PROXMOX_VM_CORES) > 32:
    raise Exception("Invalid VM Core Value")


metadata = convert_dict_values(
    {
        **get_metadata_variables(),
        **extract_metadata_variables(METADATA),
    }
)

add_default_email(metadata, PROJECT_EMAIL)

default_metadata_values_platform(metadata, PLATFORM_NAME)

body = {
    "ref": PROXMOX_VM_NAME,
    "domain": {
        "zone": DOMAIN_ROOT,
        "subdomain": DOMAIN_SUB_RESULT,
        "fieldtype": DOMAIN_FIELDTYPE,
        "ttl": int(DOMAIN_TTL),
        "target": DOMAIN_TARGET,
    },
    "vm": {
        "name": PROXMOX_VM_NAME,
        "target_node": PROXMOX_TARGET_NODE,
        "clone": PROXMOX_CLOUD_INIT_TEMPLATE,
        "vmid": int(PROXMOX_VM_ID),
        "memory": int(PROXMOX_VM_MEMORY),
        "cores": int(PROXMOX_VM_CORES),
        "network": [
            {"bridge": PROXMOX_NETWORK_BRIDGE, "tag": int(PROXMOX_NETWORK_TAG)}
        ],
    },
    "platform": {
        "name": PLATFORM_NAME,
        "metadata": prune_empty(metadata),
    },
}

print("\n")
print("Body: " + json.dumps(body, indent=4))
print("\n")

check_empty_fields(body)


LAST_LOGS = ""


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


JOB_FETCH_INTERVAL = 5  # seconds

if __name__ == "__main__":
    LAST_LOGS = ""
    resource = create_resource()
    jobid = resource["job"]["ID"]

    # Write jobid to upstream file
    write_jobid_to_upstream_file(
        platform=PLATFORM_NAME,
        upstream=UPSTREAM_BUILD_TAG,
        jobid=jobid,
        reference=PROXMOX_VM_NAME,
    )

    print(f"JobID: {jobid} ... ")
    # Wait for job status
    while True:
        if wait_job_status(jobid=jobid):
            break
        time.sleep(JOB_FETCH_INTERVAL)
