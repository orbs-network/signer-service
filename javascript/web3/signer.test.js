// const Web3 = require("Web3");
const Signer = require("./signer");
const {spawn} = require("child_process");
const { ganacheDriver, Driver } = require("@orbs-network/orbs-ethereum-contracts-v2");
const BN = require('bn.js').BN;

let signerProcess;

beforeAll(async () => {
    signerProcess = spawn("go", [
        "run", "./bootstrap/signer/main/main.go", "-config", "./javascript/web3/keys.json"
    ], {
        cwd: `${__dirname}/../../`
    });
    signerProcess.stdout.on("data", (data) => {
        console.log(`stdout: ${data}`);
    });

    await ganacheDriver.startGanache();

    await new Promise((resolve, reject) => {
        setTimeout(resolve, 3000);
    });
});

afterAll(async () => {
    await ganacheDriver.stopGanache();

    const ok = signerProcess.kill();
    console.log(`terminating child process: ${ok}`);
});

const address = "0xa328846cd5b4979d68a8c58a9bdfeee657b34de7";

it("should be able to sign messages", async () => {
    const transactionSigner = new Signer("http://localhost:7777");
    // const web3 = new Web3("http://localhost:8545");

    const tx = {
        from: address,
        gasPrice: "20000000000",
        gas: "21000",
        to: "0x3535353535353535353535353535353535353535",
        value: "1000000000000000000",
        data: "0x001"
    };

    const signedTransaction = await transactionSigner.sign(tx);
    expect(signedTransaction).toEqual({
        messageHash: "6f1c0a083b53c943b3f74ebd2bd0be8d5ee02cd6d36bf080f274feaef14bd2f1",
        v: "0x1c",
        r: "0x50b69b24790fbdf91bd0272fef54f7490fb4f61cb07a91a3d61e6c115a6fe80b",
        s: "0x76df08f4f3a5763bc721423c89c074fec9af0ed86bf889973a85499c4691cbf2",
        rawTransaction: "0xf882808b323030303030303030303085323130303094353535353535353535353535353535353535353593313030303030303030303030303030303030308200011ca050b69b24790fbdf91bd0272fef54f7490fb4f61cb07a91a3d61e6c115a6fe80ba076df08f4f3a5763bc721423c89c074fec9af0ed86bf889973a85499c4691cbf2",
        transactionHash: "0xe317b4d51a7c9ed3aec765ad32c021af031630e1e60607d49da08ac6c9a17848"
    });
});

it("should work with ganache", async () => {
    const d = await Driver.new();
    await d.erc20.assign(address, BN(500));

    const assignTx = {
        from: address,
        to: d.erc20.address,
        gas: 0x7fffffff,
        data: d.erc20.methods.assign(address, BN(200)).encodeABI(),
    };

    const { rawTransaction } = await transactionSigner.sign(assignTx);

    await d.web3.eth.sendSignedTransaction(rawTransaction);

    const balance = await d.erc20.balanceOf(address);

    expect(balance).toBe(BN(700));
});