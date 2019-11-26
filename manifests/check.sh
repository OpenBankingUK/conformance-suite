#!/bin/sh
FILE="$(pwd)/ob_3.1_accounts_transactions_fca.json"

jq -r '.scripts[].parameters | select(has("x-fapi-interaction-id")) | .["x-fapi-interaction-id"]' "${FILE}" | go run check.go
