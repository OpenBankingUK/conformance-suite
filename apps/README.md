# Functional Conformance Suite Umbrella Project

## Overview

Elixir umbrella projects are a convenience to help you organise and manage
multiple applications. While it provides a degree of separation between
applications, those applications are not fully decoupled. They share the same
configuration and the same dependencies.

The pattern of keeping multiple applications in the same repository is known as
“mono-repo”. Umbrella projects maximise this pattern by providing conveniences
to compile, test and run multiple applications at once.

## Applications in our umbrella

Currently we have the following applications in our umbrella project:

#### apps/compliance

- The core Elixir API for our application's functionality.
- Talks to `ob_api_remote` application Elixir API to make OB API calls.

#### apps/compliance_web

- Provides an HTTP API for client use - using [Phoenix Framework](https://phoenixframework.org/).
- Provides socket based communication for client use - using [Phoenix channels](https://hexdocs.pm/phoenix/channels.html#content).
- Talks to `compliance` application Elixir API to provide functionality.
- Serves our client [Vue.js](https://vuejs.org/) application.
- Client application code is in `compliance_web/assets` directory.

#### apps/log_consumer

- Consumes validation run log events from Kafka log.
- Sends log events via `compliance` application to an `Aggregate` process for report generation.
- We use [GenStage](https://elixir-lang.org/blog/2016/07/14/announcing-genstage/#genstage) to exchange events with back-pressure between Elixir processes.
- Back-pressure is a mechanism to handle a fast producer/slow consumer situation.

#### apps/ob_api_remote

- An Elixir API wrapper around HTTP calls to `ob-api-proxy` Node.js server proxy to OB API.
- Used from `compliance` application.

## Services outside umbrella

#### services/ob-api-proxy

- An HTTP API wrapper around calls to OB API.
- Provides an API facade over authentication/authorisation OB calls and stores access_tokens.
- Provides an API proxy to OB API account resource and single immediate payment endpoints.
- Written in Node.js.

## Processes overview

In Elixir, all code runs inside [processes](https://elixir-lang.org/getting-started/processes.html).  

Elixir processes:
- are isolated from each other
- run concurrent to one another
- communicate via message passing
- are not operating system processes
- are extremely lightweight in terms of memory and CPU
- not uncommon to have 10,000 to 100,000 processes running simultaneously
- are the basis for concurrency in Elixir
- provide the means for building distributed and fault-tolerant programs

### Validation Run process flow

A user initiates a validation run from the client, and sees a dynamic report of the results in the client.

Here's what happens during that process:
1. Client makes validation run HTTP request to `compliance_web` app
1. We call `compliance` Elixir API to asynchronously initiate runs
1. We start a `ValidationRun` process - keyed by `validation_run_id`
1. We return `validation_run_id` to client
1. Our `ValidationRun` process - based on user config -
  1. Finds permutations config to use for Accounts API tests.
  1. Determines which swagger files to use for validation.
  1. Coordinates a number of read/write API endpoint calls
  1. Calls the `ob_api_remote` Elixir API to initiate endpoint calls.
1. The `ob_api_remote` makes HTTP calls to the `ob-api-proxy` facade.
1. The `ob-api-proxy` Node server:
  1. Calls ASPSP auth server API HTTP endpoints.
  1. Stores access_tokens.
  1. Calls resource server API HTTP endpoints.
  1. Validates responses against provided swagger file.
  1. Writes log entry to Kafka with each request/response validation.
  1. Revokes access_token ready for next request.
1. The `log_consumer` application consumes log entries from Kafka:
  1. We set up a GenStage flow, with a producer of log entries
  1. And a consumer that adds log entries to an `Aggregate` process
  1. It starts an `Aggregate` process for a new `validation_run_id`
  1. `Aggregate` process updates endpoint reports state.
  1. Todo: We intend `Aggregate` process will persist report state.
1. Client joins via `compliance_web` a `ReportChannel` socket
1. We start a new `ReportChannel` process:
  1. We set `ReportChannel` process as a listener on an `Aggregate` process
  1. `ReportChannel` pushes report state to client, on notification of an update
1. Client listens for report updates for given `validation_run_id`
1. Client renders report dynamically as updates are received.

### Observing Processes

During development, you can open the Erlang observer tool like this:

```
cd conformance-suite
iex -S mix phx.server

iex> :observer.start()
```

You'll see the observer GUI application open.

To see a digram of running processes:

- Click on the `Applications` tab
- Click on the application name in the left pane, e.g. `compliance`
- You'll see diagram of running processes
