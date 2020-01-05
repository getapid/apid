+++
title = "Drone"
description = "How to integrate APId with Drone"
template = "docs/article.html"
weight = 2
sort_by = "weight"
+++

{{ h4(text="Docker") }}

Using the official docker image is really simple, the only thing you need to do is use it as the base image for that job.
<br><br>

```yaml
---
kind: pipeline
type: docker
name: default

steps:
- name: test
  image: getapid/apid:latest
  commands:
  - apid check -c path/to/apid.yaml
```

{{ h4(text="CLI") }}

Integrating the CLI is just as simple as using the docker image. Download the latest version of the CLI, make it executable and run it.
<br><br>

```yaml
---
kind: pipeline
type: docker
name: default

steps:
- name: test
  image: your/image:latest
  commands:
  - wget https://github.com/getapid/apid-cli/releases/download/v<version>/apid-<version>-linux-amd64.tar.gz
  - tar -xzf apid-*.tzr.gz
  - chmod u+x apid
  - apid check -c path/to/apid.yaml
```