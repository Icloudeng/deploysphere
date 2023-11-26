import os
import json
import base64
import requests

SESSION_HEADERS = os.getenv("SESSION_HEADERS")
AUTHENTIK_URL = os.getenv("AUTHENTIK_URL")
AUTHENTIK_ADMIN_EMAIL = os.getenv("AUTHENTIK_ADMIN_EMAIL")
AUTHENTIK_ADMIN_PASSWORD = os.getenv("AUTHENTIK_ADMIN_PASSWORD")


for value in (
    SESSION_HEADERS,
    AUTHENTIK_URL,
    AUTHENTIK_ADMIN_EMAIL,
    AUTHENTIK_ADMIN_PASSWORD,
):
    if not value or len(value) == 0:
        raise Exception("Envs varaible required !")

AUTHENTIK_URL = f"{AUTHENTIK_URL}/api/v3"


def main(headers: dict):
    body = {
        "email": AUTHENTIK_ADMIN_EMAIL,
        "password": AUTHENTIK_ADMIN_PASSWORD,
        "password_repeat": AUTHENTIK_ADMIN_PASSWORD,
        "component": "ak-stage-prompt",
    }

    del headers["content-type"]

    requests.post(
        f"{AUTHENTIK_URL}/flows/executor/initial-setup/?query=",
        headers=headers,
        json=body,
    ).raise_for_status()


if __name__ == "__main__":
    decoded_bytes = base64.b64decode(SESSION_HEADERS)
    main(json.loads(decoded_bytes.decode("utf-8")))
