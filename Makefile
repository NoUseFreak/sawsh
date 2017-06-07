build:
	go build awsh.go

install: build
	chmod +x awsh
	sudo mv awsh /usr/local/bin/awsh
