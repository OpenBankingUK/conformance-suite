#!/usr/bin/env bash
set -ueo pipefail

echo -e "\\033[92m ---> building openbankinguk/reference-mock-server ... \\033[0m"
if [[ "$(docker images -q openbankinguk/reference-mock-server:latest 2> /dev/null)" == "" ]]; then
  cd ./services \
    && git clone https://github.com/OpenBankingUK/reference-mock-server.git \
    && cd reference-mock-server \
    && docker-compose build
fi
