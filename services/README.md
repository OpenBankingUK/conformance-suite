# Services

These are service dependencies mostly written in different languages and provisioned as GIT submodules.

Please refer to the [GIT submodule documentation](https://git-scm.com/book/en/v2/Git-Tools-Submodules), to understand how to checkout, switch branches and get the latest updates for a submodule.

## ob-api-proxy

This is the [TPP Reference Server](https://github.com/OpenBankingUK/tpp-reference-server.git) as a submodule.

It's used to make all OB Read/Write API calls, perform some basic swagger validation logging all requests and responses for further processing.

This service is developed using [NodeJS](http://nodejs.org) and [Express](http://expressjs.com).

### Using the service

Add `OB_API_PROXY_URL` ENV the relevant `config/[dev.exs|test.exs|prod.exs]` config files in order to use the service in the `/apps` directory.

For example:

* DEV: `OB_API_PROXY_URL=http://localhost:8003`
* TEST: `OB_API_PROXY_URL=http://localhost:8003`
* PROD: `OB_API_PROXY_URL=http://<ip address|domain name>/ob-api-proxy/`

### Cloning the service

If you haven't initially cloned this correctly when cloning the root project, just run:

```sh
cd compliance-suite-server
git submodule update --init --recursive
```

### Installing dependencies

```sh
cd services/ob-api-proxy
npm i
```

### Running the service

#### Development

Setup ENVs.

In dev:

```sh
cp .env.sample .env
```

Then run using

```sh
npm run foreman
```

#### Production

We use [Nanobox](https://nanobox.io) to setup and start the service as a docker container in production.

Nanobox provides all monitoring and logging for the service.

Ensure all ENVs in the `.env` have been setup using `nanobox evar add <prod|dry-run> ENV_NAME=value`.

To test the deployment locally (DOES NOT DEPLOY TO PRODUCTION):
```sh
nanobox deploy dry-run --force --debug
```

The service should be accessible from `http://<Nanobox IP Address>/ob-api-proxy/...`
