GOFLAGS:=-v
BIN:=bin/apid
VERSION:=`cat ../version`
OSARCHLIST:=`cat svc/cli/osarchlist`

build-cli-local:
	go build $(GOFLAGS) -ldflags "-X github.com/getapid/apid/svc/cli/cmd.version=$(VERSION)" -o $(BIN) svc/cli/main.go

build-cli-release:
	echo $(OSARCHLIST)
	gox -osarch="$(OSARCHLIST)" -output="./bin/build/apid-$(VERSION)-{{.OS}}-{{.Arch}}" ./svc/cli/
	gox -osarch="$(OSARCHLIST)" -output="./bin/latest/apid-latest-{{.OS}}-{{.Arch}}" ./svc/cli/

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
