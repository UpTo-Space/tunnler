FROM golang:1.23.2-alpine

RUN mkdir /app
WORKDIR /app
ADD . /app

RUN go install -mod=mod github.com/githubnemo/CompileDaemon

ENTRYPOINT   --build="go build ." --command=./server
