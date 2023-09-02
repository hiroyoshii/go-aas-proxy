FROM golang:1.21@sha256:b490ae1f0ece153648dd3c5d25be59a63f966b5f9e1311245c947de4506981aa AS build-env

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

WORKDIR /go/src/aas-proxy
Add go.mod .
Add go.sum .
RUN  go mod download

ADD . /go/src/aas-proxy

RUN go build -o /bin/aas-proxy cmd/main.go

FROM gcr.io/distroless/static:latest-amd64@sha256:bc535c40cfde8f8f1601f6cc9b51d3387db0722a7c4756896c68e3de4f074966 AS run-env

COPY --from=build-env /bin/aas-proxy /bin/aas-proxy

CMD [ "/bin/aas-proxy" ]