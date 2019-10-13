GOENV:=CGO_ENABLED=0

build-linux:
	$(GOENV) GOOS=linux GOARCH=amd64 go build $(GOFLAGS) -o bin/test-api cmd/main.go

build-docker-compose:
	docker-compose build

start: build-linux build-docker-compose
	docker-compose up -d

stop:
	docker-compose down
