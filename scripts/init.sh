#!/bin/bash

cd ..
sudo docker build -t ci.io
sudo docker run -it --rm --name ci.io.running ci.io