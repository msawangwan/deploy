FROM docker

WORKDIR /ci.io

COPY ./lib ./lib
COPY ./bin/build ./bin/

EXPOSE 80

CMD ./bin/build && /bin/sh
