# Signer Service

[![CI](https://circleci.com/gh/orbs-network/signer-service/tree/master.svg?style=svg)](https://circleci.com/gh/orbs-network/signer-service/tree/master)

The signer service manages the node private keys (Orbs and Ethereum) securely and signs transactions on behalf of the node. Used by ONG to sign blocks and protocol messages. Also used by Ethereum Writer to sign Ethereum transactions.

The service implements [Node Sign interface of the Vault service](https://github.com/orbs-network/orbs-spec/blob/master/vchain-architecture/services/vault.md#nodesign) part of the Orbs protocol specification.

DockerHub: [https://hub.docker.com/repository/docker/orbsnetwork/signer](https://hub.docker.com/repository/docker/orbsnetwork/signer)

## Testing

`./test.sh`

## Building

`./create-docker-version.sh && ./docker-build.sh` will produce new docker image.

## CLI reference

```
Usage:
  -config value
    	path/to/config.json
  -listen string
    	ip address and port for http server (default ":7777")
  -version
    	returns information about version
```
