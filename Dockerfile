#syntax=docker/dockerfile:latest

FROM golang:1.21-alpine AS build

WORKDIR /app

COPY ../ .

RUN go mod download
RUN go build -o / .

FROM golang:1.21-alpine

WORKDIR /

COPY --link --from=build /resonator /resonator
COPY --link --from=build /app/misc /misc

ENTRYPOINT ["/resonator"]