# ☁️ APId - end-to-end API testing framework

[![Release cli](https://github.com/getapid/apid-cli/actions/workflows/release.yml/badge.svg)](https://github.com/getapid/apid-cli/actions/workflows/release.yml)
[![Coverage Status](https://coveralls.io/repos/github/getapid/cli/badge.svg?branch=main)](https://coveralls.io/github/getapid/cli?branch=main)
[![GitHub release](https://badgen.net/github/release/getapid/cli)](https://github.com/getapid/cli/releases)

## 🔭 Overview

APId is a framework that helps you write declarative, end-to-end tests for your API. The tests are written in YAML using simple, yet powerful syntax. Model tests around user flows - APId lets you verify sequential API calls as one flow.

## ⚡️ Quickstart

The quickest way to run APId is to use the docker image found [here](https://hub.docker.com/r/getapid/apid).

```bash
docker pull getapid/apid:latest
```

APId expects a file (or folder) with the defined tests to run, so as an example lets use the following basic test:

```yaml
variables:
  url: "https://httpbin.org"
transactions:
  - id: simple-request
    steps:
      - id: get
        request:
          method: "GET"
          endpoint: "{{ var.url }}/get"
        expect:
          code: 200
```

Make a new file and name it `apid.yaml` with the content above. To execute it, we need to mount it on the docker container running APId:

```bash
docker run -v "$(pwd)/apid.yaml:/apid.yaml" getapid/apid:latest check -c /apid.yaml
```

You should see the following in your terminal:

```bash
simple-request:
    OK   get
successful transactions: 1/1
failed transactions:     0/1
```

Success! You've just written your first APId test! If you change the `expect.code` from `200` to lets say `500` the test will fail and this will be the output:

```bash
simple-request:
    FAIL get

         request: GET https://httpbin.org/get
         errors:
             code:
                want 500, received 200

successful transactions: 0/1
failed transactions:     1/1
```

For more examples please check the `examples` folder in this repository.

## 📚 Documentation

You can find all APId documentation [here](https://docs.getapid.com/). The content of the docs site is located under `docs` in this repository.

## ⬇️ Installation

APId comes in two distributions - command line interface and docker image. You can use whichever you prefer.

### 💻 CLI

1. Head to the [latest release](https://github.com/getapid/apid-cli/releases/latest) and select the binary for your operating system.
2. Once downloaded open the archive and place the executable in a directory on your `$PATH`.

### ☁️ Docker

The official docker repository is located [here](https://hub.docker.com/r/getapid/apid). To pull the latest image run this in your terminal:

```bash
docker pull getapid/apid:latest
```

For more information on the docker image please visit our [documentation](https://docs.getapid.com/).

## 💻 Shell completion

Currently `apid` can generate shell completion for `bash`, `zsh` and `powershell`

```bash
apid completion bash -f /etc/bash.completion.d/apid.sh
```

## ⚙️ Development

All useful development commands can be found in the Makefile. Follow these simple steps to build and test the CLI locally:

### 1. Install Mockgen

Mockgen is used to generate mock implementations for testing

```bash
go get github.com/golang/mock/mockgen
```

### 2. Building and running CLI tests

```bash
make
```

### 3. Running CLI end-to-end tests

```bash
make e2e
```

### 4. Update the docs

The docs are located in the `docs` folder. The docs site automatically pulls the latest master version of the docs.

## 👽 Contributing

To contribute to APId, please see [CONTRIBUTING](CONTRIBUTING.md).

For questions and discussion join our [FAQ page](https://faq.getapid.com).
