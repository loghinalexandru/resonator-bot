#syntax=docker/dockerfile:latest

FROM golang:1.19-alpine AS build

WORKDIR /app

COPY ../ .

RUN go mod download
RUN go build -o /resonator

FROM golang:1.19-alpine

WORKDIR /

COPY --from=build /resonator /resonator
COPY --from=build /app/misc /misc

ENTRYPOINT ["/resonator"]