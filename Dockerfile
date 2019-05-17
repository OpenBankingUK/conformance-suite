# Image to compile go binaries
FROM golang:1.12-alpine as gobuilder
RUN apk update && apk add git make bash

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

# For caching technique, see: https://medium.com/@petomalina/using-go-mod-download-to-speed-up-golang-docker-builds-707591336888

# All these steps will be cached
RUN mkdir /app
WORKDIR /app
# COPY go.mod and go.sum files to the workspace
COPY go.mod .
COPY go.sum .
# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .

# Build the binary
RUN make build
RUN make build_cli

# Image to compile Single Page Application of the Vue.js site
FROM node:dubnium-alpine as nodebuilder
WORKDIR /app

# This is to that JavaScript code can import code defined in the Go side, e.g.,
# '../../../../../pkg/discovery/templates/ob-v3.1-generic.json'
# '../../../pkg/model/testdata/spec-config.golden.json'
ADD pkg/discovery/templates/*.json /pkg/discovery/templates/
ADD pkg/model/testdata/*.json /pkg/model/testdata/
COPY pkg/schema/spec/v3.0.0/*.json /pkg/schema/spec/v3.0.0/
COPY pkg/schema/spec/v3.1.0/*.json /pkg/schema/spec/v3.1.0/
ADD web .

ENV FORCE_COLOR=1
ENV NODE_DISABLE_COLORS=0

RUN yarn install \
	&& NODE_ENV=production yarn build

# Certificates needed if you are building a networking application
FROM alpine:latest as certs
RUN apk --update add ca-certificates

# Final image to run the binary
FROM scratch
LABEL MAINTAINER Open Banking
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /app

COPY --from=gobuilder /app/fcs_server /app/
COPY --from=gobuilder /app/fcs /app/
COPY --from=gobuilder /app/certs /app/certs
COPY --from=gobuilder /app/components /app/components
COPY --from=gobuilder /app/manifests /app/manifests
COPY --from=nodebuilder /app/dist /app/web/dist

COPY pkg/schema/spec/v3.0.0/*.json /app/pkg/schema/spec/v3.0.0/
COPY pkg/schema/spec/v3.1.0/*.json /app/pkg/schema/spec/v3.1.0/

EXPOSE 8443

ENTRYPOINT ["/app/fcs_server"]
