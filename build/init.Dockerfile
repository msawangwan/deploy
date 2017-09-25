FROM docker

WORKDIR /ci.io

COPY . .

EXPOSE 80

CMD . ./bin/listen && /bin/sh
