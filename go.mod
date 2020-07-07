module github.com/orbs-network/signer-service

go 1.13

require (
	github.com/google/go-cmp v0.4.0
	github.com/orbs-network/crypto-lib-go v1.2.0
	github.com/orbs-network/govnr v0.2.0
	github.com/orbs-network/healthcheck v1.1.0
	github.com/orbs-network/orbs-spec v0.0.0-20200624091201-eb5ca526fb87
	github.com/orbs-network/scribe v0.2.3
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.5.1
	golang.org/x/crypto v0.0.0-20200311171314-f7b00557c8c4
)

replace crypto-lib-go => ../crypto-lib-go
