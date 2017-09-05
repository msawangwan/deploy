FROM docker

WORKDIR /src

COPY . .

#VOLUME /var/run/docker.sock

EXPOSE 80

CMD ./bin/build && /bin/sh
