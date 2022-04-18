#!/usr/bin/env bash
set -e

# try and use the correct MD5 lib (depending on user OS darwin/linux)
MD5=$(which md5 || which md5sum)

# for versioning
getCurrCommit() {
  echo `git rev-parse --short HEAD| tr -d "[ \r\n\']"`
}

# for versioning
getCurrTag() {
  echo `git describe --always --tags --abbrev=0 | tr -d "[v\r\n]"`
}

# remove any previous builds that may have failed
[ -e "./build" ] && \
  echo "Cleaning up old builds..." && \
  rm -rf "./build"

# build slurp
echo "Building slurp..."
gox -ldflags="-X main.version=$(getCurrTag)
  -X main.commit=$(getCurrCommit)" \
  -osarch="linux/$(go env | grep GOARCH | sed -E 's/GOARCH="(.*)"/\1/')" \
  -output="./build/slurp"

# look through each os/arch/file and generate an md5 for each
echo "Generating md5s..."
for file in $(ls ./build); do
  cat "./build/${file}" | ${MD5} | awk '{print $1}' >> "./build/${file}.md5"
done
