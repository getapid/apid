FROM alpine:latest
ENTRYPOINT ["/usr/local/bin/apid"]
RUN apk --no-cache add ca-certificates
COPY apid /usr/local/bin/apid