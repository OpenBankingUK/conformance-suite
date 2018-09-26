# Reference: https://github.com/Financial-Times/docker-elixir-build
FROM bitwalker/alpine-elixir:1.7 as build

WORKDIR /build

COPY mix.exs .
COPY mix.lock .

ARG MIX_ENV=prod
ARG APP_VERSION=0.0.0
ARG ENVFILE
ARG ENV
ARG DEBUG_ENVS
ENV ENV $ENV
ENV MIX_ENV ${MIX_ENV}
ENV APP_VERSION ${APP_VERSION}

RUN apk --no-cache add nodejs-npm

COPY $ENVFILE .
COPY apps apps
COPY config config
COPY rel rel

RUN echo '----' \
    && echo -e "\\033[92m ---> cat $ENVFILE ... \\033[0m" \
    && cat $ENVFILE && \
    echo '----'
RUN echo '----' \
    && echo -e "\\033[92m ---> cat $ENVFILE | grep -v '^\s*#' ... \\033[0m" \
    && cat $ENVFILE | grep -v '^\s*#' \
    && echo '----'

# Uncomment line below if you have assets in the priv dir
COPY apps/compliance_web/priv apps/compliance_web/priv

# Build Phoenix assets
RUN cd apps/compliance_web/assets \
    && npm install \
    && env `cat ../../../$ENVFILE | grep -v '^\s*#'` npm run build

RUN mix deps.get
RUN env `cat $ENVFILE | grep -v '^\s*#'` mix phx.digest
RUN env `cat $ENVFILE | grep -v '^\s*#'` mix release --env=${MIX_ENV}

# Container that will be used to run the final application
FROM alpine:3.8

RUN apk --no-cache update \
    && apk --no-cache upgrade \
    && apk --no-cache add ncurses-libs openssl bash ca-certificates

RUN adduser -D app

ARG MIX_ENV=prod
ARG APP_VERSION=0.0.0
ARG ENVFILE
ARG ENV
ARG DEBUG_ENVS
ENV ENV $ENV
ENV MIX_ENV ${MIX_ENV}
ENV APP_VERSION ${APP_VERSION}
ENV PORT 4000

WORKDIR /opt/app

COPY $ENVFILE .
RUN cat $ENVFILE | awk '!/^ *#/ && NF' | sed -e 's/^/export /' >> /etc/profile.d/envs.sh

# Copy release from build stage
COPY --from=build /build/_build/${MIX_ENV}/rel/* ./
USER app

# Mutable Runtime Environment
RUN mkdir /tmp/app
ENV RELEASE_MUTABLE_DIR /tmp/app
ENV START_ERL_DATA /tmp/app/start_erl.data

# Start command
CMD echo -e "\\033[92m  ---> sleeping (10 seconds) ... \\033[0m" \
    && sleep 10s \
    && /opt/app/bin/compliance_suite_server foreground
