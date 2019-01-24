# `GO_TOOLS`
## `vscode`
### `golangci-lint`
To configure [`golangci-lint`](https://github.com/golangci/golangci-lint) in `vscode`.

#### install
```sh
# binary will be $(go env GOPATH)/bin/golangci-lint
curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.12.5
```

See:
* [GolangCI-Lint: CI Installation](https://github.com/golangci/golangci-lint#ci-installation)
* `devtools` in [`../Makefile`](../Makefile)

#### `settings.json`
```json
...
  "go.lintTool": "golangci-lint",
  "go.lintFlags": [
    "--config=${workspaceFolder}/.golangci.yml",
    "--fast",
  ],
...
```
