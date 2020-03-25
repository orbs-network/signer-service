FROM golang:1.12.9

WORKDIR /go/src/github.com/orbs-network/signer

ENV GO111MODULE=on

ADD go.* /go/src/github.com/orbs-network/signer/

RUN go get

ADD . .

RUN ./test.sh

RUN ./build-binaries.sh

FROM busybox

WORKDIR /opt/orbs

COPY --from=0 /go/src/github.com/orbs-network/signer/_bin/orbs-signer .

COPY --from=0 /go/src/github.com/orbs-network/signer/_bin/healthcheck .

VOLUME /opt/orbs/status

HEALTHCHECK CMD /opt/orbs/healthcheck --url http://localhost:7777 --output /opt/orbs/status/status.json

CMD /opt/orbs/orbs-signer