# Signer Service

[![CI](https://circleci.com/gh/orbs-network/signer-service/tree/master.svg?style=svg)](https://circleci.com/gh/orbs-network/signer-service/tree/master)

Implements [Node Sign interface of the Vault service](https://github.com/orbs-network/orbs-spec/blob/master/behaviors/services/vault.md#nodesign)

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
