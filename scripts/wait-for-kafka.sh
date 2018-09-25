#!/usr/bin/env bash
set -ueo pipefail

while ! nc -z localhost 9092; do
  echo -e "\\033[92m ---> waiting for kafka ... \\033[0m"
  sleep 1
done
