const { Transaction } = require("@ethereumjs/tx");
const Common = require("@ethereumjs/common");
const fetch = require("node-fetch");
const { encode } = require("rlp");
const { keccak256, isHexStrict, hexToNumber } = require("web3-utils");
const NodeSignInputBuilder = require("./node-sign-input-builder");
const NodeSignOutputReader = require("./node-sign-output-reader");

function getSignatureParameters(signature, chainId) {
    if (!isHexStrict(signature)) {
        throw new Error(`Given value "${signature}" is not a valid hex string.`);
    }

    const r = signature.slice(0, 66);
    const s = `0x${signature.slice(66, 130)}`;
    let v = `0x${signature.slice(130, 132)}`;
    v = hexToNumber(v);

    if (![27, 28].includes(v)) v += 27;

    if (chainId) {
        v = chainId * 2 + 8 + v;// see https://github.com/ethereum/EIPs/blob/master/EIPS/eip-155.md
    }

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

    async signEip155(transaction, chainId, expectedSenderAddress) {
        chainId = chainId || 1;

        const common = Common.default.custom({ chainId: chainId })
        const ethTx = new Transaction(transaction, { common });
        const payload = encode(ethTx.getMessageToSign(false));
        const signature = await this._sign(payload);

        const { r, s, v } = getSignatureParameters("0x" + signature.toString("hex"), chainId);
        const signedTxData = {...transaction, v,r,s}
        const signedTx = Transaction.fromTxData(signedTxData, { common })
        const from = signedTx.getSenderAddress().toString()

        console.log(`signedTx: 0x${signedTx.serialize().toString('hex')}\nfrom: ${from}`)

        if (expectedSenderAddress && from !== expectedSenderAddress) { // optional - pass undefined to disable check
            throw new Error(`Sender address mismatch after signing: expected ${expectedSenderAddress}, got ${from}`);
        }

        const validationResult = signedTx.validate(true);

        if (Array.isArray(validationResult) && validationResult.length > 0) {
            throw new Error(`TransactionSigner Error: ${validationResult}`);
        }

        const rlpEncoded = signedTx.serialize().toString('hex');
        const rawTransaction = '0x' + rlpEncoded;
        const transactionHash = keccak256(rawTransaction);

        return {
            messageHash: Buffer.from(signedTx.getMessageToSign(true)).toString('hex'),
            v: '0x' + signedTx.v.toString(16),
            r: '0x' + signedTx.r.toString(16),
            s: '0x' + signedTx.s.toString(16),
            rawTransaction,
            transactionHash
        };
    }

	async sign(transaction, chainId, expectedSenderAddress) {

		if (chainId) {
			return this.signEip155(transaction, chainId, expectedSenderAddress)
		}
		else {
			return this.signLegacy(transaction)
		}
	}

    async signLegacy(transaction, privateKey) {
        // we are going to ignore privateKey completely

        const ethTx = new Transaction(transaction);
        const signature = await this._sign(encode(ethTx.raw().slice(0, 6)));

        const { r, s, v } = getSignatureParameters("0x" + signature.toString("hex"));
        const signedTxData = {...transaction, v,r,s}
        const signedTx = Transaction.fromTxData(signedTxData)

        const validationResult = signedTx.validate(true);

        if (Array.isArray(validationResult) && validationResult.length > 0) {
            throw new Error(`TransactionSigner Error: ${validationResult}`);
        }

        const rlpEncoded = signedTx.serialize().toString('hex');
        const rawTransaction = '0x' + rlpEncoded;
        const transactionHash = keccak256(rawTransaction);

        return {
            messageHash: Buffer.from(signedTx.getMessageToSign(true)).toString('hex'),
            v: '0x' + signedTx.v.toString(16),
            r: '0x' + signedTx.r.toString(16),
            s: '0x' + signedTx.s.toString(16),
            rawTransaction,
            transactionHash
        };
    }
}

module.exports = Signer;
