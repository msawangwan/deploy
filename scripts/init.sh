#!/bin/sh

# todo:
# get version number from a config script

cd ..

echo "root dir:"
pwd

VERSION="$(/bin/cat conf/globals.json | /bin/grep version)"

echo "version:"
echo "$VERSION"

echo "build:"

docker system prune -f
docker build -t ci.io .
docker run -p 9001:80 -it --rm --name ci.io.running ci.io
