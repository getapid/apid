all: go-build

go-%:
	$(MAKE) -C go $*

e2e: go-test-api-start go-e2e-test go-test-api-stop
