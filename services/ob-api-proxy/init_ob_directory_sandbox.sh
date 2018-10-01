#!/usr/bin/env bash
set -eo pipefail
echo -e "\033[92m  ---> starting ob-api-proxy [init_ob_directory_sandbox.sh] ... \033[0m"

echo -e "\033[92m  ---> envs ... \033[0m"
printenv

if [[ -z "$OZONE_CLIENT_ID" ]]; then
  echo -e "\033[91m  ---> init_ob_directory_sandbox.sh: OZONE_CLIENT_ID is not set \033[0m" 1>&2
  exit 1
fi
if [[ -z "$OZONE_CLIENT_SECRET" ]]; then
  echo -e "\033[91m  ---> init_ob_directory_sandbox.sh: OZONE_CLIENT_SECRET is not set \033[0m" 1>&2
  exit 1
fi

function wait_for_deps {
    while ! nc -z mongo 27017; do
        echo -e "\033[92m  ---> waiting for mongo ... \033[0m"
        sleep 1
    done
    echo "mongo is UP"
}

function auth_servers_init {
    echo -e "\033[92m  ---> updating auth servers and openids ... \033[0m";
    npm run updateAuthServersAndOpenIds
}

function auth_servers_credentials {
    echo -e "\033[92m  ---> saving credentials for 3iPABZImMFEND0b9ZxSuNC ... \033[0m";
    # The values of `clientId` and `clientSecret` are obtained the dynamic registration call.
    # This is, in order:
    # 1. Get the auth server config.
    # 2. Read the value of `OpenIDConfigEndPointUri`.
    # 3. Call `https://modelobankauth2018.o3bank.co.uk:4101/.well-known/openid-configuration` and get the value of `registration_endpoint`. You can do this using `curl -s https://modelobankauth2018.o3bank.co.uk:4101/.well-known/openid-configuration | jq '.registration_endpoint'`
    # 4. Call `https://modelobank2018.o3bank.co.uk:4501/reg` and register the client to get a payload like below.
    #
    # TODO: We should refactor this code a bit to populate the code dynamically.
    npm run saveCreds authServerId=3iPABZImMFEND0b9ZxSuNC clientId="$OZONE_CLIENT_ID" clientSecret="$OZONE_CLIENT_SECRET"
}

wait_for_deps
auth_servers_init
auth_servers_credentials
