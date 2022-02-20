FROM golang:alpine as builder
WORKDIR /build
ADD . /build
RUN go version && \
    go env && \
    CGO_ENABLED=0 GOOS=linux go build

FROM alpine:3
MAINTAINER gophers.dev

WORKDIR /opt
COPY --from=builder /build/donutdns /opt

ENTRYPOINT ["/opt/donutdns"]
