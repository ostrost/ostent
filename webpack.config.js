var path = require('path'); //, webpack = require('webpack');
var rr;
module.exports = {
    resolve: {
        root: (rr = [path.join(__dirname, 'bower_components')]),
        alias: {
            jquery:          'jquery/dist/jquery'      +'.min.js',
            react:           'react/react-with-addons' +'.min.js',
            'react-dom': rr+'/react/react-dom'         +'.min.js'
        }
    }
    // plugins: [new webpack.ResolverPlugin(new webpack.ResolverPlugin.DirectoryDescriptionFilePlugin("bower.json", ["main"]))] // this will resolve with bower.json:"main"
};
