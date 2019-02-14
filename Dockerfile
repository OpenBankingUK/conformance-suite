FROM alpine:latest as certs
RUN apk --update add ca-certificates

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

RUN make init
RUN make build
RUN make build_cli

# Image to compile Single Page Application of the Vue.js site
FROM node:8.11.1-slim as nodebuilder
WORKDIR /app

# This is to that JavaScript code can import code defined in the Go side, e.g.,
# '../../../../../pkg/discovery/templates/ob-v3.1-generic.json'
# '../../../pkg/model/testdata/spec-config.golden.json'
ADD pkg/discovery/templates/*.json /pkg/discovery/templates/
ADD pkg/model/testdata/*.json /pkg/model/testdata/
ADD web .

RUN yarn install \
	&& NODE_ENV=production yarn build

# Final image to run the binary
FROM scratch
LABEL MAINTAINER Open Banking
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /app

COPY --from=gobuilder /app/fcs_server /app/
COPY --from=gobuilder /app/fcs /app/
COPY --from=gobuilder /app/certs /app/certs
COPY --from=gobuilder /app/components /app/components
COPY --from=nodebuilder /app/dist /app/web/dist

EXPOSE 8443

ENTRYPOINT ["/app/fcs_server"]
