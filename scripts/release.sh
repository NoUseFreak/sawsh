#!/usr/bin/env bash

set -e 

# Get the parent directory of where this script is.
SOURCE="${BASH_SOURCE[0]}"
while [ -h "$SOURCE" ] ; do SOURCE="$(readlink "$SOURCE")"; done
DIR="$( cd -P "$( dirname "$SOURCE" )/.." && pwd )"

cd "$DIR"

mkdir /tmp/bin
curl -sL http://bit.ly/gh-get | BIN_DIR=/tmp/bin PROJECT=NoUseFreak/letitgo bash
/tmp/bin/letitgo $(git describe --tags --abbrev=0)
