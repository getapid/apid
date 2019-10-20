GOFLAGS:=-v
BINDIR:=bin/apid
build-cli:
	go build $(GOFLAGS) -o $(BINDIR) svc/cli/main.go

install:
	go mod download

test:
	go test $(GOFLAGS) ./...

test-api-%:
	$(MAKE) -C testapi $*

e2e-test: build-cli
	$(BINDIR) check
