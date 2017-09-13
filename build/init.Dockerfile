FROM docker

WORKDIR /ci.io

COPY . .

EXPOSE 80

CMD . ./bin/build && /bin/sh
