const { InternalBuilder } = require("membuffers");

class BaseBuilder {
    constructor() {
        this.builder = new InternalBuilder();
    }

    write(buf) {
        throw new Error("not implemented");
    }

    getSize() {
        return this.builder.getSize();
    }

    calcRequiredSize() {
        this.write(null);
        return this.builder.getSize();
    }

    build() {
        const buf = new Uint8Array(this.calcRequiredSize());
        this.write(buf);
        return buf;
    }
}

module.exports = BaseBuilder;