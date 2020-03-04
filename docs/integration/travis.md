# TravisCI

Integrating APId is simple because of the flexibility it offers - you can either use the official docker image, or if that doesn't suit your taste,
you can use the CLI straight from your shell.

## Docker

Due to the nature of Travis, using a docker image has no real benefits over using the CLI. That being said, you can still use it.
<br><br>

```yaml
---
services:
  - docker
before_install:
  - docker pull getapid/apid:latest
script:
  - docker run -v /path/to/apid.yaml:/apid.yaml run getapid/apid:latest check -c /apid.yaml
```

## CLI

Integrating the CLI with Travis is the nois just as simple as using the docker image. Download the latest version of the CLI, make it executable and run it.
<br><br>

```yaml
---
before_install:
  - docker run -v /path/to/apid.yaml:/apid.yaml run getapid/apid:latest check -c /apid.yaml
  - curl -o apid https://cdn.getapid.com/cli/latest/apid-latest-linux-amd64
  - chmod u+x apid
script:
  - apid check -c path/to/apid.yaml
```
