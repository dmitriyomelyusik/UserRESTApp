FROM golang

COPY ./bin /go/restapp/

ENV DBNAME=postgres PGPASS=password PGUSER=postgres PGHOST=postgres SSLMODE=disable

CMD restapp/main

EXPOSE 8080