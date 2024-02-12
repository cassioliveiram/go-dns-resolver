FROM golang:1.20.0-alpine as build

WORKDIR /dns-lookup
COPY * .
RUN go build .

FROM alpine:3.19 as final

LABEL maintainer="Cassio Moreira (cassioliveiram@gmail.com)"
WORKDIR /app
COPY --from=build /dns-lookup/dns-lookup .
CMD ["./dns-lookup"]