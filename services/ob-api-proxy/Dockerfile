FROM node:8.11.1-alpine

ENV NODE_ENV production

RUN apk update && apk upgrade && \
  apk add --no-cache jq make ca-certificates curl bash

RUN mkdir -p /home/node/app/ob-api-proxy
WORKDIR /home/node/app/ob-api-proxy

COPY . .

RUN npm install

CMD ./init.sh && npm start
