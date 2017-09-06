FROM golang:1.9.0-alpine3.6

WORKDIR /go/src/github.com/msawangwan/ci.io/lib

COPY . .

RUN apk add --no-cache git

RUN go-wrapper download
RUN go-wrapper install

#RUN apk del git

EXPOSE 80

CMD ["go-wrapper", "run"]
