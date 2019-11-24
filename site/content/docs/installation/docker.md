+++
title = "Docker"
description = "How to use the APId docker image"
template = "docs/article.html"
weight = 2
sort_by = "weight"
+++

{{ h2(text="How to install") }}

The official APId docker image can be found at [docker hub](https://hub.docker.com/r/getapid/apid).
<br><br>

The entrypoint of the docker image is set to the APId executable. This makes it quite easy to use:
<br><br>
```sh
docker run -v /path/to/apid.yaml:/apid.yaml getapid/apid:latest check -c /apid.yaml
```