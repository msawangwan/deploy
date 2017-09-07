FROM docker

WORKDIR /ci.io

#COPY ./lib ./lib
#COPY ./bin ./bin
COPY . .

EXPOSE 80

CMD . ./bin/build && /bin/sh
