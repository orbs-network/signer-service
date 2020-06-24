// const Web3 = require("Web3");
const Signer = require("./signer");
const {spawn} = require("child_process");
const { ganacheDriver, Driver } = require("@orbs-network/orbs-ethereum-contracts-v2");
const BN = require('bn.js').BN;

let signerProcess;

beforeAll(async () => {
    signerProcess = spawn("go", [
        "run", "./bootstrap/signer/main/main.go", 
        "-config", "./javascript/web3/keys.json",
    ], {
        cwd: `${__dirname}/../../`
    });
    signerProcess.stdout.on("data", (data) => {
        console.log(`stdout: ${data}`);
    });

    // await ganacheDriver.startGanache();

    await new Promise((resolve, reject) => {
        setTimeout(resolve, 3000);
    });
});

afterAll(async () => {
    // await ganacheDriver.stopGanache();

    const ok = signerProcess.kill();
    console.log(`terminating child process: ${ok}`);
});

const address = "0xa328846cd5b4979d68a8c58a9bdfeee657b34de7";

it("should be able to sign messages", async () => {
    const transactionSigner = new Signer("http://localhost:7777");

    const tx = {
        from: address,
        gasPrice: "0x4A817C800",
        gasLimit: "0x21072",
        to: "0x3535353535353535353535353535353535353535",
        value: "0xDE0B6B3A7640000",
        data: "0x001"
    };

    const signedTransaction = await transactionSigner.sign(tx);
    expect(signedTransaction).toEqual({
        messageHash: "0b412d2c3ee990ea6bcf44a78a601a8cc2245ab5516bbb2f0a9a1199a2b36f85",
        v: "0x1b",
        r: "0x04252883c1abfe7d05ceda1691d0873c5d9525b1e0d0bb478076f7d15e551f23",
        s: "0x16a7a418e2e308006558123023b45f6b4b2733254a6ff820ce16d7d2e2523945",
        rawTransaction: "0xf86f808504a817c80083021072943535353535353535353535353535353535353535880de0b6b3a76400008200011ba004252883c1abfe7d05ceda1691d0873c5d9525b1e0d0bb478076f7d15e551f23a016a7a418e2e308006558123023b45f6b4b2733254a6ff820ce16d7d2e2523945",
        transactionHash: "0x93ede58f6a00299cfa951fcc0846a87304d9efb7576ffb7611110dd298c48f6b"
    });
});

xit("should work with ganache", async () => {
    const d = await Driver.new();
    await d.erc20.assign(address, BN(500));

    const assignTx = {
        from: address,
        to: d.erc20.address,
        gasLimit: 0x7fffffff,
        data: d.erc20.methods.assign(address, BN(200)).encodeABI(),
    };

    const { rawTransaction } = await transactionSigner.sign(assignTx);

    await d.web3.eth.sendSignedTransaction(rawTransaction);

    const balance = await d.erc20.balanceOf(address);

    expect(balance).toBe(BN(700));
});