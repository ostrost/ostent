// required Promises https://github.com/postcss/postcss-nested/issues/30
require('es6-promise').polyfill();

var path              = require('path'),
    // webpack        = require('webpack'),
    ExtractTextPlugin = require('extract-text-webpack-plugin');

var rr;
module.exports = {
    resolve: {
        root: (rr = path.join(__dirname, './bower_components')),
        //? extensions: ['', '.js', '.css', '.scss'],
        alias: {
            jquery:          'jquery/dist/jquery'      +'.min.js',
            react:           'react/react-with-addons' +'.min.js',
            'react-dom': rr+'/react/react-dom'         +'.min.js'
        }
    },
    module: {
        loaders: [
            {test: /\.css$/,  loader: ExtractTextPlugin.extract('style-loader', 'css-loader')},
            {test: /\.scss$/, loader: ExtractTextPlugin.extract('style-loader', 'css-loader!sass-loader')}
        ]
    },
    sassLoader: {includePaths: [rr+'/foundation-sites/scss/']},
    plugins: [
        // new webpack.ResolverPlugin(new webpack.ResolverPlugin.DirectoryDescriptionFilePlugin(
        //     'bower.json', ['main'])), // this will resolve with bower.json:"main"
        new ExtractTextPlugin('index.css', {allChunks: true})
    ]
};
