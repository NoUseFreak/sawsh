#!/usr/bin/env bash

# Usage: curl https://raw.githubusercontent.com/NoUseFreak/sawsh/master/scripts/get.sh | bash

PROJECT=sawsh

get_latest_release() {
  curl --silent "https://api.github.com/repos/NoUseFreak/$1/releases/latest" |
	grep '"tag_name":' |
	sed -E 's/.*"([^"]+)".*/\1/'
}

download() {
  rm -rf /usr/local/bin/$PROJECT
  curl -Ls -o /usr/local/bin/$PROJECT https://github.com/NoUseFreak/$PROJECT/releases/download/$1/`uname`_amd64.tar.gz
}

echo "Looking up latest release"
RELEASE=$(get_latest_release $PROJECT)

echo "Downloading package"
$(download $RELEASE)

echo "Making executable"
sudo chmod +x /usr/local/bin/$PROJECT

echo "Installed $PROJECT in /usr/local/bin/$PROJECT"
