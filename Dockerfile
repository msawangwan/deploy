FROM docker

WORKDIR /ci.io

COPY ./lib ./lib
COPY ./bin/build ./bin/

#WORKDIR /ci.io/src
#COPY ./lib .

#WORKDIR /ci.io/bin
#COPY ./bin/build .

#WORKDIR /ci.io

CMD ./bin/build && /bin/sh

#VOLUME /var/run/docker.sock

# EXPOSE 80

