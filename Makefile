all: go-build

go-%:
	$(MAKE) -C go $*

include go/Makefile
