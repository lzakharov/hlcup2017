FROM golang:latest

RUN apt-get update && apt-get install -y postgresql-client 

WORKDIR /go/src/github.com/lzakharov/hlcup2017

COPY . .

RUN go get -v -d ./...
RUN go install -v ./...

RUN go get github.com/oxequa/realize

CMD ["realize", "start"]

EXPOSE 8000
