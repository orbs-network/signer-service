// const SignerProvider = require("./provider");
const Web3 = require("Web3");
const Signer = require("./signer");

it("should be able to sign messages", async () => {
    const transactionSigner = new Signer("http://localhost:7777");
    const web3 = new Web3("http://localhost:8545");

    const tx = {
        from: "0xEB014f8c8B418Db6b45774c326A0E64C78914dC0",
        gasPrice: "20000000000",
        gas: "21000",
        to: '0x3535353535353535353535353535353535353535',
        value: "1000000000000000000",
        data: "0x001"
    };

    const signedTransaction = await transactionSigner.sign(tx);
    console.log(signedTransaction);
});