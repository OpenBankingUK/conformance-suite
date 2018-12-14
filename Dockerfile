# Image to compile go binaries
FROM golang:1.11-stretch as gobuilder
# disable crosscompiling
#
# A normal compiled app is dynamically linked to the libraries it needs to run (i.e., all the C libraries it binds to).
# Unfortunately, scratch is empty, so there are no libraries and no loadpath for it to look in. What we have to do is modify our build script to statically compile our app with all libraries built in.
#
# https://github.com/AlessioCoser/minimal-docker-container-for-golang
ENV CGO_ENABLED=0
# compile linux only
ENV GOOS=linux
ENV GOARCH=amd64

WORKDIR /app

ADD . .

RUN make build
RUN make build_cli

# Image to compile Single Page Application of the Vue.js site
FROM node:8.11.1-slim as nodebuilder
WORKDIR /app

# This is to that JavaScript code can import the example discovery config
# defined in the Go side using: '../../../../../pkg/discovery/templates/ob-v3.0-ozone.json'
ADD pkg/discovery/templates/*.json /pkg/discovery/templates/
ADD web .

RUN yarn install \
	&& NODE_ENV=production yarn build

# Final image to run the binary
FROM scratch
LABEL MAINTAINER Open Banking
WORKDIR /app

COPY --from=gobuilder /app/server /app/
COPY --from=gobuilder /app/fcs /app/
COPY --from=nodebuilder /app/dist /app/web/dist

EXPOSE 8080

ENTRYPOINT ["/app/server"]
