+++
title = "CircleCI"
description = "How to integrate APId with CircleCI"
template = "docs/article.html"
weight = 2
sort_by = "weight"
+++

{{ h4(text="Docker") }}

Using the official docker image is really simple, the only thing you need to do is use it as the base image for that job.
<br><br>
```yaml
---
version: 2.1
  jobs:
    test:
      docker: 
        - image: getapid/apid:latest
      steps:
        - run: apid check -c path/to/apid.yaml
```

{{ h4(text="CLI") }}

Integrating the CLI is just as simple as using the docker image. Download the latest version of the CLI, make it executable and run it.
<br><br>
```yaml
---
version: 2.1
  jobs:
    test:
      docker: 
        - image: your/image:latest
      steps:
        - run: |
            wget https://github.com/getapid/apid-cli/releases/download/v<version>/apid-<version>-linux-amd64.tar.gz
            tar -xzf apid-*.tzr.gz
            chmod u+x apid
            apid check -c path/to/apid.yaml
```
