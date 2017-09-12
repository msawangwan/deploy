FROM golang:1.9.0-alpine3.6

WORKDIR /go/src/github.com/msawangwan/ci.io

COPY . .

RUN echo 'PS1="\[$(tput setaf 3)$(tput bold)[\]appname@\\h$:\\w]#\[$(tput sgr0) \]"' >> /root/.bashrc

RUN apk add --no-cache git
RUN apk add --no-cache curl

RUN go-wrapper download
RUN go-wrapper install

#RUN apk del git

EXPOSE 80

CMD ["go-wrapper", "run"]
