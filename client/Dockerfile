FROM golang:1.11-alpine

WORKDIR /go/src/app

RUN apk add --update git

COPY docker/ipe/ .

RUN go get -d -v .
RUN go install -v .

CMD app -logtostderr=true -v=2

EXPOSE 36511