# syntax=docker/dockerfile:1
FROM golang:1.19 as build

WORKDIR /app

COPY ./ ./

RUN go mod download
RUN go build -o reverseproxy
RUN go install

FROM alpine as output

EXPOSE 8080
CMD [ "./reverseproxy" ]
