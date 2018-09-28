#!/usr/bin/env bash
set -eo pipefail

on_exit() {
  # https://stackoverflow.com/questions/360201/how-do-i-kill-background-processes-jobs-when-my-shell-script-exits
  JOBS="$(jobs -pr)"

  echo -e "\\033[92m  ---> web-serve.sh: killing background JOBS=${JOBS}  \\033[0m"
  for JOB in ${JOBS}; do
    echo -e "\\033[92m  ---> web-serve.sh: killing JOB=${JOB}  \\033[0m"
    kill -s SIGTERM "$JOB" || sleep 1 && kill -9 "$JOB"
  done
}
trap on_exit EXIT

export REFERENCE_HOST=localhost

init_server() {
  echo -e "\\033[92m  ---> web-serve.sh: starting server  \\033[0m"
  make compile_server_local
  make init_server_local
}

wait_for_deps() {
  echo -e "\\033[92m  ---> web-serve.sh: wait_for_deps  \\033[0m"

  while ! nc -z localhost 6379; do
    echo -e "\\033[92m  ---> web-serve.sh: waiting for redis ...  \033[0m"
    sleep 1
  done

  while ! nc -z localhost 27017; do
    echo -e "\\033[92m  ---> web-serve.sh: waiting for mongo ...  \\033[0m"
    sleep 1
  done

  while ! nc -z localhost 8001; do
    echo -e "\\033[92m  ---> web-serve.sh: waiting for reference-mock-server ...  \\033[0m"
    sleep 1
  done
}

# Start ob-api-proxy in dev mode with watch
start_ob_proxy() {
  # cd services/ob-api-proxy && npm i && npm run update && npm run dev
  cd services/ob-api-proxy && npm install && npm run dev
}

echo -e "\\033[92m  ---> web-serve.sh: starting ...  \\033[0m"

# Build the image if it doesn't exist
if [[ "$(docker images -q openbankinguk/reference-mock-server:latest 2> /dev/null)" == "" ]]; then
  cd ../reference-mock-server && docker-compose build
fi

# Stop containers if previously started with make serve_web_docker as they use same PORT
# Do not crash if containers are not running
docker-compose stop compliance-suite-server ob-api-proxy || true

# Start services
if [[ -n "$RECREATE" ]]; then
  echo -e "\\033[92m  ---> web-serve.sh: recreating services ... \\033[0m"
  time docker-compose up \
    -d \
    --force-recreate \
    --always-recreate-deps \
    --renew-anon-volumes \
    --remove-orphans \
    mongo redis reference-mock-server
else
  time docker-compose up \
    -d \
    mongo redis reference-mock-server
fi

wait_for_deps
echo -e "\\033[92m  ---> web-serve.sh: tail -F -n +1 $(pwd)/ob-api-proxy.log \\033[0m"
# Run ob-api-proxy in background
start_ob_proxy > ./ob-api-proxy.log 2>&1 &

init_server
make start_server_local
