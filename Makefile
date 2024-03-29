BIN:=bin/apid
OSARCHLIST:=`cat svc/cli/osarchlist`
GOFLAGS:=-ldflags "-X github.com/getapid/apid/svc/cli/cmd.version=$(VERSION)"

all: test

build:
	go build $(GOFLAGS) -o $(BIN) main.go

test: test/positive test/negative

test/ci: 
	./scripts/test.sh $(BIN)

test/positive: build
	@docker-compose -f tests/echo/docker-compose.yaml up -d &>/dev/null
	$(BIN) check -s "tests/**/*_pass.jsonnet" --silent
	@docker-compose -f tests/echo/docker-compose.yaml down  &>/dev/null

test/negative: build
	@docker-compose -f tests/echo/docker-compose.yaml up -d  &>/dev/null
	$(BIN) check -s "tests/**/*_fail.jsonnet" --silent
	@docker-compose -f tests/echo/docker-compose.yaml down  &>/dev/null

test/local: build
	@docker-compose -f tests/echo/docker-compose.yaml up -d  &>/dev/null
	$(BIN) check -s "$(target)"
	@docker-compose -f tests/echo/docker-compose.yaml down  &>/dev/null