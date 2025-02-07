const path = require('path');

module.exports = {
    entry: './src/ts/main.ts',
    output: {
        path: path.resolve(__dirname, './src/js'),
        filename: 'bundle.js',
    },
    module: {
        rules: [
            {
                test: /\.ts$/,
                use: 'ts-loader',
                exclude: /node_modules/,
            },
        ],
    },
    resolve: {
        extensions: ['.ts', '.js'],
    },
};