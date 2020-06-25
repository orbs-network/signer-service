const { Transaction } = require("ethereumjs-tx");
const fetch = require("node-fetch");
const { encode } = require("rlp");
const { keccak256, isHexStrict, hexToNumber } = require("web3-utils");
const NodeSignInputBuilder = require("./node-sign-input-builder");
const NodeSignOutputReader = require("./node-sign-output-reader");

function getSignatureParameters(signature) {
    if (!isHexStrict(signature)) {
        throw new Error(`Given value "${signature}" is not a valid hex string.`);
    }

    const r = signature.slice(0, 66);
    const s = `0x${signature.slice(66, 130)}`;
    let v = `0x${signature.slice(130, 132)}`;
    v = hexToNumber(v);

    if (![27, 28].includes(v)) v += 27;

    return {
        r,
        s,
        v
    };
}

class Signer {
    constructor(host) {
        this.host = host;
    }

    async _sign(payload) {
        const body = new NodeSignInputBuilder(payload).build();

        const res = await fetch(`${this.host}/eth-sign`, {
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

        const ethTx = new Transaction(transaction);
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