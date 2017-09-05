FROM docker

WORKDIR /ci.io

COPY ./lib ./lib
COPY ./bin/build ./bin/

CMD ./bin/build && /bin/sh
