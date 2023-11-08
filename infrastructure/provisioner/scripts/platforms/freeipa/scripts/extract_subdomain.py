import sys


def extract_subdomain(full_domain):
    parts = full_domain.split(".")
    if len(parts) > 2:
        return ".".join(parts[-(len(parts) - 1) :])
    return full_domain


if __name__ == "__main__":
    print(extract_subdomain(sys.argv[1].strip()), end="")
