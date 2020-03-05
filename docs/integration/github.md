# Github Actions

Integrating APId is simple because of the flexibility it offers - you can either use the official docker image, or if that doesn't suit your taste, you can use the CLI straight from your shell.

## Docker

Currently, there is no official APId action for github workflows, but you can always use the CLI \(see below\).

## CLI

Integrating the CLI is quite straight forward - download the latest version of the CLI, make it executable and run it.

```yaml
---
name: test
on: pull_request
jobs:
  test:
    name: end-to-end tests
    runs-on: ubuntu-latest
    steps:
      - name: run tests
        run: |
          curl -o apid https://cdn.getapid.com/cli/latest/apid-latest-linux-amd64
          chmod u+x apid
          apid check -c path/to/apid.yaml
```

