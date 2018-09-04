appname := sawsh
version := $(shell git describe --tags --abbrev=0)

.PHONY: all clean build

default: clean build all package

build:
	mkdir -p build
	go get -d
	go test
	go build -o build/sawsh *.go

all:
	go get github.com/mitchellh/gox
	mkdir -p build
	gox \
		-ldflags=-s \
		-output="build/{{.OS}}_{{.Arch}}/sawsh"

clean:
	rm -rf build/

install:
	chmod +x build/sawsh
	sudo mv build/sawsh /usr/local/bin/sawsh

build_docker:
	docker run --rm -v "$(PWD)":/usr/src/myapp -w /usr/src/myapp golang:1.8 make build_all

package:
	$(shell rm -rf build/archive)
	$(shell rm -rf build/archive)
	$(eval UNIX_FILES := $(shell ls build | grep -v sawsh | grep -v windows))
	$(eval WINDOWS_FILES := $(shell ls build | grep -v sawsh | grep windows))
	@mkdir -p build/archive
	@for f in $(UNIX_FILES); do \
		echo Packaging $$f && \
		(cd $(shell pwd)/build/$$f && tar -czf ../archive/$$f.tar.gz sawsh*); \
	done
	@for f in $(WINDOWS_FILES); do \
		echo Packaging $$f && \
		(cd $(shell pwd)/build/$$f && zip ../archive/$$f.zip sawsh*); \
	done
	ls -lah build/archive/


