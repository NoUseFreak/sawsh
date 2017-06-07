default:
	mkdir -p build
	go get -d
	go build -o build/sawsh sawsh.go

install:
	chmod +x build/sawsh
	sudo mv build/sawsh /usr/local/bin/sawsh

build_all:
	mkdir -p build
	go get -d
	for GOOS in $${GOOS_LIST:-darwin linux}; do \
		GOOS=$$GOOS GOARCH=amd64 go build -v -o build/sawsh-$$GOOS-amd64 sawsh.go; \
	done

build_docker:
	docker run --rm -v "$(PWD)":/usr/src/myapp -w /usr/src/myapp golang:1.8 make build_all

