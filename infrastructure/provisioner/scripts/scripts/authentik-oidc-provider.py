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
REDIRECT_URL = os.getenv("REDIRECT_URL")
LAUNCH_URL = os.getenv("LAUNCH_URL", "")

for value in (SESSION_HEADERS, AUTHENTIK_URL, PROVIDER_NAME, REDIRECT_URL):
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


def get_property_mappings():
    response = requests.get(
        f"{AUTHENTIK_URL}/propertymappings/scope/?ordering=scope_name", headers=headers
    )

    response.raise_for_status()
    json = response.json()
    return json["results"]


def get_authentication_flows_instances():
    response = requests.get(
        f"{AUTHENTIK_URL}/flows/instances/?designation=authentication&ordering=slug",
        headers=headers,
    )

    response.raise_for_status()
    json = response.json()
    return json["results"]


def get_authorization_flows_instances():
    response = requests.get(
        f"{AUTHENTIK_URL}/flows/instances/?designation=authorization&ordering=slug",
        headers=headers,
    )

    response.raise_for_status()
    json = response.json()
    return json["results"]


def get_certificate_keypairs():
    response = requests.get(
        f"{AUTHENTIK_URL}/crypto/certificatekeypairs/?has_key=true&include_details=false&ordering=name",
        headers=headers,
    )

    response.raise_for_status()
    json = response.json()
    return json["results"]


def get_provider_oauth2(name: str):
    response = requests.get(
        f"{AUTHENTIK_URL}/providers/all/?search={name}", headers=headers
    )

    response.raise_for_status()
    json = response.json()
    return json["results"]


def get_application(name: str):
    name = slugify(name)
    response = requests.get(
        f"{AUTHENTIK_URL}/core/applications/?search={name}&superuser_full_list=true",
        headers=headers,
    )

    response.raise_for_status()
    json = response.json()
    return json["results"]


def get_provider_oauth2_pk(pk: str):
    response = requests.get(f"{AUTHENTIK_URL}/providers/oauth2/{pk}/", headers=headers)

    response.raise_for_status()
    return response.json()


def get_provider_oauth2_setup_urls_pk(pk: str):
    response = requests.get(
        f"{AUTHENTIK_URL}/providers/oauth2/{pk}/setup_urls/", headers=headers
    )

    response.raise_for_status()
    return response.json()


def log(provider_pk):
    provider = get_provider_oauth2_pk(provider_pk)
    json_string = json.dumps(
        {
            **get_provider_oauth2_setup_urls_pk(provider_pk),
            "client_id": provider["client_id"],
            "client_secret": provider["client_secret"],
        }
    )
    base64_encoded = base64.b64encode(json_string.encode("utf-8")).decode("utf-8")

    print(f"%%%{base64_encoded}%%%")


def main():
    provider_pk = get_pk(get_provider_oauth2(PROVIDER_NAME), "name", PROVIDER_NAME)

    if provider_pk:
        log(provider_pk)
        return

    response = requests.post(
        f"{AUTHENTIK_URL}/providers/oauth2/",
        headers=headers,
        json={
            "name": PROVIDER_NAME,
            "authentication_flow": get_pk(
                get_authentication_flows_instances(),
                "slug",
                "default-authentication-flow",
            ),
            "authorization_flow": get_pk(
                get_authorization_flows_instances(),
                "slug",
                "default-provider-authorization-explicit-consent",
            ),
            "signing_key": get_pk(get_certificate_keypairs())[0],
            "property_mappings": get_pk(get_property_mappings()),
            "client_type": "confidential",
            "client_id": generate_random_key(40),
            "client_secret": generate_random_key(128),
            "access_code_validity": "minutes=1",
            "access_token_validity": "minutes=5",
            "refresh_token_validity": "days=30",
            "include_claims_in_id_token": True,
            "redirect_uris": REDIRECT_URL,
            "sub_mode": "hashed_user_id",
            "issuer_mode": "per_provider",
            "jwks_sources": [],
        },
    )

    response.raise_for_status()
    provider = response.json()

    # Create Application
    requests.post(
        f"{AUTHENTIK_URL}/core/applications/",
        headers=headers,
        json={
            "name": PROVIDER_NAME,
            "slug": slugify(PROVIDER_NAME),
            "provider": provider["pk"],
            "backchannel_providers": [],
            "open_in_new_tab": False,
            "meta_launch_url": LAUNCH_URL,
            "meta_description": "",
            "meta_publisher": "",
            "policy_engine_mode": "any",
            "group": "",
        },
    ).raise_for_status()

    log(provider["pk"])


if __name__ == "__main__":
    decoded_bytes = base64.b64decode(SESSION_HEADERS)
    headers = json.loads(decoded_bytes.decode("utf-8"))

    main()
