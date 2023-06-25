import json
import base64
import argparse
import requests
from typing import Any, List
from .utilities.dotenv import config


headers = {
    "Content-Type": "application/json",
    "Accept": "application/json"
}

DOMAIN_KEY = "domain"


def get_api_token() -> (tuple[None, None] | tuple[str, str]):
    url = config.get("NGINX_PM_URL")
    email = config.get("NGINX_PM_EMAIL")
    password = config.get("NGINX_PM_PASSWORD")

    if not url or not email or password:
        return None, None

    res = requests.post(
        f"{url}/api/tokens",
        data={"identity": email, "secret": password},
        headers=headers
    )
    if res.status_code != 200:
        return None, None

    json = res.json()

    token = json.get('token')

    headers["Authorization"] = f"Bearer {token}"

    return url, token


def get_decoded_domain(metadata: str) -> str | None:
    decoded_bytes = base64.b64decode(metadata)
    data = json.loads(decoded_bytes.decode("utf-8"))
    return data.get(DOMAIN_KEY, None)


def delete_proxy_hosts(phost: Any, url: str):
    if phost:
        requests.delete(
            f"{url}/api/nginx/proxy-hosts/{phost.id}",
            headers=headers
        )


def find_existing_proxy_host(domain: str, url: str):
    res = requests.get(
        f"{url}/api/nginx/proxy-hosts",
        headers=headers
    )
    if res.status_code != 200:
        return None
    # Check for the API Schema
    # (https://github.com/NginxProxyManager/nginx-proxy-manager/blob/develop/backend/schema/endpoints/proxy-hosts.json)

    data: List[Any] = res.json()

    for phost in data:
        domains: List[str] = phost.domain_names
        try:
            domains.index(domain)
            return phost
        except:
            continue

    return None


def get_platform_port(platform: str):
    data: dict[str, int] = {}
    with open("scripts/platform-ports.json", "r") as file:
        data = json.loads(file.read())

    return data.get(platform)


def main(action: str, metadata: str, platform: str, status: str, ip: str):
    url, token = get_api_token()
    if not url or token:
        return
    # Decode metade and get the domain value
    domain = get_decoded_domain(metadata)
    phost = find_existing_proxy_host(domain, url)

    # Check of delete action
    if action == "delete":
        delete_proxy_hosts(phost, url)
        return

    # If phost exists and has create action then no need go futher
    if phost:
        return

    # Get platform proxy
    platform_port = get_platform_port(platform)

    if not platform_port:
        return


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--action",
        choices=['create', 'delete'],
        required=True
    )
    parser.add_argument("--metadata", required=True)
    parser.add_argument("--platform", required=False)
    parser.add_argument("--status", required=False)
    parser.add_argument("--ip", required=False)
    args = parser.parse_args()

    # Process nginx pm
    main(
        action=args.action,
        metadata=args.metadata,
        platform=args.platform,
        status=args.status,
        ip=args.ip,
    )
