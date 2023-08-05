FROM golang:1.20 AS build-env

ADD . /go/src/aas-proxy
WORKDIR /go/src/aas-proxy

RUN go build -o /bin/aas-proxy cmd/main.go

FROM gcr.io/distroless/static:latest AS run-env

COPY --from=build-env /bin/aas-proxy /bin/

CMD [ "/bin/aas-proxy" ]
