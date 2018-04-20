appname := sawsh
version := $(shell git describe --tags --abbrev=0)

sources := $(wildcard *.go)

build = GOOS=$(1) GOARCH=$(2) go build -ldflags=-s -o build/$(appname)$(3)
tar = cd build && tar -cvzf $(1)_$(2).tar.gz $(appname)$(3) && rm $(appname)$(3)
zip = cd build && zip $(1)_$(2).zip $(appname)$(3) && rm $(appname)$(3)
fpm_deb = cd build && rm -rf sawsh_${version}_${2}.deb \
	  && docker run --rm -v `pwd`:/workspace -w /workspace isim/fpm -s tar -t deb -a $(2) -n ${appname} \
	  -v ${version} --deb-no-default-config-files --prefix=/usr/local/bin -p sawsh_${version}_${2}.deb $(1)_$(2).tar.gz
fpm_rpm = cd build && tar xvf $(1)_$(2).tar.gz && rm -rf sawsh_${version}_${2}.rpm \
	  && docker run --rm -v `pwd`:/workspace -w /workspace liuedy/centos-fpm fpm -s dir -t rpm -a $(2) -n ${appname} \
	  -v ${version} --prefix=/usr/local/bin -p sawsh_${version}_${2}.rpm ${appname} \
	  && rm -rf ${appname}
fpm_osx = cd build && tar xvf $(1)_$(2).tar.gz && rm -rf sawsh_${version}_${2}.pkg \
	  && fpm -s dir -t osxpkg -a $(2) -n ${appname} \
	  -v ${version} --prefix=/usr/local/bin -p sawsh_${version}_${2}.pkg ${appname} \
	  && rm -rf ${appname}

.PHONY: all windows darwin linux clean

default:
	mkdir -p build
	go get -d
	go build -o build/sawsh sawsh.go

all: windows darwin linux

clean:
	rm -rf build/

##### LINUX BUILDS #####
linux: build/linux_arm.tar.gz build/linux_arm64.tar.gz build/linux_386.tar.gz build/linux_amd64.tar.gz

build/linux_386.tar.gz: $(sources)
	$(call build,linux,386,)
	$(call tar,linux,386)

build/linux_amd64.tar.gz: $(sources)
	$(call build,linux,amd64,)
	$(call tar,linux,amd64)

build/linux_arm.tar.gz: $(sources)
	$(call build,linux,arm,)
	$(call tar,linux,arm)

build/linux_arm64.tar.gz: $(sources)
	$(call build,linux,arm64,)
	$(call tar,linux,arm64)

##### DARWIN (MAC) BUILDS #####
darwin: build/darwin_amd64.tar.gz

build/darwin_amd64.tar.gz: $(sources)
	$(call build,darwin,amd64,)
	$(call tar,darwin,amd64)

##### WINDOWS BUILDS #####
windows: build/windows_386.zip build/windows_amd64.zip

build/windows_386.zip: $(sources)
	$(call build,windows,386,.exe)
	$(call zip,windows,386,.exe)

build/windows_amd64.zip: $(sources)
	$(call build,windows,amd64,.exe)
	$(call zip,windows,amd64,.exe)

install:
	chmod +x build/sawsh
	sudo mv build/sawsh /usr/local/bin/sawsh

build_docker:
	docker run --rm -v "$(PWD)":/usr/src/myapp -w /usr/src/myapp golang:1.8 make build_all

##### Packaging #####
package: all
	$(call fpm_osx,darwin,amd64)
	$(call fpm_rpm,linux,386)
	$(call fpm_deb,linux,386)
	$(call fpm_rpm,linux,amd64)
	$(call fpm_deb,linux,amd64)
	$(call fpm_rpm,linux,arm)
	$(call fpm_deb,linux,arm)
	$(call fpm_rpm,linux,arm64)
	$(call fpm_deb,linux,arm64)

