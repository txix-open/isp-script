const subModule = require("./b.js")

function a() {
    return subModule.b()
}

module.exports = {
    a: a
}
