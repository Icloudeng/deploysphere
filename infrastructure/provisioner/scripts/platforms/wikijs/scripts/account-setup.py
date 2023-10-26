import os
import requests


BASE_URL = os.getenv("BASE_URL")
SITE_URL = os.getenv("SITE_URL")
ADMIN_EMAIL = os.getenv("ADMIN_EMAIL")
ADMIN_PASSWORD = os.getenv("ADMIN_PASSWORD")


for value in (SITE_URL, ADMIN_PASSWORD, ADMIN_EMAIL):
    if not value or len(value) == 0:
        raise Exception("Envs varaible required !")


headers = {
    "Content-Type": "application/json",
    "Accept": "*/*",
    "User-Agent": "PostmanRuntime/7.34.0",
}


def main():
    requests.post(
        f"{BASE_URL}/finalize",
        headers=headers,
        json={
            "adminEmail": ADMIN_EMAIL,
            "adminPassword": ADMIN_PASSWORD,
            "adminPasswordConfirm": ADMIN_PASSWORD,
            "siteUrl": SITE_URL,
            "telemetry": True,
        },
    )


if __name__ == "__main__":
    main()
