const EthereumTx = require("ethereumjs-tx");
const { TransactionSigner } = require("web3-eth");
const NodeSignInputBuilder = require("./node-sign-input-builder");
const NodeSignOutputReader = require("./node-sign-output-reader");
const fetch = require("node-fetch");
const { ecrecover, bufferToInt } = require("ethereumjs-util");

class Signer {
    constructor(host) {
        this.host = host;
    }

    async _sign(payload) {
        const body = new NodeSignInputBuilder(payload).build();

        const req = await fetch(`${this.host}/sign`, {
            method: "post",
            body:  body,
            headers: { "Content-Type": "application/membuffers" },
        });
        const res = await req.buffer();
        console.log(res);
        const signature = new NodeSignOutputReader(res).getSignature();
        console.log(signature);

        return signature;
    }

    async sign(transaction, privateKey) {
        // we are going to ignore privateKey completely

        const ethTx = new EthereumTx(transaction);

        const hash = ethTx.hash(false);
        console.log("payload", hash.toString("hex"))
        const signature = await this._sign(hash);

        // this signature does not conform to the acceptable format

        console.log("sig", signature.toString("hex"));

        ethTx.sig = signature;
        // should split signature according in 3 parts (r, s, v) to pass validation

        const validationResult = ethTx.validate(true);

        if (validationResult !== '') {
            throw new Error(`TransactionSigner Error: ${validationResult}`);
        }

        const rlpEncoded = ethTx.serialize().toString('hex');
        const rawTransaction = '0x' + rlpEncoded;
        const transactionHash = this.utils.keccak256(rawTransaction);

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