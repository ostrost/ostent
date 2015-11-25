// required Promises https://github.com/postcss/postcss-nested/issues/30
require('es6-promise').polyfill();

var path              = require('path'),
    webpack           = require('webpack'),
    _                 = require('lodash'),
    gulp              = require('gulp'),
    gutil             = require('gulp-util'),
    argv              = require('yargs').argv,
    ExtractTextPlugin = require('extract-text-webpack-plugin');

var rr;
var wpconf = {
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

gulp.task('wp', [], function(callback) {
    var wparg = wpconf;
    wparg = _.merge(wparg, {entry: argv.input});
    wparg = _.merge(wparg, {output: {}});
    wparg.output.filename = path.basename(argv.output);
    wparg.output.path     = path.join(__dirname, path.dirname(argv.output));
    if (argv.output.match(/\.min\.js($|\?)/i) !== null) {
        wparg.plugins.push(new webpack.optimize.UglifyJsPlugin({mangle: true}));
    }
    webpack(wparg).run(function(err, stats) {
        if(err) {
            throw new gutil.PluginError('webpack', err);
        }
        gutil.log('[webpack]', stats.toString({/* output options */}));
        callback();
    });
});
