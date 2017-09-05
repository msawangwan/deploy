FROM docker

WORKDIR /ci.io/bin
COPY ./bin/build .
CMD ./build && /bin/sh

#VOLUME /var/run/docker.sock

# EXPOSE 80

