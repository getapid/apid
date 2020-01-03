BIN:=bin/apid
OSARCHLIST:=`cat svc/cli/osarchlist`
GOFLAGS:=-ldflags "-X github.com/getapid/apid-cli/svc/cli/cmd.version=$(VERSION)"

all: test

test-api-%:
	$(MAKE) -C testapi $*

e2e: test-api-start e2e-test test-api-stop

build-e2e: test-api-build

build:
	go build $(GOFLAGS) -o $(BIN) svc/cli/main.go

release:
	./goreleaser release --snapshot

install:
	go mod download

mock:
	go generate ./...

test: mock
	go test $(GOFLAGS) ./...

e2e-test: build
	$(BIN) check -c testapi/tests/
