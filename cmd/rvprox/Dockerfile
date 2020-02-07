FROM golang:alpine AS builder

RUN mkdir /rvprox 
ADD . /rvprox
WORKDIR /rvprox

RUN apk add --no-cache git
RUN go get -d
RUN go build -o ReverseProxy .

FROM alpine:latest

COPY --from=builder /rvprox .
EXPOSE 8080
CMD [ "./ReverseProxy" ]