#syntax=docker/dockerfile:latest

FROM golang:1.19-alpine AS build

WORKDIR /app

COPY ../ .

RUN go mod download
RUN go build -o /resonator

FROM golang:1.19-alpine

WORKDIR /

COPY --link --from=build /resonator /resonator
COPY --link --from=build /app/misc /misc

ENTRYPOINT ["/resonator"]