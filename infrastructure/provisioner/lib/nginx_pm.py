import re
import json
import base64
import argparse
import requests
from typing import Any, List
from utilities.dotenv import config
from utilities.logging import logging, bingLoggingConfig


headers = {"Content-Type": "application/json", "Accept": "application/json"}

DOMAIN_KEY = "domain"
SSL_PROVIDER = "letsencrypt"


def clean_domain(url: str):
    if url.startswith("http://"):
        url = url[len("http://") :]
    elif url.startswith("https://"):
        url = url[len("https://") :]

    if url.endswith("/"):
        url = url[:-1]

    return url


def remove_none_values(input_list):
    # Filter out all None values from the input_list
    return [item for item in input_list if item is not None]


def filter_domains(domains):
    # Regular expression pattern for matching valid domain names
    pattern = r"^(?=.{1,255}$)([A-Za-z0-9](?:(?:[A-Za-z0-9\-]){0,61}[A-Za-z0-9])?\.)+[A-Za-z]{2,}$"
    # Compile the pattern for better performance
    regex = re.compile(pattern)

    # Use list comprehension to filter the domains
    matched_domains = [domain for domain in domains if regex.match(domain)]

    return matched_domains


def get_api_token():
    url = config.get("NGINX_PM_URL")
    email = config.get("NGINX_PM_EMAIL")
    password = config.get("NGINX_PM_PASSWORD")

    if not url or not email or not password:
        logging.warning(
            "Cannot read variable env (NGINX_PM_URL, NGINX_PM_EMAIL, NGINX_PM_PASSWORD)"
        )
        return None, None

    res = requests.post(
        f"{url}/api/tokens",
        json={"identity": email, "secret": password},
        headers=headers,
    )
    if res.status_code != 200:
        logging.warning("Nginx Proxy Manager authentication failed")
        return None, None

    json = res.json()

    token = json.get("token")

    headers["Authorization"] = f"Bearer {token}"

    return url, token


def get_decoded_domain(metadata: str):
    decoded_bytes = base64.b64decode(metadata)
    data: dict[str, str] = json.loads(decoded_bytes.decode("utf-8"))
    domain_main = data.get(DOMAIN_KEY, None)
    # If _mapping was passed then ignore the current _mapping process
    ignore = (
        data.get("_mapping", None) == False or data.get("_mapping", None) == "false"
    )

    if ignore:
        return None

    domains = [clean_domain(domain_main) if domain_main else None]

    for key in data.keys():
        SUBDOMAIN_FILTER = "_subdomain"
        if key.endswith(SUBDOMAIN_FILTER) and data.get(key):
            domains.append(
                clean_domain(data.get(key)),
            )

    domains = remove_none_values(domains)
    domains = filter_domains(domains)

    return domains if len(domains) > 0 else None


def delete_proxy_hosts(pHost: Any, url: str):
    if pHost:
        requests.delete(
            f"{url}/api/nginx/proxy-hosts/{pHost.get('id')}", headers=headers
        )


def find_existing_proxy_host(domains: list[str], url: str):
    fDomain = domains[0]

    res = requests.get(f"{url}/api/nginx/proxy-hosts?query={fDomain}", headers=headers)
    if res.status_code != 200:
        return None
    # Check for the API Schema
    # (https://github.com/NginxProxyManager/nginx-proxy-manager/blob/develop/backend/schema/endpoints/proxy-hosts.json)

    hosts: List[Any] = res.json()

    for host in hosts:
        validHost = True

        domain_names: List[str] = host.get("domain_names")

        if len(domain_names) != len(domains):
            validHost = False
            continue

        for domain in domains:
            try:
                domain_names.index(domain)
            except Exception as err:
                validHost = False
                logging.error(err)
                continue

        if validHost:
            return host

    return None


def find_existing_certificate(domains: list[str], url: str):
    fDomain = domains[0]

    res = requests.get(f"{url}/api/nginx/certificates?query={fDomain}", headers=headers)
    if res.status_code != 200:
        return None
    # Check for the API Schema
    # (https://github.com/NginxProxyManager/nginx-proxy-manager/blob/develop/backend/schema/endpoints/certificates.json)

    certificates: List[Any] = res.json()

    for certificate in certificates:
        validCert = True

        domain_names: List[str] = certificate.get("domain_names")

        for domain in domains:
            try:
                domain_names.index(domain)
            except Exception as err:
                validCert = False
                logging.error(err)
                continue

        if validCert:
            return certificate

    return None


def get_platform_protocol(platform: str):
    data: dict[str, Any] = {}
    with open("scripts/platform-protocols.json", "r") as file:
        data = json.loads(file.read())

    return data.get(platform)


def get_platform_nginx_pm_config(platform: str):
    data: dict[str, Any] = {}
    with open("scripts/platform-nginx-pm.json", "r") as file:
        data = json.loads(file.read())

    return data.get(platform)


def create_domains_certificate(domains: List[str], url: str):
    certificate = find_existing_certificate(domains, url)
    if certificate:
        return certificate

    body = {
        "domain_names": domains,
        "meta": {
            "letsencrypt_email": config.get("ADMIN_SYSTEM_EMAIL", "admin@example.com"),
            "letsencrypt_agree": True,
            "dns_challenge": False,
        },
        "provider": SSL_PROVIDER,
    }
    # Check for the API Schema
    # (https://github.com/NginxProxyManager/nginx-proxy-manager/blob/develop/backend/schema/endpoints/certificates.json)
    res = requests.post(f"{url}/api/nginx/certificates", json=body, headers=headers)

    if res.status_code >= 200 and res.status_code < 400:
        return res.json()

    return None


def create_proxy_host(
    url: str,
    domains: List[str],
    certificate: Any,
    platform_protocol: Any,
    ip: str,
    platform: str,
):
    certificate = certificate if certificate else {}
    certificate_id = certificate.get("id")

    advanced_config = get_platform_nginx_pm_config(platform)
    advanced_config = ("" if not advanced_config else advanced_config).replace(
        "\\n", "\n"
    )

    body = {
        "domain_names": domains,
        "forward_host": ip,
        "forward_scheme": platform_protocol.get("protocol"),
        "forward_port": platform_protocol.get("port"),
        "block_exploits": True,
        "allow_websocket_upgrade": True,
        "access_list_id": "0",
        "certificate_id": certificate_id if certificate_id else 0,
        "ssl_forced": True if certificate_id else False,
        "http2_support": True if certificate_id else False,
        "hsts_enabled": True if certificate_id else False,
        "meta": {"letsencrypt_agree": False, "dns_challenge": False},
        "advanced_config": advanced_config,
        "locations": [],
        "caching_enabled": False,
        "hsts_subdomains": False,
    }

    requests.post(f"{url}/api/nginx/proxy-hosts", json=body, headers=headers)


def main(action: str, metadata: str, platform: str, ip: str):
    # Decode metadata and get the domain value
    domains = get_decoded_domain(metadata)
    if not domains:
        return

    url, token = get_api_token()
    if not url or not token:
        return

    pHost = find_existing_proxy_host(domains, url)

    # Check of delete action
    if action == "delete":
        logging.info("Deleting... proxy host")
        delete_proxy_hosts(pHost, url)
        return

    # If pHost exists and delete it
    if pHost:
        delete_proxy_hosts(pHost, url)

    # Get platform proxy
    platform_protocol = get_platform_protocol(platform)

    if not platform_protocol:
        logging.info("Cannot found the corresponding platform protocol")
        return

    # Generate domain certificate
    certificate = create_domains_certificate(domains, url)

    # Finally create the proxy host
    create_proxy_host(
        url=url,
        domains=domains,
        certificate=certificate,
        platform_protocol=platform_protocol,
        ip=ip,
        platform=platform,
    )


if __name__ == "__main__":
    bingLoggingConfig(prefix="NGINX PM / ")

    parser = argparse.ArgumentParser()
    parser.add_argument("--action", choices=["create", "delete"], required=True)
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
