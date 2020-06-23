const NodeSignInputBuilder = require("./node-sign-input-builder");
const NodeSignOutputReader = require("./node-sign-output-reader");

it("encodes and decodes simple byte array", () => {
    const payload = new Uint8Array([1, 2, 3]);

    const builder = new NodeSignInputBuilder(payload);
    expect(builder.build()).toEqual(new Uint8Array([3, 0, 0, 0, 1, 2, 3]));

    const reader = new NodeSignOutputReader(builder.build());
    expect(reader.getSignature()).toEqual(payload);
});