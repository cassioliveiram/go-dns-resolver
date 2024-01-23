FROM golang:1.20.0-alpine

LABEL maintainer="Cassio Moreira (cassioliveiram@gmail.com)"

WORKDIR /dns-lookup

COPY * .

RUN go build .

CMD ["./dns-lookup"]