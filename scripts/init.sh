#!/bin/bash

cd ..
sudo docker build -t ci.io .
sudo docker run -p 9001:80 -it --rm --name ci.io.running ci.io
