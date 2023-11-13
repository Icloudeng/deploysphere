import os
import requests


INSTANCE_URL = os.getenv("INSTANCE_URL")

# LOGIN
LOGIN_USERNAME = os.getenv("LOGIN_USERNAME")
LOGIN_PASSWORD = os.getenv("LOGIN_PASSWORD")


# OIDC
OIDC_KEY = os.getenv("OIDC_KEY")
OIDC_CLIENT_ID = os.getenv("OIDC_CLIENT_ID")
OIDC_CLIENT_SECRET = os.getenv("OIDC_CLIENT_SECRET")
OIDC_AUTHORIZATION_URL = os.getenv("OIDC_AUTHORIZATION_URL")
OIDC_TOKEN_URL = os.getenv("OIDC_TOKEN_URL")
OIDC_USER_INFO_URL = os.getenv("OIDC_USER_INFO_URL")
OIDC_ISSUER = os.getenv("OIDC_ISSUER")
OIDC_LOGOUT_URL = os.getenv("OIDC_LOGOUT_URL")


for value in (
    INSTANCE_URL,
    LOGIN_USERNAME,
    LOGIN_PASSWORD,
    OIDC_KEY,
    OIDC_CLIENT_ID,
    OIDC_CLIENT_SECRET,
    OIDC_AUTHORIZATION_URL,
    OIDC_TOKEN_URL,
    OIDC_USER_INFO_URL,
    OIDC_ISSUER,
    OIDC_LOGOUT_URL,
):
    if not value or len(value) == 0:
        raise Exception("Envs varaible required !")


url = f"{INSTANCE_URL}/graphql"
headers = {
    "Content-Type": "application/json",
}


def login():
    body = [
        {
            "operationName": None,
            "variables": {
                "username": LOGIN_USERNAME,
                "password": LOGIN_PASSWORD,
                "strategy": "local",
            },
            "extensions": {},
            "query": "mutation ($username: String!, $password: String!, $strategy: String!) {\n  authentication {\n    login(username: $username, password: $password, strategy: $strategy) {\n      responseResult {\n        succeeded\n        errorCode\n        slug\n        message\n        __typename\n      }\n      jwt\n      mustChangePwd\n      mustProvideTFA\n      mustSetupTFA\n      continuationToken\n      redirect\n      tfaQRImage\n      __typename\n    }\n    __typename\n  }\n}\n",
        }
    ]
    response = requests.post(url, json=body)
    response.raise_for_status()

    data = response.json()[0]
    login = data["data"]["authentication"]["login"]
    responseResult = login["responseResult"]

    if not responseResult["succeeded"]:
        raise Exception(responseResult["message"])

    jwt_token = login["jwt"]
    headers["Authorization"] = f"Bearer {jwt_token}"


def create_oidc():
    body = [
        {
            "operationName": None,
            "variables": {
                "strategies": [
                    {
                        "key": "local",
                        "strategyKey": "local",
                        "displayName": "Local",
                        "order": 0,
                        "isEnabled": True,
                        "config": [],
                        "selfRegistration": False,
                        "domainWhitelist": [],
                        "autoEnrollGroups": [],
                    },
                    {
                        "key": OIDC_KEY,
                        "strategyKey": "oidc",
                        "displayName": "Generic OpenID Connect / OAuth2",
                        "order": 1,
                        "isEnabled": True,
                        "config": [
                            {
                                "key": "clientId",
                                "value": '{"v":"%s"}' % (OIDC_CLIENT_ID),
                            },
                            {
                                "key": "clientSecret",
                                "value": '{"v":"%s"}' % (OIDC_CLIENT_SECRET),
                            },
                            {
                                "key": "authorizationURL" % (OIDC_AUTHORIZATION_URL),
                                "value": '{"v":"%s"}',
                            },
                            {
                                "key": "tokenURL",
                                "value": '{"v":"%s"}' % (OIDC_TOKEN_URL),
                            },
                            {
                                "key": "userInfoURL",
                                "value": '{"v":"%s"}' % (OIDC_USER_INFO_URL),
                            },
                            {"key": "skipUserProfile", "value": '{"v":false}'},
                            {
                                "key": "issuer",
                                "value": '{"v":"%s"}' % (OIDC_ISSUER),
                            },
                            {"key": "emailClaim", "value": '{"v":"email"}'},
                            {"key": "displayNameClaim", "value": '{"v":"displayName"}'},
                            {"key": "mapGroups", "value": '{"v":false}'},
                            {"key": "groupsClaim", "value": '{"v":"groups"}'},
                            {
                                "key": "logoutURL",
                                "value": '{"v":"%s"}' % (OIDC_LOGOUT_URL),
                            },
                            {"key": "acrValues", "value": '{"v":""}'},
                        ],
                        "selfRegistration": True,
                        "domainWhitelist": [],
                        "autoEnrollGroups": [2],
                    },
                ]
            },
            "extensions": {},
            "query": "mutation ($strategies: [AuthenticationStrategyInput]!) {\n  authentication {\n    updateStrategies(strategies: $strategies) {\n      responseResult {\n        succeeded\n        errorCode\n        slug\n        message\n        __typename\n      }\n      __typename\n    }\n    __typename\n  }\n}\n",
        }
    ]

    response = requests.post(url, json=body, headers=headers)
    response.raise_for_status()


def main():
    login()
    create_oidc()


if __name__ == "__main__":
    main()
