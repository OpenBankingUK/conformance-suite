# syntax=docker/dockerfile:1

# Image to compile Go binaries
FROM --platform=$BUILDPLATFORM golang:1.21-alpine AS gobuilder

ARG TARGETOS
ARG TARGETARCH

RUN apk add --no-cache --update --upgrade \
    bash \
    git \
    make

ENV CGO_ENABLED=0

RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} make build
RUN GOOS=${TARGETOS} GOARCH=${TARGETARCH} make build_cli


# Image to compile Single Page Application
FROM --platform=$BUILDPLATFORM node:14-alpine AS nodebuilder

WORKDIR /app

COPY pkg/discovery/templates/*.json /pkg/discovery/templates/
COPY pkg/model/testdata/*.json /pkg/model/testdata/
COPY pkg/schema/spec/ /pkg/schema/spec/
COPY web .

ENV FORCE_COLOR=1
ENV NODE_DISABLE_COLORS=0

RUN yarn install --frozen-lockfile --non-interactive \
    && NODE_ENV=production yarn build


FROM alpine:latest AS certs

RUN apk add --no-cache --update --upgrade ca-certificates


FROM alpine:3.18.6

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

LABEL MAINTAINER="Open Banking"

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /app

COPY --from=gobuilder /app/fcs_server /app/
COPY --from=gobuilder /app/fcs /app/
COPY --from=gobuilder /app/certs /app/certs
COPY --from=gobuilder /app/components /app/components
COPY --from=gobuilder /app/manifests /app/manifests
COPY --from=nodebuilder /app/dist /app/web/dist

COPY pkg/schema/spec/ /app/pkg/schema/spec/

EXPOSE 8443

ENTRYPOINT ["/app/fcs_server"]
