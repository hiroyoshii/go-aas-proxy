FROM golang:1.20 AS build-env

ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

ADD . /go/src/aas-proxy
WORKDIR /go/src/aas-proxy

RUN go build -o /bin/aas-proxy cmd/main.go

FROM gcr.io/distroless/static:latest AS run-env

COPY --from=build-env /bin/aas-proxy /bin/aas-proxy
COPY --from=build-env /go/src/aas-proxy/internal/aas/query_default.tpl.sql /app/query_default.tpl.sql
ENV AAS_QUERY_SQL_PATH=/app/query_default.tpl.sql

CMD [ "/bin/aas-proxy" ]