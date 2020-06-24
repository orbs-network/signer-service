const EthereumTx = require("ethereumjs-tx");
const { TransactionSigner } = require("web3-eth");
const NodeSignInputBuilder = require("./node-sign-input-builder");
const NodeSignOutputReader = require("./node-sign-output-reader");
const fetch = require("node-fetch");
const { ecrecover } = require("ethereumjs-util");
const { hexToBytes, getSignatureParameters } = require("web3-utils");
const { encode } = require("rlp");
const { keccak256 } = require("web3-utils");

class Signer {
    constructor(host) {
        this.host = host;
    }

    async _sign(payload) {
        const body = new NodeSignInputBuilder(payload).build();

        const res = await fetch(`${this.host}/ethSign`, {
            method: "post",
            body:  body,
            headers: { "Content-Type": "application/membuffers" },
        });

        if (!res.ok) {
            throw new Error(`Bad response: ${res.statusText}`);
        }

        const data = await res.buffer();
        return new NodeSignOutputReader(data).getSignature();
    }

    async sign(transaction, privateKey) {
        // we are going to ignore privateKey completely

        const ethTx = new EthereumTx(transaction);
        const signature = await this._sign(encode(ethTx.raw.slice(0, 6)));

        const { r, s, v } = getSignatureParameters("0x" + signature.toString("hex"));

        ethTx.r = r;
        ethTx.s = s;
        ethTx.v = v;

        const validationResult = ethTx.validate(true);

        if (validationResult !== '') {
            throw new Error(`TransactionSigner Error: ${validationResult}`);
        }

        const rlpEncoded = ethTx.serialize().toString('hex');
        const rawTransaction = '0x' + rlpEncoded;
        const transactionHash = keccak256(rawTransaction);

        return {
            messageHash: Buffer.from(ethTx.hash(false)).toString('hex'),
            v: '0x' + Buffer.from(ethTx.v).toString('hex'),
            r: '0x' + Buffer.from(ethTx.r).toString('hex'),
            s: '0x' + Buffer.from(ethTx.s).toString('hex'),
            rawTransaction,
            transactionHash
        };
    }
}

module.exports = Signer;