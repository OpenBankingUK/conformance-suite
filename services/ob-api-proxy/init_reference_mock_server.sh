#!/usr/bin/env bash
set -ueo pipefail
echo -e "\033[92m  ---> starting ob-api-proxy [init_reference_mock_server.sh] ... \033[0m"

echo -e "\033[92m  ---> envs ... \033[0m"
printenv

function wait_for_deps {
    while ! nc -z redis 6379; do
        echo -e "\033[92m  ---> waiting for redis ... \033[0m"
        sleep 1
    done
    echo "redis is UP"

    while ! nc -z mongo 27017; do
        echo -e "\033[92m  ---> waiting for mongo ... \033[0m"
        sleep 1
    done
    echo "mongo is UP"

    while ! nc -z reference-mock-server 8001; do
        echo -e "\033[92m  ---> waiting for reference-mock-server ... \033[0m"
        sleep 1
    done
    echo "reference-mock-server is UP"
}

function auth_servers_init {
    echo -e "\033[92m  ---> updating auth servers and openids ... \033[0m";
    npm run updateAuthServersAndOpenIds
}

function auth_servers_credentials {
    echo -e "\033[92m  ---> saving credentials ... \033[0m";
    # reference-mock-server auth servers
    npm run saveCreds authServerId=aaaj4NmBD8lQxmLh2O clientId=spoofClientId clientSecret=spoofClientSecret
    npm run saveCreds authServerId=bbbX7tUB4fPIYB0k1m clientId=spoofClientId clientSecret=spoofClientSecret
    npm run saveCreds authServerId=cccbN8iAsMh74sOXhk clientId=spoofClientId clientSecret=spoofClientSecret
}

wait_for_deps
auth_servers_init
auth_servers_credentials
