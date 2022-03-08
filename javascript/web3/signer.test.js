const Signer = require("./signer");
const { Driver } = require("@orbs-network/orbs-ethereum-contracts-v2");
const BN = require('bn.js').BN;

const address = "0x29ce860a2247d97160d6dfc087a15f41e2349087";

it("should be able to sign messages for EIP155", async () => {

    const tx = {
        from: address,
        gasPrice: "0x4A817C800",
        gasLimit: "0x21072",
        to: "0x3535353535353535353535353535353535353535",
        value: "0xDE0B6B3A7640000",
        data: "0x001"
    };

    const transactionSigner = new Signer("http://localhost:7777");
    let signedTransaction = await transactionSigner.sign(tx, 1, address);
    expect(signedTransaction).toEqual({
        messageHash: "2fdf20b217df1f790514c2b0f5668a4628cd97a8ddb10a2983108e633dbc439c",
        v: "0x26",
        r: "0xde18f8efaf1a6b0553a7ffdd26b95216230f957cde31851dd11583d17ce27f50",
        s: "0x7d521ea616ab688f13cd828b4205e74615afe8b07fd010e6973d0dabf6b756b5",
        rawTransaction: "0xf86f808504a817c80083021072943535353535353535353535353535353535353535880de0b6b3a764000082000126a0de18f8efaf1a6b0553a7ffdd26b95216230f957cde31851dd11583d17ce27f50a07d521ea616ab688f13cd828b4205e74615afe8b07fd010e6973d0dabf6b756b5",
        transactionHash: "0xea008472307cc3ff45377e423c46cc0c2b9d3bd25134e78faa4fd891463310e1"
    });

	signedTransaction = await transactionSigner.sign(tx, 137, address);
    expect(signedTransaction).toEqual({
        messageHash: "52cc881fb05b70d9e56a03a306befa688a03732085b052c563bd450a85d1a42b",
        v: "0x135",
        r: "0xaf3933362c055de089117bca0bc28673e2dd7b1897b1a6bd99360f5e8cb5bf42",
        s: "0x8d579dcee01c70a161336625ff4321b2c5159999bcd064e25f130ab4a69c8ff",
        rawTransaction: "0xf871808504a817c80083021072943535353535353535353535353535353535353535880de0b6b3a7640000820001820135a0af3933362c055de089117bca0bc28673e2dd7b1897b1a6bd99360f5e8cb5bf42a008d579dcee01c70a161336625ff4321b2c5159999bcd064e25f130ab4a69c8ff",
        transactionHash: "0xa1e59bf1cf32e3fba6db8977ae6d63a7504f77d6bd3724c14001d332216361e0"
    });
});

it("failed address validation when signing with wrong address", async () => {
	const otherAddr = "0x9f0988Cd37f14dfe95d44cf21f9987526d6147Ba";

    const tx = {
        from: address,
        gasPrice: "0x4A817C800",
        gasLimit: "0x21072",
        to: "0x3535353535353535353535353535353535353535",
        value: "0xDE0B6B3A7640000",
        data: "0x001"
    };

    const transactionSigner = new Signer("http://localhost:7777");

    await expect(async () => {await transactionSigner.sign(tx, 1, otherAddr)})
    .rejects.toThrowError(`Sender address mismatch after signing: expected ${otherAddr}, got ${address}`)

});


it("should work with ganache for EIP155", async () => {
    const d = await Driver.new();

    const assignTx = {
        from: address,
        to: d.erc20.address,
        gasLimit: 0x7fffffff,
        data: d.erc20.web3Contract.methods.assign(address, new BN(200)).encodeABI(),
        nonce: await d.web3.eth.getTransactionCount(address),
    };

    const transactionSigner = new Signer("http://localhost:7777");
    const { rawTransaction, transactionHash } = await transactionSigner.sign(assignTx, 137, address);

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

it("should be able to sign messages", async () => {
    const tx = {
        from: address,
        gasPrice: "0x4A817C800",
        gasLimit: "0x21072",
        to: "0x3535353535353535353535353535353535353535",
        value: "0xDE0B6B3A7640000",
        data: "0x001",
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
