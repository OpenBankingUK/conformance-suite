# `E2E.md`
## run
The commands available are detailed below. In most cases you want to run `cypress` against the local server:.

Terminal 1:

```sh
make serve_web
```

Terminal 2
```sh
make test_integration_local
```

## ci
For CI we use:

```sh
make test_integration
```

## `cypress` references
* Good walkthrough with runnable code: <https://docs.cypress.io/examples/examples/tutorials.html>
* Installation guide: <https://docs.cypress.io/guides/getting-started/installing-cypress.html#>.
* How to integrate with CircleCI: <https://docs.cypress.io/examples/examples/docker.html#>.
