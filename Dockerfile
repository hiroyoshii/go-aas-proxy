FROM golang:1.21 AS build-env

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

ADD . /go/src/aas-proxy
WORKDIR /go/src/aas-proxy

RUN go build -o /bin/aas-proxy cmd/main.go

FROM gcr.io/distroless/static:latest-amd64 AS run-env

COPY --from=build-env /bin/aas-proxy /bin/aas-proxy

CMD [ "/bin/aas-proxy" ]