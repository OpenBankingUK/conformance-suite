### Installation via Docker - for quick start with mocked API

To install as a container-based app we assume
[Docker](https://www.docker.com/community-edition) ver17.12+ is installed.

If not installed you can find [Docker Community Edition downloads here](https://www.docker.com/community-edition#/download).

__BEFORE PROCEEDING FURTHER__ install the
[reference mock server via Docker](https://github.com/OpenBankingUK/reference-mock-server),

Then proceed with installing this server via Docker as follows.

```sh
cd tpp-reference-server
docker build -t ob/tpp-server --build-arg TAG_VERSION=v0.x.0 .
```

Use `docker images` command to check `ob/tpp-server` has been created:

```sh
docker images
# REPOSITORY                 TAG                 ...
# ob/tpp-server              latest              ...
# ob/reference-mock-server   latest              ...
# node                       8.4-alpine          ...
```

Use `docker-compose up` to install `mongo` and `redis` images, and run both
the reference mock server and the TPP reference server:

```sh
ASPSP_AUTH_HOST_IP=localhost docker-compose up
```

The TPP reference server should now be running on localhost:8003.


To open shell in running container if needed:

```sh
docker ps | grep tpp-server # to find CONTAINER_ID
open shell > docker exec -it [CONTAINER_ID] bash
```
