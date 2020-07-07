FROM golang:1.12.9

WORKDIR /go/src/github.com/orbs-network/signer

ENV GO111MODULE=on

ADD go.* /go/src/github.com/orbs-network/signer/

RUN go get

ADD . .

RUN ./test.sh

RUN ./build-binaries.sh

FROM alpine:3.12

WORKDIR /opt/orbs

RUN apk add --no-cache daemontools --repository http://dl-cdn.alpinelinux.org/alpine/edge/testing

COPY --from=0 /go/src/github.com/orbs-network/signer/_bin/orbs-signer .

COPY --from=0 /go/src/github.com/orbs-network/signer/_bin/healthcheck .

VOLUME /opt/orbs/status

COPY ./entrypoint.sh /opt/orbs/service

HEALTHCHECK CMD /opt/orbs/healthcheck --url http://localhost:7777 --output /opt/orbs/status/status.json

CMD /opt/orbs/service