#!/usr/bin/env python3

import base64
import time
import requests
import os
import json


INSTALLER_URL = os.environ.get("INSTALLER_URL", "")
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


def check_empty_fields(data, path=[]):
    for key, value in data.items():
        current_path = path + [key]

        if isinstance(value, dict):
            check_empty_fields(value, current_path)
        elif value is None or (isinstance(value, str) and value.strip() == ""):
            raise ValueError(f"Empty field at path: {'.'.join(current_path)}")


def join_and_replace(*args):
    # Replace dots and spaces with hyphens in each argument
    replaced_args = [arg.replace(".", "-").replace(" ", "-") for arg in args]
    # Join the modified arguments with hyphens
    result = "-".join(replaced_args)
    return result


# ========== Requests =================


def create_resource():
    response = requests.put(
        f"{INSTALLER_URL}/resources/domain",
        headers=headers,
        json=body,
    )
    print("Create Resource Response: " + json.dumps(response.json(), indent=4))
    response.raise_for_status()

    return response.json()


# DOMAIN  VARIABLES
DOMAIN_ROOT = os.environ.get("DOMAIN_ROOT", "").strip()
DOMAIN_SUB = os.environ.get("DOMAIN_SUB", "").strip()
DOMAIN_FIELDTYPE = os.environ.get("DOMAIN_FIELDTYPE", "").strip()
DOMAIN_TTL = os.environ.get("DOMAIN_TTL", "3600").strip()
DOMAIN_TARGET = os.environ.get("DOMAIN_TARGET", "")

if not DOMAIN_SUB or DOMAIN_SUB == "":
    raise Exception("DOMAIN_SUB Cannot be empty")


body = {
    "ref": join_and_replace(DOMAIN_SUB, DOMAIN_FIELDTYPE, DOMAIN_ROOT),
    "domain": {
        "zone": DOMAIN_ROOT,
        "subdomain": DOMAIN_SUB,
        "fieldtype": DOMAIN_FIELDTYPE,
        "ttl": int(DOMAIN_TTL),
        "target": DOMAIN_TARGET,
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

    print(f"JobID: {jobid} ... ")
    # Wait for job status
    while True:
        if wait_job_status(jobid=jobid):
            break
        time.sleep(JOB_FETCH_INTERVAL)
