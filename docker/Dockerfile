FROM golang:alpine as builder
WORKDIR /build
ADD ./Dockerfile /build

RUN apk add git
RUN git clone https://github.com/shoenig/donutdns

WORKDIR /build/donutdns
RUN go version && \
    go env && \
    CGO_ENABLED=0 GOOS=linux go build

FROM alpine:3
MAINTAINER gophers.dev

WORKDIR /opt
COPY --from=builder /build/donutdns /opt

ENTRYPOINT ["/opt/donutdns"]
