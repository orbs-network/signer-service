const BaseBuilder = require("./base-builder");

class NodeSignInputBuilder extends BaseBuilder {
    constructor(payload) {
        super();

        this.fields = {
            payload,
        };
    }

    write(buf) {
        this.builder.reset();
        this.builder.writeBytes(buf, this.fields.payload);
    }
}

module.exports = NodeSignInputBuilder;