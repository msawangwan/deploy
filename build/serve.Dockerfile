FROM golang:1.9.0-alpine3.6

WORKDIR /go/src/github.com/msawangwan/ci.io

COPY . .

RUN apk add --no-cache git
RUN apk add --no-cache curl

RUN go-wrapper download
RUN go-wrapper install

EXPOSE 80

CMD ["go-wrapper", "run"]
