FROM golang:1.13.5-alpine3.10 as builder
RUN apk update && apk upgrade && \
    apk add \
    xz-dev \
    musl-dev \
    gcc
RUN mkdir -p /go/src/github.com/canyanio/rating-agent-hep
COPY . /go/src/github.com/canyanio/rating-agent-hep
RUN cd /go/src/github.com/canyanio/rating-agent-hep && env CGO_ENABLED=1 go build

FROM alpine:3.10
RUN apk update && apk upgrade && \
        apk add --no-cache ca-certificates xz
RUN mkdir -p /etc/rating-agent-hep
COPY ./config.yaml /etc/rating-agent-hep
COPY --from=builder /go/src/github.com/canyanio/rating-agent-hep/rating-agent-hep /usr/bin
ENTRYPOINT ["/usr/bin/rating-agent-hep", "--config", "/etc/rating-agent-hep/config.yaml"]

EXPOSE 9060/udp
