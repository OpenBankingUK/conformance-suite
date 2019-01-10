#!/usr/bin/env bash
echo -e '\033[92m  ---> Starting server ... \033[0m'
PORT=8080 go run cmd/server/main.go
