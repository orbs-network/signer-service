const Signer = require("./signer");
const { Driver } = require("@orbs-network/orbs-ethereum-contracts-v2");
const BN = require('bn.js').BN;

const address = "0x29ce860a2247d97160d6dfc087a15f41e2349087";

it("should be able to sign messages", async () => {
    const tx = {
        from: address,
        gasPrice: "0x4A817C800",
        gasLimit: "0x21072",
        to: "0x3535353535353535353535353535353535353535",
        value: "0xDE0B6B3A7640000",
        data: "0x001"
    };

    const transactionSigner = new Signer("http://localhost:7777");
    const signedTransaction = await transactionSigner.sign(tx);
    expect(signedTransaction).toEqual({
        messageHash: "0b412d2c3ee990ea6bcf44a78a601a8cc2245ab5516bbb2f0a9a1199a2b36f85",
        v: "0x1b",
        r: "0xb2de7a48e95b56570efbea466d732f381d5bfe9f6aeac08a5eac8f4d01b59136",
        s: "0x4aca32813e60365d4ba17bad4afdf28b7f52c4e0eaee151b7e5719d9d1a1ee37",
        rawTransaction: "0xf86f808504a817c80083021072943535353535353535353535353535353535353535880de0b6b3a76400008200011ba0b2de7a48e95b56570efbea466d732f381d5bfe9f6aeac08a5eac8f4d01b59136a04aca32813e60365d4ba17bad4afdf28b7f52c4e0eaee151b7e5719d9d1a1ee37",
        transactionHash: "0xc5be0d5a6c01c869b9e9713f7d1bf03802070c71f6e8fcdbca300a8b21cf08af"
    });
});

it("should work with ganache", async () => {
    const d = await Driver.new();

    const assignTx = {
        from: address,
        to: d.erc20.address,
        gasLimit: 0x7fffffff,
        data: d.erc20.web3Contract.methods.assign(address, new BN(200)).encodeABI(),
        nonce: await d.web3.eth.getTransactionCount(address),
    };

    const transactionSigner = new Signer("http://localhost:7777");
    const { rawTransaction, transactionHash } = await transactionSigner.sign(assignTx);

    await d.web3.eth.sendSignedTransaction(rawTransaction);

    while (true) {
        const status = await d.web3.eth.getTransactionReceipt(transactionHash);
        if (status.status) {
            console.log(status);
            break;
        } else {
            console.log(logs);
            throw new Error(logs[0]);
        }
    }

    const balance = await d.erc20.balanceOf(address);
    expect(balance).toEqual(new BN(200).toString());
});