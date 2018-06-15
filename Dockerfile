FROM golang

COPY ./bin /go/restapp/

ENV DBNAME=postgres PGPASS=password PGUSER=postgres PGHOST=127.0.0.1 SSLMODE=disable

CMD restapp/main

EXPOSE 8080