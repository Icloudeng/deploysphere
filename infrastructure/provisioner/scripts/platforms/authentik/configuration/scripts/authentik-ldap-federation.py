import os
import json
import base64
import requests
import secrets
import string
import re

SESSION_HEADERS = os.getenv("SESSION_HEADERS")
AUTHENTIK_URL = os.getenv("AUTHENTIK_URL")
PROVIDER_NAME = os.getenv("PROVIDER_NAME")

# LDAP
LDAP_SERVER_URL = os.getenv("LDAP_SERVER_URL")
LDAP_BASE_CN = os.getenv("LDAP_BASE_CN")
LDAP_BIND_DN = os.getenv("LDAP_BIND_DN")
LDAP_BIND_PASSWORD = os.getenv("LDAP_BIND_PASSWORD")


for value in (
    SESSION_HEADERS,
    AUTHENTIK_URL,
    PROVIDER_NAME,
    LDAP_SERVER_URL,
    LDAP_BASE_CN,
    LDAP_BIND_DN,
    LDAP_BIND_PASSWORD,
):
    if not value or len(value) == 0:
        raise Exception("Envs varaible required !")

AUTHENTIK_URL = f"{AUTHENTIK_URL}/api/v3"

headers = {}


def slugify(text):
    text = re.sub(r"[^\w\s-]", "", text).strip().lower()
    text = re.sub(r"[-\s]+", "-", text)
    return text


def generate_random_key(length: int):
    characters = string.ascii_letters + string.digits
    random_key = "".join(secrets.choice(characters) for _ in range(length))
    return random_key


def get_pk(results=[], field=None, value=None):
    pks = []
    for v in results:
        pks.append(v.get("pk"))

        if field and value and v.get(field) == value:
            return v.get("pk")

    if field and value:
        return None

    return pks


def get_ldap_property_mappings():
    response = requests.get(
        f"{AUTHENTIK_URL}/propertymappings/ldap/?ordering=managed%2Cobject_field",
        headers=headers,
    )

    response.raise_for_status()
    json = response.json()
    return json["results"]


def get_source_ldap(name: str):
    response = requests.get(
        f"{AUTHENTIK_URL}/sources/all/?search={name}", headers=headers
    )

    response.raise_for_status()
    json = response.json()
    return json["results"]


def main():
    provider_name_slugified = slugify(PROVIDER_NAME)
    provider_pk = get_pk(
        get_source_ldap(provider_name_slugified), "name", PROVIDER_NAME
    )

    mappings = get_ldap_property_mappings()
    body = {
        "name": PROVIDER_NAME,
        "slug": provider_name_slugified,
        "enabled": True,
        "user_path_template": "goauthentik.io/sources/%(slug)s",
        "server_uri": LDAP_SERVER_URL,
        "peer_certificate": "",
        "client_certificate": "",
        "bind_cn": LDAP_BIND_DN,
        "bind_password": LDAP_BIND_PASSWORD,
        "start_tls": False,
        "sni": False,
        "base_dn": LDAP_BASE_CN,
        "additional_user_dn": "cn=users",
        "additional_group_dn": "cn=groups",
        "user_object_filter": "(objectClass=person)",
        "group_object_filter": "(objectClass=nestedGroup)",
        "group_membership_field": "member",
        "object_uniqueness_field": "ipaUniqueID",
        "sync_users": True,
        "sync_users_password": True,
        "sync_groups": True,
        "sync_parent_group": "",
        "property_mappings": get_pk(mappings),
        "property_mappings_group": [
            get_pk(mappings, "managed", "goauthentik.io/sources/ldap/default-name"),
            get_pk(mappings, "managed", "goauthentik.io/sources/ldap/openldap-cn"),
        ],
    }

    if provider_pk:
        requests.patch(
            f"{AUTHENTIK_URL}/sources/ldap/{provider_name_slugified}/",
            headers=headers,
            json=body,
        ).raise_for_status()
    else:
        requests.post(
            f"{AUTHENTIK_URL}/sources/ldap/",
            headers=headers,
            json=body,
        ).raise_for_status()
        # Send another request to synchonize LDAP entries
        requests.patch(
            f"{AUTHENTIK_URL}/sources/ldap/{provider_name_slugified}/",
            headers=headers,
            json=body,
        ).raise_for_status()


if __name__ == "__main__":
    decoded_bytes = base64.b64decode(SESSION_HEADERS)
    headers = json.loads(decoded_bytes.decode("utf-8"))

    main()
