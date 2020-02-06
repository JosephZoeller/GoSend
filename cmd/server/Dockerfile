FROM golang:alpine AS builder

RUN mkdir /fileserver 
ADD . /fileserver
WORKDIR /fileserver

RUN apk add --no-cache git
RUN go get -d
RUN go build -o FileServer .

FROM alpine:latest

COPY --from=builder /fileserver .

CMD [ "./FileServer" ]