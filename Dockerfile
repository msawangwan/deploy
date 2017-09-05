FROM docker

WORKDIR /go/src/github.com/msawangwan/ci.io

COPY . .

VOLUME /var/run/docker.sock

EXPOSE 80

CMD ["/bin/sh"]

