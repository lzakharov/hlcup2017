FROM golang:latest

WORKDIR /go/src/github.com/lzakharov/hlcup2017

COPY . .

RUN go get -v -d ./...
RUN go install -v ./...

RUN go get github.com/oxequa/realize

ENTRYPOINT ["realize", "start"]

EXPOSE 8000
