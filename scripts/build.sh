#!/usr/bin/env bash
#
# This script builds the application from source for multiple platforms.

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

cd "$DIR"

# Determine the arch/os combos we're building for
XC_ARCH=${XC_ARCH:-"386 amd64 arm"}
XC_OS=${XC_OS:-linux darwin windows}
XC_EXCLUDE_OSARCH="!darwin/arm !darwin/386"


# Delete the build dir
echo "==> Removing build directory..."
rm -rf build/{bin,pkg}
mkdir -p build/{bin,pkg}

if ! which gox > /dev/null; then
    echo "==> Installing gox..."
    (cd /tmp && go get -u github.com/mitchellh/gox)
fi

if ! which vembed > /dev/null; then
    echo "==> Installing vembed..."
    (cd /tmp && go get -u github.com/NoUseFreak/go-vembed/vembed)
fi

# Instruct gox to build statically linked binaries
export CGO_ENABLED=0

LD_FLAGS="-s -w "
LD_FLAGS+="`vembed`"

# Ensure all remote modules are downloaded and cached before build so that
# the concurrent builds launched by gox won't race to redundantly download them.
go mod download

# Build!
echo "==> Building..."
if [ "${DEV}" == "" ]; then
    gox \
        -os="${XC_OS}" \
        -arch="${XC_ARCH}" \
        -osarch="${XC_EXCLUDE_OSARCH}" \
        -ldflags "${LD_FLAGS}" \
        -output "build/bin/{{.OS}}_{{.Arch}}/${PWD##*/}" \
        .
else 
    go build --ldflags "${LD_FLAGS}" -o build/bin/  .
fi

# Done!
echo
echo "==> Results:"
du -ah build/*