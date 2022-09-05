# syntax=docker/dockerfile:1
FROM golang:1.19 as build

WORKDIR /app

COPY ./ ./

ENV GOOS=linux
ENV GOARCH=amd64

RUN go mod download
RUN go build -o reverseproxy

FROM alpine as output

COPY --from=build app/reverseproxy /reverseproxy

EXPOSE 8080
CMD [ "/reverseproxy" ]
