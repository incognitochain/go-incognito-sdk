FROM golang:latest

RUN go get golang.org/x/tools/cmd/godoc

RUN mkdir -p $GOROOT/src/source

WORKDIR $GOROOT/src/source

COPY . .

EXPOSE 6060
CMD godoc -http=:6060