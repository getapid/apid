GOFLAGS:=-v
build:
	go build $(GOFLAGS) -o bin/apid svc/cli/main.go

run: build
	go run $(GOFLAGS) svc/cli/main.go check 

install:
	go mod download

test:
	go test $(GOFLAGS) ./...
