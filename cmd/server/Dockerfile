FROM golang:alpine AS builder

RUN mkdir /srvr 
ADD . /srvr
WORKDIR /srvr

RUN apk add --no-cache git
RUN go get -d
RUN go build -o FileServer .

FROM alpine:latest
COPY --from=builder /srvr .

CMD [ "./FileServer" ]