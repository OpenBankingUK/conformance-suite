[![CircleCI](https://circleci.com/gh/OpenBankingUK/compliance-suite-server.svg?style=svg&circle-token=7042965fb04fa83d7cafa5b2d43a2c0f0febabf6)](https://circleci.com/gh/OpenBankingUK/compliance-suite-server)

# Functional Conformance Suite

## Technical Overview

You can read an [overview of our suite umbrella project and process flow here](./apps/README.md).

## Cloning the repos

You can git clone the repositories as follows:

```sh
$ git clone https://github.com/OpenBankingUK/compliance-suite-server.git
$ git clone https://github.com/OpenBankingUK/reference-mock-server.git
```

## Getting Started
- Install `elixir`, `docker` and `nodeJS`
```sh
brew install elixir node
```

On a mac to install `docker` you can follow this [page](https://docs.docker.com/docker-for-mac/install/) or run:
```
brew cask install docker
```

- Copy `.env` from `.env.sample` in the `root` and in `services/ob-api-proxy`

## Commands
Generate a new secret:
```sh
mix phx.gen.secret
# Output: kro6Z/FRTZhnsjw0TTSs4tqZNmUv6zLJOmWWj7g0m1Fgp2ZOSzfo6cImpceDH1vD
```

Start local server:
```sh
make serve_web
```

Start `ob-api-proxy` running (in isolation, `make serve_web` starts it in background):
```sh
cd services/ob-api-proxy
npm i
npm run update # Add Auth Servers and save credentials
npm run dev # Start the node app in watch mode
```

Start local server in Docker:
```sh
make serve_web_docker
```

You can go to [localhost:4000](http://localhost:4000) and try it out!

## Local development

Run tests automatically when files change:

```sh
mix test.watch
```

Run Elixir formatter to auto-format files (configuration is in root and
/apps/*/ `.formatter.exs` files):

```sh
mix format
```

Run credo static code analysis checks (configuration is in `.credo.exs` file):

```sh
mix credo
```

If you use [Atom](https://atom.io/) editor, formatter and credo checks can be
run in editor on file save by installing these packages:

```sh
apm install atom-elixir-formatter
apm install linter-elixir-credo
```

## Deployment & CI

### Deployment

TODO

### CI (continuous deployment)
Our CI server is [CircleCI](https://circleci.com) and is configured in [.circleci/config.yml](.circleci/config.yml).
Also, it continuously deploys to our hosted AWS instance on a successful merge to master on Github.

### Debugging

#### Check containers status
Execute the command to get logs into the service:
```sh
$ docker-compose ps
```

#### Accessing container logs
Execute the command to get logs into the service:
```sh
$ docker-compose logs -f

# Or for a single service:
$ docker-compose logs -f reference-mock-server
$ docker-compose logs -f ob-api-proxy
```
