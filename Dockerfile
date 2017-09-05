FROM docker

# WORKDIR /go/src/github.com/msawangwan/ci.io
WORKDIR /src

COPY . .

VOLUME /var/run/docker.sock

EXPOSE 80

CMD ["/bin/sh", "./bin/build"]