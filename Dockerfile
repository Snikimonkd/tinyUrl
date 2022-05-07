FROM golang:1.18 AS build

ADD . /opt/app
WORKDIR /opt/app
RUN go build ./cmd/tinyUrl/main.go

FROM ubuntu:20.04

ENV DEBIAN_FRONTEND noninteractive
ENV PGVER 12
RUN apt-get update -y && apt-get install -y postgresql postgresql-contrib

USER postgres
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    /etc/init.d/postgresql stop

EXPOSE 5432

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

WORKDIR /usr/src/app

COPY . .
COPY --from=build /opt/app/main .

EXPOSE 5000

ENV PGPASSWORD docker
RUN service postgresql start && psql -h localhost -d docker -U docker -p 5432 -a -q -f ./init.sql

ENV DB=NOSQL
CMD ["sh", "-c", "env $DB ./main"]