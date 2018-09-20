#!/usr/bin/env bash
set -ueox pipefail

update_auth_servers() {
    echo -e "\033[92m  ---> updating auth servers and openids ... \033[0m";
    npm run updateAuthServersAndOpenIds
}

save_creds() {
    echo -e "\033[92m  ---> saving credentials ... \033[0m";
    # reference-mock-server auth servers
    npm run saveCreds authServerId=aaaj4NmBD8lQxmLh2O clientId=spoofClientId clientSecret=spoofClientSecret
    npm run saveCreds authServerId=bbbX7tUB4fPIYB0k1m clientId=spoofClientId clientSecret=spoofClientSecret
    npm run saveCreds authServerId=cccbN8iAsMh74sOXhk clientId=spoofClientId clientSecret=spoofClientSecret
}

"$@"
