+++
title = "Travis"
description = "How to integrate APId with Travis CI"
template = "docs/article.html"
weight = 2
sort_by = "weight"
+++

{{ h4(text="Docker") }}

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

{{ h4(text="CLI") }}

Integrating the CLI with Travis is the nois just as simple as using the docker image. Download the latest version of the CLI, make it executable and run it.
<br><br>

```yaml
---
before_install:
  - docker run -v /path/to/apid.yaml:/apid.yaml run getapid/apid:latest check -c /apid.yaml
  - wget https://github.com/getapid/apid-cli/releases/download/v<version>/apid-<version>-linux-amd64.tar.gz
  - tar -xzf apid-*.tzr.gz
  - chmod u+x apid
script:
  - apid check -c path/to/apid.yaml
```