#!/usr/bin/env bash
if [[ -z "$ENV" ]]; then
  ENVFILE=".env"
else
  ENVFILE=".env.${ENV}"
fi

echo -e "\\033[92m ---> create-envs, pwd=$(pwd), ENV=${ENV}, ENVFILE=${ENVFILE} ... \\033[0m"

ENVS=(
  DATA_DB_HOST
  DATA_DB_NAME
  ENDPOINT_URL_HOST
  ENDPOINT_URL_PORT
  GOOGLE_OAUTH_CLIENT_ID
  GOOGLE_OAUTH_TOKENINFO_URL
  GOOGLE_OAUTH_CLIENT_SECRET
  GUARDIAN_SECRET_KEY
  KAFKA_HOST
  KAFKA_PORT
  OB_API_PROXY_URL
  VALIDATION_KAFKA_BROKER
  VALIDATION_KAFKA_TOPIC
)

if [ ! -f $ENVFILE ]; then
  for key in "${ENVS[@]}"; do
    value=$key
    eval pair="$key=\$$value"
    echo $pair >> $ENVFILE
  done
fi

export $(cat $ENVFILE | grep -v '^\s*#' | xargs)
