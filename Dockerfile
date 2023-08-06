FROM golang:1.20 AS build-env

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

ADD . /go/src/aas-proxy
WORKDIR /go/src/aas-proxy

RUN go build -o /bin/aas-proxy cmd/main.go

FROM alpine:latest AS run-env

COPY --from=build-env /bin/aas-proxy /bin/aas-proxy

CMD [ "/bin/aas-proxy" ]