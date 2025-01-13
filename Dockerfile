FROM golang:1.15

WORKDIR /go/src/github.com/orbs-network/signer

ENV GO111MODULE=on

ADD go.* /go/src/github.com/orbs-network/signer/

RUN go mod download

ADD . .

RUN ./test.sh

RUN ./build-binaries.sh

FROM alpine:3.13

#RUN apk add --no-cache daemontools --repository http://dl-cdn.alpinelinux.org/alpine/edge/testing
RUN apk add --no-cache socat
RUN wget https://dl-cdn.alpinelinux.org/alpine/v3.0/testing/x86_64/daemontools-0.76-r1.apk
RUN apk add --allow-untrusted --no-cache daemontools-0.76-r1.apk || true

WORKDIR /opt/orbs

COPY --from=0 /go/src/github.com/orbs-network/signer/_bin/orbs-signer .

COPY --from=0 /go/src/github.com/orbs-network/signer/_bin/healthcheck .

ADD ./boyar/service /opt/orbs/service

VOLUME /opt/orbs/status
VOLUME /opt/orbs/logs

HEALTHCHECK CMD /opt/orbs/healthcheck --url http://localhost:7777 --output /opt/orbs/status/status.json --log /opt/orbs/logs/healthcheck

CMD /opt/orbs/service
