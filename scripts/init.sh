#!/bin/bash

# todo:
# get version number from a config script

cd ..

VERSION=$("cat conf/globals.json | grep version")

echo "working dir:"
pwd

echo "version:"
echo "$VERSION"

echo "build:"

docker system prune -f
docker build -t ci.io .
docker run -p 9001:80 -it --rm --name ci.io.running ci.io
