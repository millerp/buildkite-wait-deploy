FROM golang:1.11

WORKDIR /go/src/deploy-server

COPY main.go main.go

RUN go get -d -v .
RUN go install -v .

CMD deploy-server

EXPOSE 8000