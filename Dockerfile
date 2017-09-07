FROM docker

WORKDIR /ci.io

COPY ./lib ./lib
COPY ./bin ./bin
#COPY ./bin/build ./bin/
#COPY ./bin/parseip ./bin/

EXPOSE 80

CMD . ./bin/build && /bin/sh
