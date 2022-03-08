FROM golang:1.15

WORKDIR /go/src/github.com/orbs-network/signer

ENV GO111MODULE=on

ADD go.* /go/src/github.com/orbs-network/signer/

RUN go mod download

ADD . .

RUN ./test.sh

RUN ./build-binaries.sh

FROM alpine:3.12

RUN apk add --no-cache daemontools --repository http://dl-cdn.alpinelinux.org/alpine/edge/testing

WORKDIR /opt/orbs

COPY --from=0 /go/src/github.com/orbs-network/signer/_bin/orbs-signer .

COPY --from=0 /go/src/github.com/orbs-network/signer/_bin/healthcheck .

ADD ./boyar/service /opt/orbs/service

VOLUME /opt/orbs/status
VOLUME /opt/orbs/logs

HEALTHCHECK CMD /opt/orbs/healthcheck --url http://localhost:7777 --output /opt/orbs/status/status.json --log /opt/orbs/logs/healthcheck

CMD /opt/orbs/service
