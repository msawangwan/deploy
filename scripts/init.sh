#!/bin/bash

cd ..

docker build -t ci.io .
docker run -p 9001:80 -it --rm --name ci.io.running ci.io

pwd