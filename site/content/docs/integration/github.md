+++
title = "Github"
description = "How to integrate APId with Github workflows"
template = "docs/article.html"
weight = 2
sort_by = "weight"
+++

{{ h4(text="Docker") }}

Currently, there is no official APId action for github workflows, but you can always use the CLI (see below).

{{ h4(text="CLI") }}

Integrating the CLI is quite straight forward - download the latest version of the CLI, make it executable and run it.
<br><br>

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
            wget https://github.com/getapid/apid-cli/releases/download/v<version>/apid-<version>-linux-amd64.tar.gz
            tar -xzf apid-*.tzr.gz
            chmod u+x apid
            apid check -c path/to/apid.yaml
```