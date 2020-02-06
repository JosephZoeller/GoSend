FROM golang:alpine AS builder

RUN mkdir /fileclient 
ADD . /fileclient
WORKDIR /fileclient

RUN apk add --no-cache git
RUN go get -d
RUN go build -o FileClient .

FROM alpine:latest

COPY --from=builder /fileclient .

CMD [ "./FileClient" ]