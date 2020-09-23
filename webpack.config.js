var HtmlWebpackPlugin = require('html-webpack-plugin')
var path = require('path')

var devDir = 'dev'

module.exports = {
    mode: 'development',
    entry: path.resolve(devDir, 'index.js'),
    plugins: [new HtmlWebpackPlugin({ template: path.resolve(devDir, 'index.html') })],
    module: {
        rules: [
            {
                test: /\.css$/i,
                use: [
                    'style-loader',
                    {
                        loader: 'css-loader', options: {
                            sourceMap: true
                        }
                    }
                ]
            }
        ]
    },
    devServer: {
        host: '0.0.0.0',
        port: 3000,
        disableHostCheck: true
    }
}
