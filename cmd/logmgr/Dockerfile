FROM golang:alpine AS builder

RUN mkdir /log 
ADD . /log
WORKDIR /log

RUN apk add --no-cache git
RUN go get -d
RUN go build -o LogManager .

FROM alpine:latest

COPY --from=builder /log .
CMD [ "./LogManager" ]