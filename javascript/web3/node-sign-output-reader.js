const { InternalMessage } = require("membuffers");

const ARGUMENT_TYPE_BYTES_VALUE = 3;

class NodeSignOutputReader {
    constructor(buf) {
        this.message = new InternalMessage(buf, buf.byteLength, ARGUMENT_TYPE_BYTES_VALUE, []);
    }

    getSignature() {
        return this.message.getBytesInOffset(0);
    }
}

module.exports = NodeSignOutputReader;