# Docker

The official APId docker image can be found at [docker hub](https://hub.docker.com/r/getapid/apid).

```sh
docker pull getapid/apid:latest
```

The entrypoint of the docker image is set to the APId executable. This makes it quite easy to use:

```sh
docker run -v /path/to/apid.yaml:/apid.yaml getapid/apid:latest check -c /apid.yaml
```
