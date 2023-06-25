import json
import base64
import argparse
import requests
from typing import Any, List
from utilities.dotenv import config
from utilities.logging import logging, bingLoggingConfig


headers = {
    "Content-Type": "application/json",
    "Accept": "application/json"
}

DOMAIN_KEY = "domain"
SSL_PROVIDER = "letsencrypt"


def clean_domain(url: str):
    if url.startswith("http://"):
        url = url[len("http://"):]
    elif url.startswith("https://"):
        url = url[len("https://"):]

    if url.endswith("/"):
        url = url[:-1]

    return url


def get_api_token() -> (tuple[None, None] | tuple[str, str]):
    url = config.get("NGINX_PM_URL")
    email = config.get("NGINX_PM_EMAIL")
    password = config.get("NGINX_PM_PASSWORD")

    if not url or not email or not password:
        logging.warn(
            "Cannot read variable env (NGINX_PM_URL, NGINX_PM_EMAIL, NGINX_PM_PASSWORD)"
        )
        return None, None

    res = requests.post(
        f"{url}/api/tokens",
        json={"identity": email, "secret": password},
        headers=headers
    )
    if res.status_code != 200:
        logging.warn(
            "Nginx proxy mananger authentication failed"
        )
        return None, None

    json = res.json()

    token = json.get('token')

    headers["Authorization"] = f"Bearer {token}"

    return url, token


def get_decoded_domain(metadata: str) -> str | None:
    decoded_bytes = base64.b64decode(metadata)
    data = json.loads(decoded_bytes.decode("utf-8"))
    domain = data.get(DOMAIN_KEY, None)

    return clean_domain(domain) if domain else None


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
        domains: List[str] = phost.get('domain_names')
        try:
            domains.index(domain)
            return phost
        except Exception as err:
            logging.error(err)
            continue

    return None


def find_existing_certificate(domain: str, url: str):
    res = requests.get(
        f"{url}/api/nginx/certificates",
        headers=headers
    )
    if res.status_code != 200:
        return None
    # Check for the API Schema
    # (https://github.com/NginxProxyManager/nginx-proxy-manager/blob/develop/backend/schema/endpoints/certificates.json)

    data: List[Any] = res.json()

    for phost in data:
        domains: List[str] = phost['domain_names']
        try:
            domains.index(domain)
            return phost
        except:
            continue

    return None


def get_platform_schema(platform: str):
    data: dict[str, Any] = {}
    with open("scripts/platform-schemas.json", "r") as file:
        data = json.loads(file.read())

    return data.get(platform)


def create_domain_certificate(domain: str, url: str):
    certificate = find_existing_certificate(domain, url)
    if certificate:
        return certificate

    body = {
        "domain_names": [domain],
        "meta": {
            "letsencrypt_email": config.get("ADMIN_SYSTEM_EMAIL", "admin@example.com"),
            "letsencrypt_agree": True,
            "dns_challenge": False
        },
        "provider": SSL_PROVIDER
    }
    # Check for the API Schema
    # (https://github.com/NginxProxyManager/nginx-proxy-manager/blob/develop/backend/schema/endpoints/certificates.json)
    res = requests.post(
        f"{url}/api/nginx/certificates",
        json=body,
        headers=headers
    )

    if res.status_code >= 200 and res.status_code < 400:
        return res.json()

    return None


def create_proxy_host(url: str, domain: str, certificate: Any, platform_schema: Any, ip: str):
    certificate = certificate if certificate else {}
    certificate_id = certificate.get('id')

    body = {
        "domain_names": [domain],
        "forward_host": ip,
        "forward_scheme": platform_schema.get('protocol'),
        "forward_port": platform_schema.get('port'),
        "block_exploits": True,
        "allow_websocket_upgrade": True,
        "access_list_id": "0",
        "certificate_id": certificate_id if certificate_id else 0,
        "ssl_forced": True if certificate_id else False,
        "http2_support":  True if certificate_id else False,
        "hsts_enabled":  True if certificate_id else False,
        "meta": {
            "letsencrypt_agree": False,
            "dns_challenge": False
        },
        "advanced_config": "",
        "locations": [],
        "caching_enabled": False,
        "hsts_subdomains": False
    }

    requests.post(
        f"{url}/api/nginx/proxy-hosts",
        json=body,
        headers=headers
    )


def main(action: str, metadata: str, platform: str, ip: str):
    url, token = get_api_token()
    if not url or not token:
        return

    # Decode metade and get the domain value
    domain = get_decoded_domain(metadata)
    phost = find_existing_proxy_host(domain, url)

    # Check of delete action
    if action == "delete":
        logging.info(
            "Deleting... proxy host"
        )
        delete_proxy_hosts(phost, url)
        return

    # If phost exists and has create action then no need go futher
    if phost:
        logging.info(
            "Exited proxy host exist while has create action"
        )
        return

    # Get platform proxy
    platform_schema = get_platform_schema(platform)

    if not platform_schema:
        logging.info(
            "Cannot found the corresponding platform schema"
        )
        return

    # Generate domain certificate
    certificate = create_domain_certificate(domain, url)

    # Finally create the proxy host
    create_proxy_host(
        url=url,
        domain=domain,
        certificate=certificate,
        platform_schema=platform_schema,
        ip=ip
    )


if __name__ == '__main__':
    bingLoggingConfig(prefix="NGINX PM / ")

    parser = argparse.ArgumentParser()
    parser.add_argument(
        "--action",
        choices=['create', 'delete'],
        required=True
    )
    parser.add_argument("--metadata", required=True)
    parser.add_argument("--platform", required=False)
    parser.add_argument("--ip", required=False)
    args = parser.parse_args()

    logging.info(args)

    # Process nginx pm
    main(
        action=args.action,
        metadata=args.metadata,
        platform=args.platform,
        ip=args.ip,
    )
