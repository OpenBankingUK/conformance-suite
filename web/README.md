# web

## Contents

[TOC]

## Project setup

```sh
yarn install --frozen-lockfile --non-interactive
```

### Compiles and hot-reloads for development

```sh
yarn run serve
```

### Compiles and minifies for production

```sh
yarn run build
```

### Run your tests

```sh
yarn run test
```

### Lints and fixes files

```sh
yarn run lint
```

### Run your unit tests

```sh
yarn run test:unit
```

---

## UI

### Info

The UI is a SPA built with [Vue](https://vuejs.org/v2/guide/). The app is served by the [Phoenix Framework](https://phoenixframework.org/), Phoenix also handles the websocket.
Vue components are written using [single files](https://vuejs.org/v2/guide/single-file-components.html) (with the HTML template, JS and CSS).
Webpack transpiles all the files from modern JS and bundles the files for the browser.

#### State management

We are using [Vuex](https://vuex.vuejs.org/) as state management with modules to divide the store. Each module contains `actions`, `mutations`, `mutation-types`, `getters` and initial `state`. We are storing the store in `localStorage` (just the user info to handle the Google ID session) using [vuex-persistedstate](https://github.com/robinvdvleuten/vuex-persistedstate). This is where we set what to save in `localStorage`: [store/index.js](https://bitbucket.org/openbankingteam/conformance-suite/src/develop/apps/compliance_web/assets/src/store/index.js#lines-13).

#### SPA Routing

We are using the default router for Vue: [Vue Router](https://router.vuejs.org/)

#### Other Libraries

- `axios` and `vue-axios` as HTTP client [DOCS](https://github.com/axios/axios)
- `antd` as UI framework [DOCS](https://vuecomponent.github.io/ant-design-vue/docs/vue/introduce/)
- `brace` and `vue2-brace-editor` for the JSON editor [DOCS](https://github.com/Hector101/vue2-brace-editor)

#### Auth

We are using [Google Sign In](https://developers.google.com/identity/sign-in/web/sign-in) to authenticate the users on both Client and Server. On the server we verify the integrity of the ID token. There is a [vuex module](https://bitbucket.org/openbankingteam/conformance-suite/src/develop/apps/compliance_web/assets/src/store/modules/user/) to handle all the logic.
The router checks if a route is private using the `beforeEnter` method: [implementation](https://bitbucket.org/openbankingteam/conformance-suite/src/develop/apps/compliance_web/assets/src/router/index.js#lines-13:21)

#### Testing

The testing framework is [Jest](https://jestjs.io/docs/en/getting-started). The strategy is to add a file `*.spec.js` next to each file we want to test.
In the `__mocks__` folder there are 2 files: [fileMock.js](https://bitbucket.org/openbankingteam/conformance-suite/src/develop/apps/compliance_web/assets/__mocks__/fileMock.js) to globally mock all the import of static files (`css`, `svg`, etc...) in the tests, and [vue.js](https://bitbucket.org/openbankingteam/conformance-suite/src/develop/apps/compliance_web/assets/__mocks__/vue.js) to mock `Vue.axios` so we don't make any real http call in the tests.

## Getting started

### Install dependencies

In `apps/compliance_web/assets`, that is the root of the UI, run `npm i` to install all the node deps.

### Start the app

To start the app run `make serve_web` in the root of the application. Phoenix will serve the app in watch mode, that means that Webpack will re-compile the app when a file changes and refresh the browser at `http://localhost:4000`.

### Test the app

In a separate terminal you can run `npm t` in `apps/compliance_web/assets`. This will run all the tests. To run the tests in watch mode append `-- --watch`, for instance: `npm t -- --watch`. To see the code coverage you can append `-- --coverage`. It's possible to combine both flags: `npm t -- --watch --coverage`. You can also see the code coverage in the `coverage` folder: `open apps/compliance_web/assets/coverage/lcov-report/index.html` to open in a browser.

---

## assets

## inline fonts and images

`web/vue.config.js`: inlines the fonts and images into the app. If we don't want to inline the
fonts and images into the final app, simply remove this file.
