# Image to compile go binaries
FROM golang:1.15-alpine as gobuilder
RUN apk add --no-cache --update --upgrade \
	bash \
	git \
	make

# 1. disable crosscompiling
# 2. compile linux only
# 3. target x64_64
#
# A normal compiled app is dynamically linked to the libraries it needs to run (i.e., all the C libraries it binds to).
# Unfortunately, scratch is empty, so there are no libraries and no loadpath for it to look in. What we have to do is modify our build script to statically compile our app with all libraries built in.
#
# https://github.com/AlessioCoser/minimal-docker-container-for-golang
ENV CGO_ENABLED=0
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
COPY pkg/discovery/templates/*.json /pkg/discovery/templates/
COPY pkg/model/testdata/*.json /pkg/model/testdata/
COPY pkg/schema/spec/v3.0.0/*.json /pkg/schema/spec/v3.0.0/
COPY pkg/schema/spec/v3.1.0/*.json /pkg/schema/spec/v3.1.0/
COPY pkg/schema/spec/v3.1.1/*.json /pkg/schema/spec/v3.1.1/
COPY pkg/schema/spec/v3.1.2/*.json /pkg/schema/spec/v3.1.2/
COPY pkg/schema/spec/v3.1.3/*.json /pkg/schema/spec/v3.1.3/
COPY pkg/schema/spec/v3.1.4/*.json /pkg/schema/spec/v3.1.4/
COPY pkg/schema/spec/v3.1.5/*.json /pkg/schema/spec/v3.1.5/
COPY pkg/schema/spec/v3.1.6/*.json /pkg/schema/spec/v3.1.6/
COPY web .

ENV FORCE_COLOR=1
ENV NODE_DISABLE_COLORS=0

RUN yarn install --frozen-lockfile --non-interactive \
	&& NODE_ENV=production yarn build

# Certificates needed if you are building a networking application
FROM alpine:latest as certs
RUN apk add --no-cache --update --upgrade ca-certificates

# Final image to run the binary
FROM alpine:3.9.4
RUN apk add --no-cache --update --upgrade \
	bash \
	coreutils \
	curl \
	emacs \
	git \
	jq \
	openssl \
	tree \
	wget \
	vim

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
COPY pkg/schema/spec/v3.1.1/*.json /app/pkg/schema/spec/v3.1.1/
COPY pkg/schema/spec/v3.1.2/*.json /app/pkg/schema/spec/v3.1.2/
COPY pkg/schema/spec/v3.1.3/*.json /app/pkg/schema/spec/v3.1.3/
COPY pkg/schema/spec/v3.1.4/*.json /app/pkg/schema/spec/v3.1.4/
COPY pkg/schema/spec/v3.1.5/*.json /app/pkg/schema/spec/v3.1.5/
COPY pkg/schema/spec/v3.1.6/*.json /app/pkg/schema/spec/v3.1.6/

EXPOSE 8443

ENTRYPOINT ["/app/fcs_server"]
