#!/usr/bin/env python3

import re
import base64
import time
import requests
import os
import json
import tempfile


LAST_LOGS = ""
JOB_FETCH_INTERVAL = 5  # seconds

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
    "zimbra": "50",
}


platforms_subdomain = {
    "zimbra": "mail",
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


def is_memory_binary_format(number):
    # Check if the number is a power of 2
    return number & (number - 1) == 0 and number != 0


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


def first_split_domain(subdomain: str):
    arr = subdomain.split(".")[1:]
    return ".".join(arr)


# ========== Requests =================


def create_resource(type: str = "", body={}):
    response = requests.put(
        f"{INSTALLER_URL}/resources{type}", headers=headers, json=body
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
    "PROXMOX_CLOUD_INIT_TEMPLATE", "ubuntu-20.04-cloudinit-template"
).strip()
PROXMOX_TARGET_NODE = os.environ.get("PROXMOX_TARGET_NODE", "auto").strip()
PROXMOX_VM_CORES = os.environ.get("PROXMOX_VM_CORES", "2").strip()
PROXMOX_NETWORK_BRIDGE = os.environ.get("PROXMOX_NETWORK_BRIDGE", "vmbr3").strip()
PROXMOX_NETWORK_TAG = os.environ.get("PROXMOX_NETWORK_TAG", "11").strip()

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


# Resources
fqdn_domain = f"{DOMAIN_SUB_RESULT}.{DOMAIN_ROOT}"

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
            {"bridge": PROXMOX_NETWORK_BRIDGE},
        ],
    },
    "platform": {
        "name": PLATFORM_NAME,
        "metadata": {
            "zimbra_fqdn": fqdn_domain,
            "zimbra_timezone": "Africa/Central",
            "mapping": False,
        },
    },
}

mx_body = {
    "ref": f"mx-{PROXMOX_VM_NAME}",
    "domain": {
        "zone": DOMAIN_ROOT,
        "subdomain": first_split_domain(DOMAIN_SUB_RESULT),
        "fieldtype": "MX",
        "ttl": 3600,
        "target": f"{fqdn_domain}.",
    },
}

print("\n")
print("Body: " + json.dumps(body, indent=4))
print("\n")

check_empty_fields(body)


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


def wait_for_job(jobid):
    global LAST_LOGS
    LAST_LOGS = ""

    while True:
        if wait_job_status(jobid=jobid):
            break
        time.sleep(JOB_FETCH_INTERVAL)


if __name__ == "__main__":
    # Create MX Domain Resource
    resource = create_resource(body=mx_body, type="domain")
    jobid = resource["job"]["ID"]
    print(f"Create Domain resource JobID: {jobid} ... ")
    wait_for_job(jobid)

    # Create all resources
    resource = create_resource(body=body)
    jobid = resource["job"]["ID"]

    # Write jobid to upstream file
    write_jobid_to_upstream_file(
        platform=PLATFORM_NAME,
        upstream=UPSTREAM_BUILD_TAG,
        jobid=jobid,
        reference=PROXMOX_VM_NAME,
    )

    print(f"Create resources JobID: {jobid} ... ")
    wait_for_job(jobid)
