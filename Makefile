

default: test build

.PHONY: test
test:
	sh -c "'$(CURDIR)/scripts/test.sh'"

.PHONY: build
build:
	sh -c "'$(CURDIR)/scripts/build.sh'"
