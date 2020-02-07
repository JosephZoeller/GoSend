FROM golang:alpine AS builder

RUN mkdir /reverseproxy 
ADD . /reverseproxy
WORKDIR /reverseproxy

RUN apk add --no-cache git
RUN go get -d
RUN go build -o ReverseProxy .

FROM alpine:latest

COPY --from=builder /reverseproxy .
EXPOSE 8080
ENTRYPOINT [ "./ReverseProxy" ]