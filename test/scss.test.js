const { resolve } = require("path")
const { sync }  = require("glob")
const { runSass } = require("sass-true")

const testPath = "test/**/*.test.scss"
const loadPaths = ["src"]

sync(resolve(process.cwd(), testPath)).forEach(function testFile(file) {
  runSass({ describe, it }, file, { loadPaths })
})
