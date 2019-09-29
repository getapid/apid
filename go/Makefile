GOFLAGS:=-v
build:
	go build $(GOFLAGS) svc/cli/main.go

run: build
	go run $(GOFLAGS) svc/cli/main.go

install:
	go mod download

test:
	go test $(GOFLAGS) ./...
