module.exports = {
    outputDir: '../static',
    assetsDir: 'assets',
    devServer: {
        port: 8081,
        proxy: 'http://localhost:8080'
    }
}
