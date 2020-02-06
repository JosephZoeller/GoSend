FROM golang:alpine AS builder

RUN mkdir /logmanager 
ADD . /logmanager
WORKDIR /logmanager

RUN apk add --no-cache git
RUN go get -d
RUN go build -o LogManager .

FROM alpine:latest

COPY --from=builder /logmanager .

CMD [ "./LogManager" ]