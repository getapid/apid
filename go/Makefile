GOFLAGS:=-v
BINDIR:=bin/apid

build-cli:
	go build $(GOFLAGS) -o $(BINDIR) svc/cli/main.go

install:
	go mod download

mock:
	go get github.com/golang/mock/mockgen
	go generate ./...

test: mock
	go test $(GOFLAGS) ./...

test-api-%: mock
	$(MAKE) -C testapi $*

e2e-test: build-cli
	$(BINDIR) check
