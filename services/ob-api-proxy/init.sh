#!/usr/bin/env bash
set -eo pipefail

echo -e "\033[92m  ---> starting ob-api-proxy [init.sh] ... \033[0m"

if [[ "$INIT_REFERENCE_MOCK_SERVER" == true ]]; then
  echo -e "\033[92m  ---> using init_reference_mock_server.sh ... \033[0m"
  /home/node/app/ob-api-proxy/init_reference_mock_server.sh
else
  echo -e "\033[92m  ---> using init_ob_directory_sandbox.sh ... \033[0m"
  /home/node/app/ob-api-proxy/init_ob_directory_sandbox.sh
fi
