FROM golang:1.9.0-alpine3.6

WORKDIR /go/src/github.com/msawangwan/ci.io

COPY . .

RUN go-wrapper download
RUN go-wrapper install

#RUN go install .

EXPOSE 80

CMD ["go-wrapper", "run"]
#ENTRYPOINT /go/src/github.com/msawangwan/ci.io
