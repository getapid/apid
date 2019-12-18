BIN:=bin/apid
OSARCHLIST:=`cat svc/cli/osarchlist`
GOFLAGS:=-ldflags "-X github.com/getapid/apid/svc/cli/cmd.version=$(VERSION)"

build-cli:
	go build $(GOFLAGS) -o $(BIN) svc/cli/main.go

build-cli-release:
	gox -osarch="$(OSARCHLIST)" -output="./bin/build/apid-$(VERSION)-{{.OS}}-{{.Arch}}" $(GOFLAGS) ./svc/cli/
	gox -osarch="$(OSARCHLIST)" -output="./bin/latest/apid-latest-{{.OS}}-{{.Arch}}" $(GOFLAGS) ./svc/cli/

install:
	go mod download

mock:
	go generate ./...

test: mock
	go test $(GOFLAGS) ./...

test-api-%:
	$(MAKE) -C testapi $*

e2e-test: build-cli
	$(BIN) check -c testapi/tests/
