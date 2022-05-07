FROM golang:1.18 AS build

ADD . /opt/app
WORKDIR /opt/app
RUN go build ./cmd/tinyUrl/main.go

FROM ubuntu:20.04

WORKDIR /usr/src/app

COPY . .
COPY --from=build /opt/app/main .

EXPOSE 5000

ENV DB=INMEM
ENTRYPOINT [ "./main" ]
CMD []